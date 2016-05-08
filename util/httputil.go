package util

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

func HTTPClient(insecure bool, cas ...string) (*http.Client, error) {
	client := new(http.Client)
	var err error
	client.Transport, err = httpTransport(insecure, cas...)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func httpTransport(insecure bool, cas ...string) (*http.Transport, error) {
	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		//ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure,
		},
	}
	if !insecure {
		caPool := systemRootsPool()
		for _, ca := range cas {
			if ca == "" {
				continue
			}
			ok := caPool.AppendCertsFromPEM([]byte(ca))
			if !ok {
				return nil, fmt.Errorf("TLS CA provided, but we could not parse a PEM encoded cert from it. CA provided: \n%s", ca)
			}
		}
		transport.TLSClientConfig.RootCAs = caPool
	}
	return &transport, nil
}
