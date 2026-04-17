package provider

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPProvider struct {
	addr   string
	client *http.Client
}

func NewHTTPProvider(addr string) (Provider, error) {
	if addr == "" {
		return nil, fmt.Errorf("addr is empty for http provider")
	}

	return &HTTPProvider{
		addr: addr,
		client: &http.Client{
			Timeout: 5,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 2,
				IdleConnTimeout:     30 * time.Second,

				TLSHandshakeTimeout:   3 * time.Second,
				ResponseHeaderTimeout: 3 * time.Second,
			},
		},
	}, nil
}

func (h *HTTPProvider) Disconnect() {}

func (h *HTTPProvider) Get() ([]byte, error) {
	response, err := h.client.Get(h.addr)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var body []byte
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
