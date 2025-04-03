package prometheusservice

//go:generate mockgen -destination=./mocks/mock_prometheus_client.go -package=prometheusservice github.com/motain/of-catalog/internal/services/prometheusservice PrometheusClientInterface

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"

	"github.com/motain/of-catalog/internal/services/configservice"
	"github.com/motain/of-catalog/internal/utils/awsutils"
)

// PrometheusClientInterface defines the contract for Prometheus client operations.
// It provides methods for querying Prometheus data in different formats.
type PrometheusClientInterface interface {
	// Query executes an instant query at a specific timestamp
	Query(query string, timestamp time.Time) (float64, error)
	// QueryRange executes a query over a time range
	QueryRange(query string, r v1.Range) (model.Value, error)
}

// PrometheusClient implements PrometheusClientInterface to interact with AWS Managed Prometheus.
// It handles authentication and provides methods to query Prometheus data.
type PrometheusClient struct {
	api    v1.API // Underlying Prometheus API client
	region string // AWS region for the Prometheus workspace
}

// NewPrometheusClient creates a new Prometheus client with AWS authentication.
// It configures the client with either existing AWS credentials or assumes a specified role.
//
// Parameters:
//   - cfg: Configuration service providing AWS region, Prometheus URL, and role ARN
//
// Returns:
//   - PrometheusClientInterface: Configured Prometheus client
//   - Panics if configuration is invalid or client creation fails
func NewPrometheusClient(cfg configservice.ConfigServiceInterface) PrometheusClientInterface {
	ctx := context.Background()
	region := cfg.GetAWSRegion()
	workspaceURL := cfg.GetPrometheusURL()

	if workspaceURL == "" {
		panic("Prometheus workspace URL not configured")
	}

	// Initialize AWS configuration
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		panic(fmt.Errorf("failed to load AWS config: %w", err))
	}

	// Set up credentials provider
	credProvider := getCredentialsProvider(ctx, awsCfg, cfg.GetAWSRole())

	// Create authenticated HTTP client
	httpClient := &http.Client{
		Transport: &awsutils.SigV4RoundTripper{
			Transport:   http.DefaultTransport,
			Region:      region,
			Service:     "aps",
			Credentials: credProvider,
		},
	}

	// Initialize Prometheus client
	promClient, err := api.NewClient(api.Config{
		Address: workspaceURL,
		Client:  httpClient,
	})
	if err != nil {
		panic(fmt.Errorf("failed to create Prometheus client: %w", err))
	}

	return &PrometheusClient{
		api:    v1.NewAPI(promClient),
		region: region,
	}
}

// getCredentialsProvider determines the appropriate AWS credentials provider.
// It either uses existing credentials or assumes a specified role.
func getCredentialsProvider(ctx context.Context, awsCfg aws.Config, roleARN string) aws.CredentialsProvider {
	if roleARN == "" {
		stsClient := sts.NewFromConfig(awsCfg)
		_, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			panic(fmt.Errorf("failed to get caller identity: %w", err))
		}
		return awsCfg.Credentials
	}
	return stscreds.NewAssumeRoleProvider(sts.NewFromConfig(awsCfg), roleARN)
}

// Query executes an instant query against Prometheus at the specified timestamp.
// Returns the query result as a Prometheus model.Value.
func (pc *PrometheusClient) Query(query string, timestamp time.Time) (float64, error) {
	result, _, err := pc.api.Query(context.Background(), query, timestamp)
	if err != nil {
		fmt.Printf("Query: %s, Timestamp: %s\n", query, timestamp)
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}
	response := 0.0
	if vector, ok := result.(model.Vector); ok {
		for _, sample := range vector {
			response = float64(sample.Value)
		}
	} else {
		fmt.Println("Result is not a vector for query:", query)
	}

	return response, nil
}

// QueryRange executes a range query against Prometheus over the specified time range.
// Returns the query result as a Prometheus model.Value.
func (pc *PrometheusClient) QueryRange(query string, r v1.Range) (model.Value, error) {
	result, _, err := pc.api.QueryRange(context.Background(), query, r)
	if err != nil {
		return nil, fmt.Errorf("failed to execute range query: %w", err)
	}
	return result, nil
}
