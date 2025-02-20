package utils

import (
	"math"
	"time"

	"github.com/marceljk/pv_tracker_go/internal/model"
)

func MapDailyForecast(u model.ForecastResponseModel) model.DailyForecast {
	daily := model.DailyForecast{
		DailyForecast: make(map[int64]model.DailyForecastEntry),
	}
	for idx := range u.Forecasts {
		momentForecast := u.Forecasts[idx]
		timestamp := momentForecast.PeriodEnd
		// Set to date and 00:00 time
		date := time.Date(
			timestamp.Year(),
			timestamp.Month(),
			timestamp.Day(),
			0,
			0,
			0,
			0,
			time.UTC,
		)
		unixTimestamp := date.UnixMilli()
		if val, ok := daily.DailyForecast[unixTimestamp]; ok {
			daily.DailyForecast[unixTimestamp] = model.DailyForecastEntry{
				Estimate: val.Estimate + float64(momentForecast.PvEstimate),
				Day:      val.Day,
			}
		} else {
			estimate := float64(momentForecast.PvEstimate)
			day := weekdayToGermanDays(date.Weekday())
			daily.DailyForecast[unixTimestamp] = model.DailyForecastEntry{
				Estimate: estimate,
				Day:      day,
			}
		}
	}

	// Calculate correct value. The valus in estimate need to be divided by 2 because the period for them in the forecast are 30 minutes
	for idx, val := range daily.DailyForecast {
		daily.DailyForecast[idx] = model.DailyForecastEntry{
			Day:      val.Day,
			Estimate: math.Round(val.Estimate / 2),
		}
	}
	return daily
}

func weekdayToGermanDays(w time.Weekday) string {
	weekdays := [7]string{"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag"}
	if 0 < int(w) && int(w) < len(weekdays) {
		return weekdays[w]
	}
	return w.String()
}
