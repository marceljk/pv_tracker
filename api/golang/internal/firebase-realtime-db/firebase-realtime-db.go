package firebaserealtimedb

import (
	"context"
	"fmt"
	"sync"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/marceljk/pv_tracker/api/golang/internal/model"
	"google.golang.org/api/option"
)

type Database struct {
	live           *db.Ref
	hourlyForecast *db.Ref
	dailyForecast  *db.Ref
	history        *db.Ref
	dailySum       *db.Ref
	cache          *cacheStruct
}

type cacheStruct struct {
	nextUpdate time.Time
	history    map[string]model.PvData
}

type DailySumStruct struct {
	model.PvData
	GridPowerIn  int `json:"gridPowerIn"`
	GridPowerOut int `json:"gridPowerOut"`
	Count        int `json:"count"`
}

func NewFirebaseDbClient(ctx context.Context, config *firebase.Config, opt ...option.ClientOption) (*Database, error) {
	app, err := firebase.NewApp(ctx, config, opt...)
	if err != nil {
		return nil, fmt.Errorf("failed to create firebase app: %w", err)
	}

	db, err := app.Database(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to firebase db: %w", err)
	}

	cache := &cacheStruct{}

	return &Database{
		live:           db.NewRef("live"),
		hourlyForecast: db.NewRef("hourlyForecast"),
		dailyForecast:  db.NewRef("dailyForecast"),
		history:        db.NewRef("history"),
		dailySum:       db.NewRef("dailySum"),
		cache:          cache,
	}, nil
}

func (d *Database) GetLive(ctx context.Context) (*model.PvData, error) {
	var pvData model.PvData
	if err := d.live.Get(ctx, &pvData); err != nil {
		return nil, fmt.Errorf("failed get live data: %w", err)
	}
	return &pvData, nil
}

func (d *Database) SetLive(ctx context.Context, pvData *model.PvData) error {
	if err := d.live.Set(ctx, pvData); err != nil {
		return fmt.Errorf("failed set live data: %w", err)
	}
	return nil
}

func (d *Database) SetHourlyForecast(ctx context.Context, hourlyForecast *model.ForecastResponseModel) error {
	if err := d.hourlyForecast.Set(ctx, hourlyForecast.Forecasts); err != nil {
		return fmt.Errorf("failed set hourly forecast: %w", err)
	}
	return nil
}

func (d *Database) SetDailyForecast(ctx context.Context, dailyForecast *model.DailyForecast) error {
	if err := d.dailyForecast.Set(ctx, dailyForecast.DailyForecast); err != nil {
		return fmt.Errorf("failed set daily forecast: %w", err)
	}
	return nil
}

func (d *Database) SetHistory(ctx context.Context, pvData *model.PvData) error {
	t := time.Now().Format(time.RFC3339)[:19] // timestamp in RFC3339 format without zone
	if err := d.history.Child(t).Set(ctx, pvData); err != nil {
		return fmt.Errorf("failed set history: %w", err)
	}
	if d.cache.history != nil {
		d.cache.history[t] = *pvData
	}
	return nil
}

func (d *Database) UpdateDailySum(ctx context.Context) error {
	dailySum, err := d.calcDailySum(ctx)
	if err != nil {
		return fmt.Errorf("failed calc dailySum: %w", err)
	}
	err = d.saveDailySum(ctx, dailySum)
	if err != nil {
		return fmt.Errorf("failed saving dailySum: %w", err)
	}
	return nil
}

func (d *Database) CleanHistoryUntil(ctx context.Context, t time.Time) error {
	formattedT := t.Format(time.RFC3339)[:19]
	historyChilds := make(map[string]model.PvData)
	err := d.history.OrderByKey().EndAt(formattedT).LimitToLast(5000).Get(ctx, &historyChilds)
	if err != nil {
		return fmt.Errorf("failed get history until %q: %w", formattedT, err)
	}
	var wg sync.WaitGroup
	for childKey := range historyChilds {
		wg.Add(1)
		go func(innerCtx context.Context, key string) {
			err = d.history.Child(key).Delete(innerCtx)
			if err != nil {
				fmt.Printf("[cleanup] - failed deleting child %q: %v\n", key, err)
			}
			wg.Done()
		}(ctx, childKey)
	}
	wg.Wait()

	return nil
}

// caches the data from history for an hour. It also contains the data from `SetHistory()`. This is to avoid high
// consumption of the network traffic on firebase
func (d *Database) getHistoryCached(ctx context.Context) (map[string]model.PvData, error) {
	source := "cache"
	if len(d.cache.history) == 0 || d.cache.nextUpdate.Before(time.Now()) {
		historyChilds := make(map[string]model.PvData)
		if err := d.history.OrderByKey().Get(ctx, &historyChilds); err != nil {
			return nil, fmt.Errorf("failed get history: %w", err)
		}
		d.cache.history = historyChilds
		d.cache.nextUpdate = time.Now().Add(time.Hour)
		source = "database"
	}

	fmt.Printf("[daily sum] - get data from %s\n", source)

	return d.cache.history, nil
}

func (d *Database) calcDailySum(ctx context.Context) (map[string]DailySumStruct, error) {
	historyChilds, err := d.getHistoryCached(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed set daily sum: %w", err)
	}

	dailySum := make(map[string]DailySumStruct)

	for key, child := range historyChilds {
		day := key[:11]
		var gridPowerIn int
		var gridPowerOut int
		if child.GridPowerW > 0 {
			gridPowerIn = child.GridPowerW
		} else {
			gridPowerOut = child.GridPowerW
		}

		if existing, ok := dailySum[day]; ok {
			dailySum[day] = DailySumStruct{
				PvData: model.PvData{
					BatteryPercent:    existing.BatteryPercent + child.BatteryPercent,
					BatteryPowerW:     existing.BatteryPowerW + child.BatteryPowerW,
					PowerConsumptionW: existing.PowerConsumptionW + child.PowerConsumptionW,
					PvPowerW:          existing.PvPowerW + child.PvPowerW,
				},
				GridPowerIn:  existing.GridPowerIn + gridPowerIn,
				GridPowerOut: existing.GridPowerOut + gridPowerOut,
				Count:        existing.Count + 1,
			}
		} else {
			dailySum[day] = DailySumStruct{
				PvData: model.PvData{
					BatteryPercent:    child.BatteryPercent,
					BatteryPowerW:     child.BatteryPowerW,
					PowerConsumptionW: child.PowerConsumptionW,
				},
				GridPowerIn:  gridPowerIn,
				GridPowerOut: gridPowerOut,
				Count:        1,
			}
		}
	}

	return dailySum, nil
}

func (d *Database) saveDailySum(ctx context.Context, dailySum map[string]DailySumStruct) error {
	var errs []error
	for day, value := range dailySum {
		if err := d.dailySum.Child(day).Set(ctx, value); err != nil {
			errs = append(errs, fmt.Errorf("failed set %q in dailySum: %w", day, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("failed saving dailySums: %v", errs)
	}
	return nil
}
