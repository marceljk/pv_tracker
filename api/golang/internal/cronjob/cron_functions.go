package cronjob

import (
	"context"
	"fmt"

	"github.com/marceljk/pv_tracker_go/internal"
	"github.com/marceljk/pv_tracker_go/internal/utils"
)

func updateLiveData(ctx context.Context, pvRepo internal.PvRepository, db internal.Database) {
	data, err := pvRepo.GetLiveData()
	if err != nil {
		fmt.Printf("failed fetching live data: %v\n", err)
		return
	}

	err = db.SetLive(ctx, data)
	if err != nil {
		fmt.Printf("failed setting live data in db: %v\n", err)
	}
}

func pushHistoryData(ctx context.Context, pvRepo internal.PvRepository, db internal.Database) {
	data, err := pvRepo.GetLiveData()
	if err != nil {
		fmt.Printf("failed fetching live data: %v\n", err)
		return
	}

	err = db.SetHistory(ctx, data)
	if err != nil {
		fmt.Printf("failed setting history data in db: %v\n", err)
	}
}

func updateForecast(ctx context.Context, db internal.Database, forecastRepo internal.ForecastRepository) {
	hourlyData, err := forecastRepo.GetForecast(ctx)
	if err != nil {
		fmt.Printf("failed fetching forecast data: %v\n", err)
		return
	}

	dailyData := utils.MapDailyForecast(*hourlyData)

	err = db.SetHourlyForecast(ctx, hourlyData)
	if err != nil {
		fmt.Printf("failed updating hourly forecast: %v\n", err)
	} else {
		fmt.Printf("updated hourly forecast\n")
	}

	err = db.SetDailyForecast(ctx, &dailyData)
	if err != nil {
		fmt.Printf("failed updating daily forecast: %v\n", err)
	} else {
		fmt.Printf("updated daily forecast\n")
	}
}
