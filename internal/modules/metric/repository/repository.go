package repository

//go:generate mockgen -destination=./mock_repository.go -package=repository github.com/motain/of-catalog/internal/modules/metric/repository RepositoryInterface

import (
	"context"
	"errors"
	"fmt"
	"log"

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
				Success                bool                          `json:"success"`
				Errors                 []compassservice.CompassError `json:"errors"`
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

	for _, err := range response.Compass.CreateMetricDefinition.Errors {
		fmt.Printf("Error: %v\n", err.Message)
	}
	if compassservice.HasAlreadyExistsError(response.Compass.CreateMetricDefinition.Errors) {
		remoteMetric, err := r.Search(metric)
		if err != nil {
			return "", err
		}

		metric.ID = remoteMetric.ID
		updateError := r.Update(ctx, metric)

		return remoteMetric.ID, updateError
	} else {
		if !response.Compass.CreateMetricDefinition.Success {
			return "", errors.New("failed to create metric")
		}
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
		"id":          metric.ID,
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

	if !response.Compass.DeleteMetricDefinition.Success {
		return errors.New("failed to delete metric")
	}

	return nil
}

func (r *Repository) Search(metric resources.Metric) (*resources.Metric, error) {
	query := `
		query searchMetricDefinition($cloudId: ID!) {
			compass {
				metricDefinitions(query: {cloudId: $cloudId, first: 100}) {
					... on CompassMetricDefinitionsConnection {
						nodes{
							id
							name
						}
					}
				}
			}
		}`

	variables := map[string]interface{}{
		"cloudId": r.compass.GetCompassCloudId(),
		"name":    metric.Name,
	}

	var response struct {
		Compass struct {
			MetricDefinitions struct {
				Nodes []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"nodes"`
			} `json:"metricDefinitions"`
		} `json:"compass"`
	}

	if err := r.compass.Run(context.Background(), query, variables, &response); err != nil {
		log.Printf("Failed to search metric: %v", err)
		return nil, err
	}

	if len(response.Compass.MetricDefinitions.Nodes) == 0 {
		return nil, errors.New("metric not found")
	}

	for _, node := range response.Compass.MetricDefinitions.Nodes {
		if node.Name == metric.Name {
			return &resources.Metric{ID: node.ID}, nil
		}
	}

	return nil, errors.New("metric not found")
}
