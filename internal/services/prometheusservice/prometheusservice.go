package prometheusservice

//go:generate mockgen -destination=./mocks/mock_prometheus_service.go -package=prometheusservice github.com/motain/of-catalog/internal/services/prometheusservice PrometheusServiceInterface

import (
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// PrometheusServiceInterface defines the contract for Prometheus service operations.
// It provides methods for executing both instant and range queries against Prometheus.
type PrometheusServiceInterface interface {
	// InstantQuery executes a PromQL query at the current time.
	// Returns the query result as a Prometheus model.Value.
	InstantQuery(queryString string) (float64, error)

	// RangeQuery executes a PromQL query over a specified time range.
	// Returns the query result as a Prometheus model.Value.
	RangeQuery(queryString string, start, end time.Time, step time.Duration) (model.Value, error)
}

// PrometheusService implements PrometheusServiceInterface to provide a high-level interface
// for interacting with Prometheus. It wraps the underlying PrometheusClientInterface to
// provide simplified query methods.
type PrometheusService struct {
	// client is the underlying Prometheus client that handles the actual API communication
	client PrometheusClientInterface
}

// NewPrometheusService creates a new PrometheusService instance with the provided client.
// This constructor ensures proper initialization of the service with its dependencies.
//
// Parameters:
//   - client: The PrometheusClientInterface implementation to use for API communication
//
// Returns:
//   - *PrometheusService: A new service instance ready to handle Prometheus queries
func NewPrometheusService(client PrometheusClientInterface) *PrometheusService {
	return &PrometheusService{client: client}
}

// InstantQuery executes a PromQL query at the current time.
// This method provides a simplified interface for executing instant queries
// without requiring the caller to specify a timestamp.
//
// Parameters:
//   - queryString: The PromQL query to execute
//
// Returns:
//   - model.Value: The query result in Prometheus model format
//   - error: Any error that occurred during query execution
func (ps *PrometheusService) InstantQuery(queryString string) (float64, error) {
	return ps.client.Query(queryString, time.Now())
}

// RangeQuery executes a PromQL query over a specified time range.
// This method provides a simplified interface for executing range queries
// by accepting start time, end time, and step duration as separate parameters.
//
// Parameters:
//   - queryString: The PromQL query to execute
//   - start: The start time of the query range
//   - end: The end time of the query range
//   - step: The time step between data points
//
// Returns:
//   - model.Value: The query result in Prometheus model format
//   - error: Any error that occurred during query execution
func (ps *PrometheusService) RangeQuery(queryString string, start, end time.Time, step time.Duration) (model.Value, error) {
	r := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}
	return ps.client.QueryRange(queryString, r)
}
