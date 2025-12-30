package firebaserealtimedb

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
	databaseURL    string
}

type DailySumStruct struct {
	model.PvData
	GridPowerIn  int `json:"gridPowerIn"`
	GridPowerOut int `json:"gridPowerOut"`
	Count        int `json:"count"`
}

type sseEvent struct {
	Event string          `json:"event"`
	Path  string          `json:"path"`
	Data  json.RawMessage `json:"data"`
}

func NewFirebaseDbClient(ctx context.Context, config *firebase.Config, opt ...option.ClientOption) (*Database, error) {
	app, err := firebase.NewApp(ctx, config, opt...)
	if err != nil {
		return nil, fmt.Errorf("failed to create firebase app: %w", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to firebase db: %w", err)
	}

	return &Database{
		live:           client.NewRef("live"),
		hourlyForecast: client.NewRef("hourlyForecast"),
		dailyForecast:  client.NewRef("dailyForecast"),
		history:        client.NewRef("history"),
		dailySum:       client.NewRef("dailySum"),
		databaseURL:    config.DatabaseURL,
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
	return nil
}

func (d *Database) UpdateDailySum(ctx context.Context, CtxCancel context.CancelFunc) {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02T15:04:05")

	u, err := url.Parse(d.databaseURL)
	if err != nil {
		fmt.Printf("Stream failed: %v\n", err)
		CtxCancel()
		return
	}
	u.Path = "/history.json"
	q := u.Query()
	q.Set("orderBy", `"$key"`)
	q.Set("startAt", fmt.Sprintf(`"%s"`, yesterday))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		fmt.Printf("Stream failed: %v\n", err)
		CtxCancel()
		return
	}
	req.Header.Set("Accept", "text/event-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Stream failed: %v\n", err)
		CtxCancel()
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	// Set a larger buffer size for the scanner to prevent "token too long" errors
	// when processing large data payloads from the Firebase stream.
	const maxTokenSize = 10 * 1024 * 1024 // 10MB
	scanner.Buffer(make([]byte, 0, 64*1024), maxTokenSize)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		jsonData := strings.TrimPrefix(line, "data: ")
		if jsonData == "null" { // for keep-alive events
			continue
		}

		var event sseEvent
		if err := json.Unmarshal([]byte(jsonData), &event); err != nil {
			fmt.Printf("UpdateDailySum: failed to unmarshal event data: %v\n", err)
			continue
		}

		if event.Path == "/" {
			var data map[string]model.PvData
			if err := json.Unmarshal(event.Data, &data); err != nil {
				fmt.Printf("UpdateDailySum: failed to unmarshal initial data: %v\n", err)
				continue
			}

			//do initial calculation
			for key, value := range data {
				if len(key) >= 10 {
					day := key[:10]
					if err := d.updateDailySum(ctx, day, &value); err != nil {
						fmt.Printf("UpdateDailySum: failed to update sum for day %s: %v\n", day, err)
					}
				}
			}
		} else {
			if len(event.Path) < 11 || !strings.HasPrefix(event.Path, "/") {
				continue
			}
			if string(event.Data) == "null" { // Data for deleted node
				continue
			}

			var data model.PvData
			if err := json.Unmarshal(event.Data, &data); err != nil {
				fmt.Printf("UpdateDailySum: failed to unmarshal update data for path %s: %v\n", event.Path, err)
				continue
			}

			day := event.Path[1:11]
			if strings.Contains(day, "/") {
				continue
			}
			if err := d.updateDailySum(ctx, day, &data); err != nil {
				fmt.Printf("UpdateDailySum: failed to update sum for day %s: %v\n", day, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Stream failed: %v\n", err)
	}
	CtxCancel()
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

func (d *Database) updateDailySum(ctx context.Context, day string, pvData *model.PvData) error {
	var dailySum DailySumStruct
	if err := d.dailySum.Child(day).Get(ctx, &dailySum); err != nil {
		return fmt.Errorf("failed get dailySum for day %s: %w", day, err)
	}

	var gridPowerIn int
	var gridPowerOut int
	if pvData.GridPowerW > 0 {
		gridPowerIn = pvData.GridPowerW
	} else {
		gridPowerOut = pvData.GridPowerW
	}

	if dailySum.Count > 0 {
		dailySum = DailySumStruct{
			PvData: model.PvData{
				BatteryPercent:    dailySum.BatteryPercent + pvData.BatteryPercent,
				BatteryPowerW:     dailySum.BatteryPowerW + pvData.BatteryPowerW,
				PowerConsumptionW: dailySum.PowerConsumptionW + pvData.PowerConsumptionW,
				PvPowerW:          dailySum.PvPowerW + pvData.PvPowerW,
			},
			GridPowerIn:  dailySum.GridPowerIn + gridPowerIn,
			GridPowerOut: dailySum.GridPowerOut + gridPowerOut,
			Count:        dailySum.Count + 1,
		}
	} else {
		dailySum = DailySumStruct{
			PvData: model.PvData{
				BatteryPercent:    pvData.BatteryPercent,
				BatteryPowerW:     pvData.BatteryPowerW,
				PowerConsumptionW: pvData.PowerConsumptionW,
				PvPowerW:          pvData.PvPowerW,
			},
			GridPowerIn:  gridPowerIn,
			GridPowerOut: gridPowerOut,
			Count:        1,
		}
	}

	if err := d.dailySum.Child(day).Set(ctx, dailySum); err != nil {
		return fmt.Errorf("failed to set dailySum for day %s: %w", day, err)
	}
	return nil
}

func (e *sseEvent) UnmarshalJSON(b []byte) error {
	var eventData struct {
		Path string          `json:"path"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(b, &eventData); err != nil {
		return err
	}
	e.Path = eventData.Path
	e.Data = eventData.Data

	return nil
}
