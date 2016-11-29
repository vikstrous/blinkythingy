package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/vikstrous/blinkythingy"
	"github.com/vikstrous/blinkythingy/fetcher"
	"github.com/vikstrous/blinkythingy/util"
)

type HTTPConfig struct {
	blinkythingy.HTTPClientConfig
	URL string
}

func MapToHTTPConfig(mapConfig blinkythingy.MapConfig) (HTTPConfig, error) {
	config := HTTPConfig{}
	marshalled, err := yaml.Marshal(mapConfig)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(marshalled, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

type httpFetcher struct {
	url    string
	colors []blinkythingy.Color
	client *http.Client
}

func New(mapConfig blinkythingy.MapConfig) (fetcher.Fetcher, error) {
	config, err := MapToHTTPConfig(mapConfig)
	if err != nil {
		return nil, err
	}
	httpClient, err := util.HTTPClient(config.InsecureTLS, config.CA)
	if err != nil {
		return nil, err
	}
	return &httpFetcher{
		url:    config.URL,
		client: httpClient,
	}, nil
}

func (f *httpFetcher) FetchStatuses() error {
	req, err := http.NewRequest("GET", f.url, nil)
	if err != nil {
		return err
	}
	res, err := f.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("Bad status: %d", res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(&f.colors)
}

func (f *httpFetcher) ListStatuses() []blinkythingy.Color {
	return f.colors
}
