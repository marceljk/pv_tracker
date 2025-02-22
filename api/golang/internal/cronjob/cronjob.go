package cronjob

import (
	"context"
	"fmt"
	"time"

	"github.com/marceljk/pv_tracker_go/internal"
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
		fmt.Printf("[live] - update live\n")
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

	// Every 6 hours cleans history which is older than 24 hours
	c.AddFunc("0 0 */6 * * *", func() {
		ctx := context.Background()
		durationOneDay := time.Hour * 24
		t := time.Now().Add(durationOneDay * (-1))
		err := db.CleanHistoryUntil(ctx, t)
		if err != nil {
			fmt.Printf("[cleanup] - failed cleaning history until %q %v\n", t.String(), err)
		} else {
			fmt.Printf("[cleanup] - cleaned history until %q\n", t.String())
		}
	})
}

func (c *Cronjob) Start() {
	c.cron.Start()
}
