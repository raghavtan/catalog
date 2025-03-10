package compassservice

//go:generate mockgen -destination=./mock_compass_service.go -package=compassservice github.com/motain/of-catalog/internal/services/compassservice CompassServiceInterface

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/machinebox/graphql"
	"github.com/motain/of-catalog/internal/services/configservice"
)

type CompassServiceInterface interface {
	Run(ctx context.Context, query string, variables map[string]interface{}, response interface{}) error
	SendMetric(body map[string]string) (string, error)
	GetCompassCloudId() string
}

type CompassService struct {
	gqlClient  GraphQLClientInterface
	httpClient HTTPClientInterface
	token      string
	cloudId    string
}

func NewCompassService(
	config configservice.ConfigServiceInterface,
	gqlClient GraphQLClientInterface,
	httpClient HTTPClientInterface,
) *CompassService {
	return &CompassService{
		gqlClient:  gqlClient,
		httpClient: httpClient,
		token:      config.GetCompassToken(),
		cloudId:    config.GetCompassCloudId(),
	}
}

func (c *CompassService) Run(ctx context.Context, query string, variables map[string]interface{}, response interface{}) error {
	req := graphql.NewRequest(query)
	for key, value := range variables {
		req.Var(key, value)
	}

	req.Header.Set("Authorization", "Basic "+c.token)

	if err := c.gqlClient.Run(ctx, req, response); err != nil {
		log.Printf("Failed to execute query: %v", err)
		return err
	}

	return nil
}

func (c *CompassService) SendMetric(body map[string]string) (string, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal body: %v", err)
	}
	req, err := http.NewRequest("POST", "/gateway/api/compass/v1/metrics", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response body: %v", err)
		}
		return "", fmt.Errorf("response body: %s", string(body))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(respBody), nil
}

func (c *CompassService) GetCompassCloudId() string {
	return c.cloudId
}
