package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vikstrous/blinkythingy"
	"github.com/vikstrous/blinkythingy/fetcher"
)

type serverFetcher struct {
	url    string
	colors []blinkythingy.Color
	client *http.Client
}

func New(url string, httpClient *http.Client) fetcher.Fetcher {
	return &serverFetcher{
		url:    url,
		client: httpClient,
	}
}

func (f *serverFetcher) FetchStatuses() error {
	req := http.NewRequest("GET", url, nil)
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("Bad status: %s")
	}
	return json.Decoder(res).Decode(d.colors)
}

func (f *serverFetcher) ListStatuses() []blinkythingy.Color {
	return d.colors
}
