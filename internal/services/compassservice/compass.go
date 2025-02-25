package compassservice

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/machinebox/graphql"
	"github.com/motain/fact-collector/internal/services/configservice"
)

type CompassServiceInterface interface {
	Run(ctx context.Context, query string, variables map[string]interface{}, response interface{}) error
	SendMetric(body map[string]string) string
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

func (c *CompassService) SendMetric(body map[string]string) string {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("Failed to marshal body: %v", err)
		return ""
	}
	req, err := http.NewRequest("POST", "/gateway/api/compass/v1/metrics", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return ""
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-OK response: %v", resp.Status)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			return ""
		}
		log.Printf("Response body: %s", string(body))
		return ""
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return ""
	}

	return string(respBody)
}

func (c *CompassService) GetCompassCloudId() string {
	return c.cloudId
}
