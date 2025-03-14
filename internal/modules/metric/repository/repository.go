package repository

//go:generate mockgen -destination=./mock_repository.go -package=repository github.com/motain/of-catalog/internal/modules/metric/repository RepositoryInterface

import (
	"context"
	"errors"
	"log"

	"github.com/motain/of-catalog/internal/modules/metric/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/metric/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
)

type RepositoryInterface interface {
	Create(ctx context.Context, metric resources.Metric) (string, error)
	Update(ctx context.Context, metric resources.Metric) error
	Delete(ctx context.Context, id string) error
}

type Repository struct {
	compass compassservice.CompassServiceInterface
}

func NewRepository(compass compassservice.CompassServiceInterface) *Repository {
	return &Repository{compass: compass}
}

func (r *Repository) Create(ctx context.Context, metric resources.Metric) (string, error) {
	metricDto := dtos.CreateMetricOutput{}
	query := metricDto.GetQuery()
	variables := metricDto.SetVariables(r.compass.GetCompassCloudId(), metric)

	if err := r.compass.Run(ctx, query, variables, &metricDto); err != nil {
		log.Printf("Failed to create metric: %v", err)
		return "", err
	}

	if compassservice.HasAlreadyExistsError(metricDto.Compass.CreateMetric.Errors) {
		remoteMetric, err := r.Search(metric)
		if err != nil {
			return "", err
		}

		metric.ID = remoteMetric.ID
		updateError := r.Update(ctx, metric)

		return remoteMetric.ID, updateError
	}

	if !metricDto.IsSuccessful() {
		return "", errors.New("failed to create metric")
	}

	return metricDto.Compass.CreateMetric.Definition.ID, nil
}

func (r *Repository) Update(ctx context.Context, metric resources.Metric) error {
	metricDto := dtos.UpdateMetricOutput{}
	query := metricDto.GetQuery()
	variables := metricDto.SetVariables(r.compass.GetCompassCloudId(), metric)

	if err := r.compass.Run(ctx, query, variables, &metricDto); err != nil {
		log.Printf("Failed to update metric: %v", err)
		return err
	}

	if !metricDto.IsSuccessful() {
		return errors.New("failed to update metric")
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	metricDto := dtos.DeleteMetricOutput{}
	query := metricDto.GetQuery()
	variables := metricDto.SetVariables(id)

	if err := r.compass.Run(ctx, query, variables, &metricDto); err != nil {
		log.Printf("Failed to create metric: %v", err)
		return err
	}

	if !metricDto.IsSuccessful() {
		return errors.New("failed to delete metric")
	}

	return nil
}

func (r *Repository) Search(metric resources.Metric) (*resources.Metric, error) {
	metricsDto := dtos.SearchMetricsOutput{}
	query := metricsDto.GetQuery()
	variables := metricsDto.SetVariables(r.compass.GetCompassCloudId(), metric)

	if err := r.compass.Run(context.Background(), query, variables, &metricsDto); err != nil {
		log.Printf("Failed to search metric: %v", err)
		return nil, err
	}

	if !metricsDto.IsSuccessful() {
		return nil, errors.New("metric not found")
	}

	for _, node := range metricsDto.Compass.Definitions.Nodes {
		if node.Name == metric.Name {
			return &resources.Metric{ID: node.ID}, nil
		}
	}

	return nil, errors.New("metric not found")
}
