package jenkinsapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type client struct {
	host     string
	username string
	password string
	client   *http.Client
}

type Client interface {
	JobDataWithFilter(job, filter string) (interface{}, error)
}

func New(host, username, password string, httpClient *http.Client) Client {
	return &client{
		host:     host,
		username: username,
		password: password,
		client:   httpClient,
	}
}

func (c *client) JobDataWithFilter(job, filter string) (interface{}, error) {
	res, err := c.makeRequest("GET", fmt.Sprintf("/job/%s/api/json?tree=%s", job, filter), nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Non-200 status code: %s", res.StatusCode)
	}
	parsed := map[string]interface{}{}
	err = json.NewDecoder(res.Body).Decode(&parsed)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func (c *client) makeRequest(method, route string, payload interface{}) (*http.Response, error) {
	var reader io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("https://%s%s", c.host, route), reader)
	if err != nil {
		return nil, err
	}
	if method != "HEAD" && method != "GET" {
		req.Header.Add("Content-Type", "application/json")
	}

	req.SetBasicAuth(c.username, c.password)

	// always close the connection after sending the request so that we
	// don't get bitten by net/http's bugs with reusing connections
	req.Close = true

	return c.client.Do(req)
}
