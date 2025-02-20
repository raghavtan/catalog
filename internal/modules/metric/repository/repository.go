package repository

import (
	"context"
	"errors"
	"log"

	"github.com/motain/fact-collector/internal/modules/metric/resources"
	"github.com/motain/fact-collector/internal/services/compassservice"
)

type RepositoryInterface interface {
	CreateMetric(ctx context.Context, metric resources.Metric) (string, error)
	UpdateMetric(ctx context.Context, metric resources.Metric) error
	DeleteMetric(ctx context.Context, id string) error
}

type Repository struct {
	compass compassservice.CompassServiceInterface
}

func NewRepository(
	compass compassservice.CompassServiceInterface,
) *Repository {
	return &Repository{compass: compass}
}

func (r *Repository) CreateMetric(ctx context.Context, metric resources.Metric) (string, error) {
	query := `
		mutation createMetricDefinition ($cloudId: ID!, $name: String!, $description: String!, $unit: String!) {
			compass {
				createMetricDefinition(
					input: {
						cloudId: $cloudId
						name: $name
						description: $description
						format: {
							suffix: { suffix: $unit }
						}
					}
				) {
					success
					createdMetricDefinition {
						id
					}
					errors {
						message
					}
				}
			}
		}`

	variables := map[string]interface{}{
		"cloudId":     r.compass.GetCompassCloudId(),
		"name":        metric.Name,
		"description": metric.Description,
		"unit":        metric.Format.Unit,
	}

	var response struct {
		Compass struct {
			CreateMetricDefinition struct {
				Success                bool `json:"success"`
				CreateMetricDefinition struct {
					ID string `json:"id"`
				} `json:"createdMetricDefinition"`
			} `json:"createMetricDefinition"`
		} `json:"compass"`
	}

	if err := r.compass.Run(ctx, query, variables, &response); err != nil {
		log.Printf("Failed to create metric: %v", err)
		return "", err
	}

	return response.Compass.CreateMetricDefinition.CreateMetricDefinition.ID, nil
}

func (r *Repository) UpdateMetric(ctx context.Context, metric resources.Metric) error {
	query := `
		mutation updateMetricDefinition ($cloudId: ID!, $id: ID!, $name: String!, $description: String!, $unit: String!) {
			compass {
				updateMetricDefinition(
					input: {
						id: $id
						cloudId: $cloudId
						name: $name
						description: $description
						format: {
							suffix: { suffix: $unit }
						}
					}
				) {
					success
					errors {
						message
					}
				}
			}
		}`

	variables := map[string]interface{}{
		"cloudId":     r.compass.GetCompassCloudId(),
		"id":          *metric.ID,
		"name":        metric.Name,
		"description": metric.Description,
		"unit":        metric.Format.Unit,
	}

	var response struct {
		Compass struct {
			UpdateMetricDefinition struct {
				Success bool `json:"success"`
			} `json:"updateMetricDefinition"`
		} `json:"compass"`
	}

	if err := r.compass.Run(ctx, query, variables, &response); err != nil {
		log.Printf("Failed to update metric: %v", err)
		return err
	}

	if !response.Compass.UpdateMetricDefinition.Success {
		return errors.New("failed to update metric")
	}

	return nil
}

func (r *Repository) DeleteMetric(ctx context.Context, id string) error {
	query := `
		mutation deleteMetricDefinition($id: ID!) {
			compass {
				deleteMetricDefinition(input: {id: $id}) {
					deletedMetricDefinitionId
					errors {
						message
					}
					success
				}
			}
		}`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		Compass struct {
			DeleteMetricDefinition struct {
				Success bool `json:"success"`
			} `json:"deleteMetricDefinition"`
		} `json:"compass"`
	}

	if err := r.compass.Run(ctx, query, variables, &response); err != nil {
		log.Printf("Failed to create metric: %v", err)
		return err
	}

	return nil
}
