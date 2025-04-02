package compassservice

//go:generate mockgen -destination=./mocks/mock_compass_service.go -package=compassservice github.com/motain/of-catalog/internal/services/compassservice CompassServiceInterface

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/machinebox/graphql"
	"github.com/motain/of-catalog/internal/services/compassservice/dtos"
	"github.com/motain/of-catalog/internal/services/configservice"
)

type InputDTOInterface interface {
	GetQuery() string
	GetPreValidationFunc() dtos.ValidationFunc
	SetVariables() map[string]interface{}
}

type OutputDTOInterface interface {
	IsSuccessful() bool
	GetErrors() []string
}

type CompassServiceInterface interface {
	Run(ctx context.Context, query string, variables map[string]interface{}, response interface{}) error
	RunWithDTOs(ctx context.Context, input InputDTOInterface, output OutputDTOInterface) error
	SendMetric(ctx context.Context, body map[string]string) (string, error)
	SendAPISpecifications(ctx context.Context, input dtos.APISpecificationsInput) (string, error)
	GetCompassCloudId() string
}

const (
	metricsV1Endpoint  = "/gateway/api/compass/v1/metrics"
	apiSpecsV1Endpoint = "/gateway/api/compass/v1/component/:componentId/api_specs"
)

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

func (c *CompassService) RunWithDTOs(ctx context.Context, input InputDTOInterface, output OutputDTOInterface) error {
	query := input.GetQuery()
	operation := strings.TrimSpace(query[:strings.Index(query, "(")])

	if runErr := c.Run(ctx, query, input.SetVariables(), output); runErr != nil {
		log.Printf("failed to run %s: %v", operation, runErr)
		return runErr
	}

	preValidationFunc := input.GetPreValidationFunc()
	if preValidationFunc != nil {
		if runErr := preValidationFunc(); runErr != nil {
			return runErr
		}
	}

	if !output.IsSuccessful() {
		return fmt.Errorf("failed to execute %s: %v", operation, output.GetErrors())
	}

	return nil
}

func (c *CompassService) SendMetric(ctx context.Context, body map[string]string) (string, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, metricsV1Endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return c.do(req)
}

func (c *CompassService) SendAPISpecifications(ctx context.Context, input dtos.APISpecificationsInput) (string, error) {
	endpoint := strings.Replace(apiSpecsV1Endpoint, ":componentId", input.ComponentID, 1)

	body, contentType, buildBodyErr := c.buildMultiPartBody(input)
	if buildBodyErr != nil {
		return "", fmt.Errorf("failed to build multipart body: %w", buildBodyErr)
	}

	req, requestErr := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, body)
	if requestErr != nil {
		return "", fmt.Errorf("failed to create request: %w", requestErr)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")

	return c.do(req)
}

func (c *CompassService) GetCompassCloudId() string {
	return c.cloudId
}

func (c *CompassService) do(req *http.Request) (string, error) {
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

func (c *CompassService) buildMultiPartBody(input dtos.APISpecificationsInput) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, createFormFileErr := writer.CreateFormFile("file", input.FileName)
	if createFormFileErr != nil {
		return body, "", fmt.Errorf("failed to create form file: %w", createFormFileErr)
	}

	_, writeErr := part.Write([]byte(input.ApiSpecs))
	if writeErr != nil {
		return body, "", fmt.Errorf("failed to write to form file: %w", writeErr)
	}

	closeWriterErr := writer.Close()
	if closeWriterErr != nil {
		return body, "", fmt.Errorf("failed to close writer: %w", closeWriterErr)
	}

	return body, writer.FormDataContentType(), nil
}
