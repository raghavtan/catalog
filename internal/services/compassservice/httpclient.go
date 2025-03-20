package compassservice

//go:generate mockgen -destination=./mocks/mock_http_client.go -package=compassservice github.com/motain/of-catalog/internal/services/compassservice HTTPClientInterface

import (
	"net/http"
	"strings"
	"time"

	"github.com/motain/of-catalog/internal/services/configservice"
)

type HTTPClientInterface interface {
	Do(*http.Request) (*http.Response, error)
}

type CompassTransport struct {
	Transport http.RoundTripper
	BaseURL   string
	Host      string
	AuthToken string
}

func NewHTTPClient(config configservice.ConfigServiceInterface) HTTPClientInterface {
	baseTransport := &http.Transport{
		MaxIdleConns:      10,
		IdleConnTimeout:   30 * time.Second,
		DisableKeepAlives: false,
	}

	return &http.Client{
		Transport: &CompassTransport{
			Transport: baseTransport,
			Host:      config.GetCompassHost(),
			AuthToken: config.GetCompassToken(),
		},
	}
}

func (c *CompassTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !strings.HasPrefix(req.URL.String(), "http") {
		req.URL.Scheme = "https"
		req.URL.Host = c.Host
		req.URL.Path = c.BaseURL + req.URL.Path
	}

	req.Host = c.Host
	req.Header.Set("Authorization", "Basic "+c.AuthToken)

	return c.Transport.RoundTrip(req)
}
