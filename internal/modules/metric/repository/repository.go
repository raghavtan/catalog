package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/motain/fact-collector/internal/modules/metric/resources"
	"github.com/motain/fact-collector/internal/services/compassservice"
)

type RepositoryInterface interface {
	Create(ctx context.Context, metric resources.Metric) (string, error)
	Update(ctx context.Context, metric resources.Metric) error
	Delete(ctx context.Context, id string) error
	CreateMetricSource(ctx context.Context, metricID string, componentID string, intentifier string) (string, error)
	DeleteMetricSource(ctx context.Context, metricSourceID string) error
	Push(ctx context.Context, metricSourceID string, value float64, recordedAt time.Time) error
}

type Repository struct {
	compass compassservice.CompassServiceInterface
}

func NewRepository(
	compass compassservice.CompassServiceInterface,
) *Repository {
	return &Repository{compass: compass}
}

func (r *Repository) Create(ctx context.Context, metric resources.Metric) (string, error) {
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

func (r *Repository) Update(ctx context.Context, metric resources.Metric) error {
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

func (r *Repository) Delete(ctx context.Context, id string) error {
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

func (r *Repository) CreateMetricSource(ctx context.Context, metricID string, componentID string, intentifier string) (string, error) {
	query := `
		mutation createMetricSource($metricId: ID!, $componentId: ID!, $externalId: ID!) {
			compass {
				createMetricSource(input: {metricDefinitionId: $metricId, componentId: $componentId, externalMetricSourceId: $externalId}) {
					success
					createdMetricSource {
						id
					}
					errors {
						message
					}
				}
			}
		}`

	variables := map[string]interface{}{
		"metricId":    metricID,
		"componentId": componentID,
		"externalId":  intentifier,
	}

	var response struct {
		Compass struct {
			CreateMetricSource struct {
				Success            bool `json:"success"`
				CreateMetricSource struct {
					ID string `json:"id"`
				} `json:"createdMetricSource"`
			} `json:"createMetricSource"`
		} `json:"compass"`
	}

	if err := r.compass.Run(ctx, query, variables, &response); err != nil {
		log.Printf("Failed to create metric source: %v", err)
		return "", err
	}

	return response.Compass.CreateMetricSource.CreateMetricSource.ID, nil
}

func (r *Repository) DeleteMetricSource(ctx context.Context, metricSourceID string) error {
	query := `
		mutation deleteMetricSource($id: ID!) {
			compass {
				deleteMetricSource(input: {id: $id}) {
					deletedMetricSourceId
					errors {
						message
					}
					success
				}
			}
		}`

	variables := map[string]interface{}{
		"id": metricSourceID,
	}

	var response struct {
		Compass struct {
			DeleteMetricSource struct {
				Success bool `json:"success"`
			} `json:"deleteMetricSource"`
		} `json:"compass"`
	}

	if err := r.compass.Run(ctx, query, variables, &response); err != nil {
		log.Printf("Failed to delete metric source: %v", err)
		return err
	}

	return nil
}

func (r *Repository) Push(ctx context.Context, metricSourceID string, value float64, recordedAt time.Time) error {
	requestBody := map[string]string{
		"metricSourceId": metricSourceID,
		"value":          fmt.Sprintf("%f", value),
		"timestamp":      recordedAt.UTC().Format(time.RFC3339),
	}

	r.compass.SendMetric(requestBody)

	return nil
}
