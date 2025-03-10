package jsonservice

import (
	"net/http"
	"time"

	"github.com/motain/of-catalog/internal/services/configservice"
)

type JSONServiceInterface interface {
	Do(*http.Request) (*http.Response, error)
}

type JSONTransport struct {
	Transport http.RoundTripper
}

func NewJSONService(config configservice.ConfigServiceInterface) JSONServiceInterface {
	baseTransport := &http.Transport{
		MaxIdleConns:      10,
		IdleConnTimeout:   30 * time.Second,
		DisableKeepAlives: false,
	}

	return &http.Client{
		Transport: &JSONTransport{
			Transport: baseTransport,
		},
	}
}

func (c *JSONTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return c.Transport.RoundTrip(req)
}
