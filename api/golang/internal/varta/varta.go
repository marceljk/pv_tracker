package varta

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/marceljk/pv_tracker/api/golang/internal/model"
)

const (
	baseUrl   = "http://varta130104162"
	loginPath = "/cgi/login"
	dataPath  = "/cgi/data"

	loginHeaderCred = "set-cookie"
)

type Repo struct {
	baseUrl     string
	headers     http.Header
	username    string
	password    string
	credentials string
}

func NewVartaRepo(username, password string) *Repo {
	repo := &Repo{
		baseUrl:  baseUrl,
		headers:  nil,
		username: username,
		password: password,
	}
	return repo
}

func (r *Repo) Login() error {
	client := http.Client{}

	form := url.Values{}
	form.Add("username", r.username)
	form.Add("password", r.password)

	reqUrl := fmt.Sprintf("%s%s", baseUrl, loginPath)
	req, err := http.NewRequest(http.MethodPost, reqUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("failed create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("api call: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("api responde with status code %d", resp.StatusCode)
	}

	if header := resp.Header; header != nil {
		r.credentials = header.Get("set-cookie")
	}

	return nil
}

func (r *Repo) GetLiveData() (*model.PvData, error) {
	resp, err := r.fetchLiveData()
	if err != nil {
		return nil, fmt.Errorf("could not fetch live data: %w", err)
	}

	result := resp.mapToPvData()
	return &result, nil
}

func (r *Repo) fetchLiveData() (*rawApiResponse, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	url := fmt.Sprintf("%s%s", baseUrl, dataPath)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("cookie", r.credentials)
	response, err := client.Do(req)
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
		return nil, fmt.Errorf("cannot unmarshal response: %w", err)
	}

	return resp, nil
}
