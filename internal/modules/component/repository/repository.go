package repository

//go:generate mockgen -destination=./mock_repository.go -package=repository github.com/motain/of-catalog/internal/modules/component/repository RepositoryInterface

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
)

type RepositoryInterface interface {
	Create(ctx context.Context, component resources.Component) (resources.Component, error)
	Update(ctx context.Context, component resources.Component) error
	Delete(ctx context.Context, id string) error
	GetBySlug(slug string) (*resources.Component, error)
	// MetricSource operations
	BindMetric(ctx context.Context, componentID string, metricID string, intentifier string) (string, error)
	UnBindMetric(ctx context.Context, metricSourceID string) error
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

func (r *Repository) Create(ctx context.Context, component resources.Component) (resources.Component, error) {
	query := `
		mutation createComponent ($cloudId: ID!, $componentDetails: CreateCompassComponentInput!) {
			compass {
				createComponent(cloudId: $cloudId, input: $componentDetails) {
					success
					componentDetails {
						id
					}
					errors {
						message
					}
				}
			}
		}`

	links := make([]map[string]string, 0)
	for _, link := range component.Links {
		links = append(links, map[string]string{
			"type": link.Type,
			"name": link.Name,
			"url":  link.URL,
		})
	}

	variables := map[string]interface{}{
		"cloudId": r.compass.GetCompassCloudId(),
		"componentDetails": map[string]interface{}{
			"name":        component.Name,
			"slug":        component.Slug,
			"description": component.Description,
			"typeId":      component.TypeID,
			"links":       links,
			"labels":      component.Labels,
		},
	}

	if component.OwnerID != "" {
		variables["componentDetails"].(map[string]interface{})["ownerId"] = component.OwnerID
	}

	var response struct {
		Compass struct {
			CreateComponent struct {
				Success          bool                          `json:"success"`
				Errors           []compassservice.CompassError `json:"errors"`
				ComponentDetails struct {
					ID            string `json:"id"`
					MetricSources struct {
						Nodes []struct {
							ID               string `json:"id"`
							MetricDefinition struct {
								ID   string `json:"id"`
								Name string `json:"name"`
							} `json:"metricDefinition"`
						} `json:"nodes"`
					} `json:"metricSources"`
				} `json:"componentDetails"`
			} `json:"createComponent"`
		} `json:"compass"`
	}

	if err := r.compass.Run(ctx, query, variables, &response); err != nil {
		log.Printf("Failed to create component: %v", err)
		return component, err
	}

	if compassservice.HasAlreadyExistsError(response.Compass.CreateComponent.Errors) {
		remoteComponent, err := r.GetBySlug(component.Slug)
		if err != nil {
			return component, err
		}

		component.ID = remoteComponent.ID
		component.MetricSources = remoteComponent.MetricSources
		updateError := r.Update(ctx, component)

		return component, updateError
	} else {
		if !response.Compass.CreateComponent.Success {
			return component, errors.New("failed to create component")
		}
	}

	metricSources := make(map[string]*resources.MetricSource)
	for _, node := range response.Compass.CreateComponent.ComponentDetails.MetricSources.Nodes {
		metricSources[node.MetricDefinition.Name] = &resources.MetricSource{
			ID:     node.ID,
			Metric: node.MetricDefinition.ID,
		}
	}

	component.ID = response.Compass.CreateComponent.ComponentDetails.ID
	component.MetricSources = metricSources
	return component, nil
}

func (r *Repository) Update(ctx context.Context, component resources.Component) error {
	query := `
		mutation updateComponent ($componentDetails: UpdateCompassComponentInput!) {
			compass {
				updateComponent(input: $componentDetails) {
					success
					errors {
						message
					}
				}
			}
		}`

	variables := map[string]interface{}{
		"cloudId": r.compass.GetCompassCloudId(),
		"componentDetails": map[string]interface{}{
			"id":          component.ID,
			"name":        component.Name,
			"slug":        component.Slug,
			"description": component.Description,
		},
	}

	if component.OwnerID != "" {
		variables["componentDetails"].(map[string]interface{})["ownerId"] = component.OwnerID
	}
	var response struct {
		Compass struct {
			UpdateComponentDefinition struct {
				Success bool `json:"success"`
			} `json:"updateComponent"`
		} `json:"compass"`
	}

	if err := r.compass.Run(ctx, query, variables, &response); err != nil {
		log.Printf("Failed to update component: %v", err)
		return err
	}

	if !response.Compass.UpdateComponentDefinition.Success {
		return errors.New("failed to update component")
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `
		mutation deleteComponent($id: ID!) {
			compass {
				deleteComponent(input: {id: $id}) {
					deletedComponentId
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
			DeleteComponent struct {
				Success bool `json:"success"`
			} `json:"deleteComponent"`
		} `json:"compass"`
	}

	if err := r.compass.Run(ctx, query, variables, &response); err != nil {
		log.Printf("Failed to create component: %v", err)
		return err
	}

	if !response.Compass.DeleteComponent.Success {
		return errors.New("failed to delete component")
	}

	return nil
}

func (r *Repository) GetBySlug(slug string) (*resources.Component, error) {
	query := `
		query getComponentBySlug($cloudId: ID!, $slug: String!) {
			compass {
				componentByReference(reference: {slug: {slug: $slug, cloudId: $cloudId}}) {
					... on CompassComponent {
						id
						metricSources {
							... on CompassComponentMetricSourcesConnection {
								nodes {
									id,
									metricDefinition {
										name
									}
								}
							}
						}
					}
				}
			}
		}`

	variables := map[string]interface{}{
		"cloudId": r.compass.GetCompassCloudId(),
		"slug":    slug,
	}

	var response struct {
		Compass struct {
			ComponentByReference struct {
				ID            string `json:"id"`
				MetricSources struct {
					Nodes []struct {
						ID               string `json:"id"`
						MetricDefinition struct {
							ID   string `json:"id"`
							Name string `json:"name"`
						} `json:"metricDefinition"`
					} `json:"nodes"`
				} `json:"metricSources"`
			} `json:"componentByReference"`
		} `json:"compass"`
	}

	if err := r.compass.Run(context.Background(), query, variables, &response); err != nil {
		log.Printf("Failed to get component by slug: %v", err)
		return nil, err
	}

	metricSources := make(map[string]*resources.MetricSource)
	for _, node := range response.Compass.ComponentByReference.MetricSources.Nodes {
		metricSources[node.MetricDefinition.Name] = &resources.MetricSource{
			ID:     node.ID,
			Metric: node.MetricDefinition.ID,
		}
	}
	component := resources.Component{
		ID:            response.Compass.ComponentByReference.ID,
		MetricSources: metricSources,
	}

	return &component, nil
}

func (r *Repository) BindMetric(ctx context.Context, componentID string, metricID string, intentifier string) (string, error) {
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

	if !response.Compass.CreateMetricSource.Success {
		return "", errors.New("failed to create metric source")
	}

	return response.Compass.CreateMetricSource.CreateMetricSource.ID, nil
}

func (r *Repository) UnBindMetric(ctx context.Context, metricSourceID string) error {
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

	if !response.Compass.DeleteMetricSource.Success {
		return errors.New("failed to delete metric source")
	}

	return nil
}

func (r *Repository) Push(ctx context.Context, metricSourceID string, value float64, recordedAt time.Time) error {
	requestBody := map[string]string{
		"metricSourceId": metricSourceID,
		"value":          fmt.Sprintf("%f", value),
		"timestamp":      recordedAt.UTC().Format(time.RFC3339),
	}

	_, errSend := r.compass.SendMetric(requestBody)

	return errSend
}
