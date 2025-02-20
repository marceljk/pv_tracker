package model

import "time"

type ForecastResponseModel struct {
	Forecasts []Forecast `json:"forecasts"`
}

type Forecast struct {
	PvEstimate float32   `json:"pv_estimate"`
	PeriodEnd  time.Time `json:"period_end"`
}

type DailyForecast struct {
	DailyForecast map[int64]DailyForecastEntry `json:"dailyForecast"`
}

type DailyForecastEntry struct {
	Estimate float64 `json:"estimate"`
	Day      string  `json:"day"`
}
