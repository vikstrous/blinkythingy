package http

import (
	"encoding/json"
	"net"
	"net/http"
	"path"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/Sirupsen/logrus"
	"github.com/vikstrous/blinkythingy"
	"github.com/vikstrous/blinkythingy/display"
)

type HTTPConfig struct {
	Addr string
	Path string
}

func MapToBlinkyConfig(mapConfig blinkythingy.MapConfig) (HTTPConfig, error) {
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

type httpDisplay struct {
	listener *net.TCPListener
	colors   []blinkythingy.Color
}

func New(mapConfig blinkythingy.MapConfig) (display.Display, error) {
	config, err := MapToBlinkyConfig(mapConfig)
	if err != nil {
		return nil, err
	}

	addr := config.Addr
	servePath := config.Path

	if addr == "" {
		addr = ":1337"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	listener := ln.(*net.TCPListener)

	d := &httpDisplay{
		listener: listener,
	}

	mux := http.NewServeMux()
	mux.HandleFunc(path.Join("/", servePath), func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(d.colors)
	})
	server := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// TODO: provide a mechanism for killing this by calling listener.Close()
	go func() {
		err := server.Serve(listener)
		if err != nil {
			logrus.Warn(err)
		}
	}()
	return d, nil
}

func (d *httpDisplay) Flush(colors []blinkythingy.Color) error {
	d.colors = colors
	return nil
}
