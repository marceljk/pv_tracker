package cronjob

import (
	"context"
	"fmt"
	"time"

	"github.com/marceljk/pv_tracker/api/golang/internal"
	"github.com/robfig/cron/v3"
)

type Cronjob struct {
	cron *cron.Cron
}

func NewCronjob(pvRepo internal.PvRepository, db internal.Database, forecastRepo internal.ForecastRepository) *Cronjob {
	c := cron.New(cron.WithSeconds())

	initCronJobs(c, pvRepo, db, forecastRepo)

	return &Cronjob{
		cron: c,
	}
}

func initCronJobs(c *cron.Cron, pvRepo internal.PvRepository, db internal.Database, forecastRepo internal.ForecastRepository) {
	// Every 3 seconds update live data
	c.AddFunc("*/3 * * * * *", func() {
		ctx := context.Background()
		// fmt.Printf("[live] - update live\n")
		updateLiveData(ctx, pvRepo, db)
	})

	// Every 30 seconds push history data
	c.AddFunc("*/30 * * * * *", func() {
		ctx := context.Background()
		fmt.Printf("[history] - update history\n")
		pushHistoryData(ctx, pvRepo, db)
		fmt.Printf("[daily sum] - update daily sum\n")
		db.UpdateDailySum(ctx)
	})

	// Every 30 min update forecast
	c.AddFunc("0 */30 * * * *", func() {
		ctx := context.Background()
		fmt.Printf("[forecast] - update forecast\n")
		updateForecast(ctx, db, forecastRepo)
	})

	// Cleans history that is from the past day
	c.AddFunc("0 0 3,6 * * *", func() {
		ctx := context.Background()
		now := time.Now()
		year, month, day := now.Date()
		midnight := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
		err := db.CleanHistoryUntil(ctx, midnight)
		if err != nil {
			fmt.Printf("[cleanup] - failed cleaning history until %q %v\n", midnight.String(), err)
		} else {
			fmt.Printf("[cleanup] - cleaned history until %q\n", midnight.String())
		}
	})

	c.AddFunc("0 0 */6 * * *", func() {
		if err := pvRepo.Login(); err != nil {
			fmt.Printf("[login] - failed to refresh credentials\n")
		} else {
			fmt.Printf("[login] - refreshed credentials\n")
		}
	})
}

func (c *Cronjob) Start() {
	c.cron.Start()
}
