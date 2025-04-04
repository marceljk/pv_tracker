package internal

import (
	"context"
	"time"

	"github.com/marceljk/pv_tracker/api/golang/internal/model"
)

type PvRepository interface {
	GetLiveData() (*model.PvData, error)
}

type Database interface {
	SetLive(context.Context, *model.PvData) error
	SetHourlyForecast(context.Context, *model.ForecastResponseModel) error
	SetDailyForecast(context.Context, *model.DailyForecast) error
	SetHistory(context.Context, *model.PvData) error
	UpdateDailySum(context.Context) error
	CleanHistoryUntil(context.Context, time.Time) error
}

type ForecastRepository interface {
	GetForecast(ctx context.Context) (resp *model.ForecastResponseModel, err error)
}
