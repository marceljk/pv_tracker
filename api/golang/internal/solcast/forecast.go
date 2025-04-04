package solcast

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/marceljk/pv_tracker/api/golang/internal/model"
)

type SolcastClient struct {
	endpoint string
	apiKey   string
	useFile  bool
}

func NewSolcastClient(endpoint, apiKey string, debugMode bool) *SolcastClient {
	return &SolcastClient{
		endpoint: endpoint,
		apiKey:   apiKey,
		useFile:  debugMode,
	}
}

func (f *SolcastClient) GetForecast(ctx context.Context) (resp *model.ForecastResponseModel, err error) {
	if f.useFile {
		resp, err = f.readFromFile()
	} else {
		resp, err = f.fetchForecast()
	}
	return
}

func (f *SolcastClient) readFromFile() (*model.ForecastResponseModel, error) {
	data, err := os.ReadFile("forecast_response.json")
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	var resp model.ForecastResponseModel
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal forecast_response.json: %w", err)
	}
	return &resp, nil
}

func (f *SolcastClient) fetchForecast() (*model.ForecastResponseModel, error) {
	req, err := http.NewRequest(http.MethodGet, f.endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed building request: %w", err)
	}
	req.Header.Add("Authorization", f.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %w", err)
	}
	var forecastResponse model.ForecastResponseModel
	err = json.Unmarshal(data, &forecastResponse)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal response: %w", err)
	}
	os.WriteFile("forecast_response.json", data, 0o777)
	return &forecastResponse, nil
}
