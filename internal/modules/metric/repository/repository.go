package repository

//go:generate mockgen -destination=./mocks/mock_repository.go -package=repository github.com/motain/of-catalog/internal/modules/metric/repository RepositoryInterface

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/motain/of-catalog/internal/interfaces/repositoryinterfaces"
	"github.com/motain/of-catalog/internal/modules/metric/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/metric/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
)

type RepositoryInterface interface {
	Create(ctx context.Context, metric resources.Metric) (string, error)
	Update(ctx context.Context, metric resources.Metric) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, metric resources.Metric) (*resources.Metric, error)
}

type Repository struct {
	compass compassservice.CompassServiceInterface
}

func NewRepository(compass compassservice.CompassServiceInterface) *Repository {
	return &Repository{compass: compass}
}

func (r *Repository) Create(ctx context.Context, metric resources.Metric) (string, error) {
	input := &dtos.CreateMetricInput{Metric: metric}
	output := &dtos.CreateMetricOutput{}

	// This function is executed before the validation of the operation
	// That is before checking if the operation was successful
	// If the metric already exists, it searches for the metric in the remote service
	// and updates the local metric with the remote metric ID
	// and clears the errors
	preValidationFunc := func() error {
		if !compassservice.HasAlreadyExistsError(output.Compass.CreateMetric.Errors) {
			return nil
		}

		remoteMetric, err := r.Search(ctx, metric)
		if err != nil {
			return err
		}

		metric.ID = remoteMetric.ID
		updateError := r.Update(ctx, metric)
		if updateError != nil {
			return updateError
		}

		output.Compass.CreateMetric.Definition.ID = remoteMetric.ID
		output.Compass.CreateMetric.Errors = nil
		output.Compass.CreateMetric.Success = true

		return nil
	}

	runErr := r.run(ctx, input, output, preValidationFunc)
	if runErr != nil {
		return "", fmt.Errorf("Create error for %s: %s", metric, runErr)
	}

	return output.Compass.CreateMetric.Definition.ID, nil
}

func (r *Repository) Update(ctx context.Context, metric resources.Metric) error {
	input := &dtos.UpdateMetricInput{CompassCloudID: r.compass.GetCompassCloudId(), Metric: metric}
	output := &dtos.UpdateMetricOutput{}
	if runErr := r.run(ctx, input, output, nil); runErr != nil {
		return fmt.Errorf("Update error for %s: %s", metric, runErr)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	input := &dtos.DeleteMetricInput{MetricID: id}
	output := &dtos.DeleteMetricOutput{}
	if runErr := r.run(ctx, input, output, nil); runErr != nil {
		return fmt.Errorf("Delete error for %s: %s", id, runErr)
	}
	return nil
}

func (r *Repository) Search(ctx context.Context, metric resources.Metric) (*resources.Metric, error) {
	input := &dtos.SearchMetricsInput{Metric: metric}
	output := &dtos.SearchMetricsOutput{}
	runErr := r.run(ctx, input, output, nil)
	if runErr != nil {
		return nil, fmt.Errorf("Search error for %s: %s", metric.Name, runErr)
	}

	for _, node := range output.Compass.Definitions.Nodes {
		if node.Name == metric.Name {
			return &resources.Metric{ID: node.ID}, nil
		}
	}

	return nil, fmt.Errorf("Search error for %s: %s", metric.Name, "metric not found")
}

func (r *Repository) run(
	ctx context.Context,
	input repositoryinterfaces.InputDTOInterface,
	output repositoryinterfaces.OutputDTOInterface,
	preValidationFunc repositoryinterfaces.ValidationFunc,
) error {
	query := input.GetQuery()
	operation := strings.TrimSpace(query[:strings.Index(query, "(")])

	if runErr := r.compass.Run(ctx, query, input.SetVariables(), output); runErr != nil {
		log.Printf("failed to run %s: %v", operation, runErr)
		return runErr
	}

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
