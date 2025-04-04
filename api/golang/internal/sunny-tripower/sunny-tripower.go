package sunnytripower

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/marceljk/pv_tracker/api/golang/internal/model"
)

type Repo struct {
	baseUrl string
}

func NewRepo(baseUrl string) *Repo {
	return &Repo{
		baseUrl: baseUrl,
	}
}

func (r *Repo) GetLiveData() (*model.PvData, error) {
	result := &model.PvData{}

	resp, err := r.fetchLiveData()
	if err != nil {
		return nil, fmt.Errorf("could not fetch live data: %w", err)
	}

	result.PvPowerW = resp.PvPowerW
	return result, nil
}

func (r *Repo) fetchLiveData() (*livePowerResponse, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	url := fmt.Sprintf("%s%s", r.baseUrl, "/dyn/getDashValues.json")
	response, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch live data failed: %w", err)
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body: %w", err)
	}
	resp := &rawApiResponse{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal body: %w", err)
	}

	emptyResponse := &livePowerResponse{
		PvPowerW: 0,
	}

	if len(resp.Result.Zero17AXxxxx6B9.Six1000046C200.Num1) != 1 {
		return emptyResponse, nil
	}

	if val := resp.Result.Zero17AXxxxx6B9.Six1000046C200.Num1[0].Val; val != nil {
		return &livePowerResponse{*val}, nil
	}
	return emptyResponse, nil
}
