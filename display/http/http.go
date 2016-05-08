package http

import (
	"encoding/json"
	"net"
	"net/http"
	"path"
	"time"

	"github.com/vikstrous/blinkythingy"
	"github.com/vikstrous/blinkythingy/display"
)

type httpDisplay struct {
	listener *net.TCPListener
	colors   []blinkythingy.Color
}

func New(addr, serverPath string) (display.Display, error) {
	if addr == "" {
		addr = ":1337"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	listener := tcpKeepAliveListener{ln.(*net.TCPListener)}

	d := &httpDisplay{
		listener: listener,
	}

	mux := http.NewServeMux()
	mux.HandleFunc(path.Join("/", servePath), func(w http.ResponseWriter, r *http.Request) {
		json.Encoder(w).Encode(d.colors)
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
			fmt.Warn(err)
		}
	}()
	return d, nil
}

func (d *httpDisplay) Flush(colors []caashttp.Color) error {
	d.colors = colors
}
