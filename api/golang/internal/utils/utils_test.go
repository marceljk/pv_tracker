package utils

import (
	"math"
	"testing"
	"time"

	"github.com/marceljk/pv_tracker/api/golang/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestMapDailyForecast(t *testing.T) {
	testCases := []struct {
		name           string
		input          model.ForecastResponseModel
		expectedOutput model.DailyForecast
	}{
		{
			name: "Single day forecast",
			input: model.ForecastResponseModel{
				Forecasts: []model.Forecast{
					{
						PeriodEnd:  time.Date(2023, 10, 27, 12, 0, 0, 0, time.UTC),
						PvEstimate: 100,
					},
					{
						PeriodEnd:  time.Date(2023, 10, 27, 12, 30, 0, 0, time.UTC),
						PvEstimate: 200,
					},
				},
			},
			expectedOutput: model.DailyForecast{
				DailyForecast: map[int64]model.DailyForecastEntry{
					time.Date(2023, 10, 27, 0, 0, 0, 0, time.UTC).UnixMilli(): {
						Estimate: math.Round((float64(100) + float64(200)) / 2),
						Day:      "Freitag",
					},
				},
			},
		},
		{
			name: "Multiple day forecast",
			input: model.ForecastResponseModel{
				Forecasts: []model.Forecast{
					{
						PeriodEnd:  time.Date(2023, 10, 27, 12, 0, 0, 0, time.UTC),
						PvEstimate: 100,
					},
					{
						PeriodEnd:  time.Date(2023, 10, 27, 12, 30, 0, 0, time.UTC),
						PvEstimate: 200,
					},
					{
						PeriodEnd:  time.Date(2023, 10, 28, 12, 0, 0, 0, time.UTC),
						PvEstimate: 300,
					},
					{
						PeriodEnd:  time.Date(2023, 10, 28, 12, 30, 0, 0, time.UTC),
						PvEstimate: 400,
					},
				},
			},
			expectedOutput: model.DailyForecast{
				DailyForecast: map[int64]model.DailyForecastEntry{
					time.Date(2023, 10, 27, 0, 0, 0, 0, time.UTC).UnixMilli(): {
						Estimate: math.Round((float64(100) + float64(200)) / 2),
						Day:      "Freitag",
					},
					time.Date(2023, 10, 28, 0, 0, 0, 0, time.UTC).UnixMilli(): {
						Estimate: math.Round((float64(300) + float64(400)) / 2),
						Day:      "Samstag",
					},
				},
			},
		},
		{
			name: "Empty forecast",
			input: model.ForecastResponseModel{
				Forecasts: []model.Forecast{},
			},
			expectedOutput: model.DailyForecast{
				DailyForecast: map[int64]model.DailyForecastEntry{},
			},
		},
		{
			name: "Multiple entries on same day",
			input: model.ForecastResponseModel{
				Forecasts: []model.Forecast{
					{
						PeriodEnd:  time.Date(2023, 10, 27, 12, 0, 0, 0, time.UTC),
						PvEstimate: 100,
					},
					{
						PeriodEnd:  time.Date(2023, 10, 27, 12, 30, 0, 0, time.UTC),
						PvEstimate: 200,
					},
					{
						PeriodEnd:  time.Date(2023, 10, 27, 13, 00, 0, 0, time.UTC),
						PvEstimate: 300,
					},
					{
						PeriodEnd:  time.Date(2023, 10, 27, 13, 30, 0, 0, time.UTC),
						PvEstimate: 400,
					},
				},
			},
			expectedOutput: model.DailyForecast{
				DailyForecast: map[int64]model.DailyForecastEntry{
					time.Date(2023, 10, 27, 0, 0, 0, 0, time.UTC).UnixMilli(): {
						Estimate: math.Round((float64(100) + float64(200) + float64(300) + float64(400)) / 2),
						Day:      "Freitag",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput := MapDailyForecast(tc.input)
			assert.Equal(t, tc.expectedOutput, actualOutput)
		})
	}
}

func TestWeekdayToGermanDays(t *testing.T) {
	testCases := []struct {
		name           string
		input          time.Weekday
		expectedOutput string
	}{
		{
			name:           "Sunday",
			input:          time.Sunday,
			expectedOutput: "Sonntag",
		},
		{
			name:           "Monday",
			input:          time.Monday,
			expectedOutput: "Montag",
		},
		{
			name:           "Tuesday",
			input:          time.Tuesday,
			expectedOutput: "Dienstag",
		},
		{
			name:           "Wednesday",
			input:          time.Wednesday,
			expectedOutput: "Mittwoch",
		},
		{
			name:           "Thursday",
			input:          time.Thursday,
			expectedOutput: "Donnerstag",
		},
		{
			name:           "Friday",
			input:          time.Friday,
			expectedOutput: "Freitag",
		},
		{
			name:           "Saturday",
			input:          time.Saturday,
			expectedOutput: "Samstag",
		},
		{
			name:           "Invalid Weekday",
			input:          time.Weekday(7),
			expectedOutput: "%!Weekday(7)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput := weekdayToGermanDays(tc.input)
			assert.Equal(t, tc.expectedOutput, actualOutput)
		})
	}
}
