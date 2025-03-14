package repository

//go:generate mockgen -destination=./mock_repository.go -package=repository github.com/motain/of-catalog/internal/modules/component/repository RepositoryInterface

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
)

type RepositoryInterface interface {
	Create(ctx context.Context, component resources.Component) (resources.Component, error)
	Update(ctx context.Context, component resources.Component) error
	Delete(ctx context.Context, id string) error
	GetBySlug(slug string) (*resources.Component, error)
	// Dependency operations
	SetDependency(ctx context.Context, dependentId, providerId string) error
	UnSetDependency(ctx context.Context, dependentId, providerId string) error
	// Documents operations
	AddDocument(ctx context.Context, componentID string, document resources.Document) (resources.Document, error)
	UpdateDocument(ctx context.Context, componentID string, document resources.Document) error
	RemoveDocument(ctx context.Context, componentID, documentID string) error
	// MetricSource operations
	BindMetric(ctx context.Context, componentID string, metricID string, intentifier string) (string, error)
	UnBindMetric(ctx context.Context, metricSourceID string) error
	Push(ctx context.Context, metricSourceID string, value float64, recordedAt time.Time) error
}

type Repository struct {
	compass            compassservice.CompassServiceInterface
	DocumentCategories map[string]string
}

func NewRepository(
	compass compassservice.CompassServiceInterface,
) *Repository {
	return &Repository{compass: compass, DocumentCategories: nil}
}

func (r *Repository) Create(ctx context.Context, component resources.Component) (resources.Component, error) {
	componentDTO := dtos.CreateComponentOutput{}
	query := componentDTO.GetQuery()
	variables := componentDTO.SetVariables(r.compass.GetCompassCloudId(), component)

	if err := r.compass.Run(ctx, query, variables, &componentDTO); err != nil {
		log.Printf("Failed to create component: %v", err)
		return component, err
	}

	if compassservice.HasAlreadyExistsError(componentDTO.Compass.CreateComponent.Errors) {
		remoteComponent, err := r.GetBySlug(component.Slug)
		if err != nil {
			return component, err
		}

		component.ID = remoteComponent.ID
		component.MetricSources = remoteComponent.MetricSources
		updateError := r.Update(ctx, component)

		return component, updateError
	}

	if !componentDTO.IsSuccessful() {
		return component, errors.New("failed to create component")
	}

	metricSources := make(map[string]*resources.MetricSource)
	for _, node := range componentDTO.Compass.CreateComponent.Details.MetricSources.Nodes {
		metricSources[node.MetricDefinition.Name] = &resources.MetricSource{
			ID:     node.ID,
			Metric: node.MetricDefinition.ID,
		}
	}

	createdLinks := make([]resources.Link, len(componentDTO.Compass.CreateComponent.Details.Links))
	for i, link := range componentDTO.Compass.CreateComponent.Details.Links {
		createdLinks[i] = resources.Link{
			ID:   link.ID,
			Type: link.Type,
			Name: link.Name,
			URL:  link.URL,
		}
	}
	component.ID = componentDTO.Compass.CreateComponent.Details.ID
	component.MetricSources = metricSources
	component.Links = createdLinks

	return component, nil
}

func (r *Repository) Update(ctx context.Context, component resources.Component) error {
	componentDTO := dtos.UpdateComponentOutput{}
	query := componentDTO.GetQuery()
	variables := componentDTO.SetVariables(component)

	if err := r.compass.Run(ctx, query, variables, &componentDTO); err != nil {
		log.Printf("Failed to update component: %v", err)
		return err
	}

	if !componentDTO.IsSuccessful() {
		return errors.New("failed to update component")
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	componentDTO := dtos.DeleteComponentOutput{}
	query := componentDTO.GetQuery()
	variables := componentDTO.SetVariables(id)

	if err := r.compass.Run(ctx, query, variables, &componentDTO); err != nil {
		log.Printf("Failed to delete component: %v", err)
		return err
	}

	if compassservice.HasNotFoundError(componentDTO.Compass.DeleteComponent.Errors) {
		return nil
	}

	if !componentDTO.IsSuccessful() {
		return errors.New("failed to delete component")
	}

	return nil
}

func (r *Repository) SetDependency(ctx context.Context, dependentId, providerId string) error {
	componentDTO := dtos.CreateDependencyOutput{}
	query := componentDTO.GetQuery()
	variables := componentDTO.SetVariables(dependentId, providerId)

	if err := r.compass.Run(ctx, query, variables, &componentDTO); err != nil {
		log.Printf("failed to create component dependency: %v", err)
		return err
	}

	if !componentDTO.IsSuccessful() {
		return errors.New("failed to create component dependency")
	}

	return nil
}

func (r *Repository) UnSetDependency(ctx context.Context, dependentId, providerId string) error {
	componentDTO := dtos.DeleteDependencyOutput{}
	query := componentDTO.GetQuery()
	variables := componentDTO.SetVariables(dependentId, providerId)

	if err := r.compass.Run(ctx, query, variables, &componentDTO); err != nil {
		log.Printf("Failed to delete component dependency: %v", err)
		return err
	}

	if !componentDTO.IsSuccessful() {
		return errors.New("failed to delete component dependency")
	}

	return nil
}

func (r *Repository) GetBySlug(slug string) (*resources.Component, error) {
	componentDTO := dtos.ComponentByReferenceOutput{}
	query := componentDTO.GetQuery()
	variables := componentDTO.SetVariables(r.compass.GetCompassCloudId(), slug)

	if err := r.compass.Run(context.Background(), query, variables, &componentDTO); err != nil {
		log.Printf("failed to get component by slug: %v", err)
		return nil, err
	}

	if !componentDTO.IsSuccessful() {
		return nil, errors.New("failed to get component by slug")
	}

	metricSources := make(map[string]*resources.MetricSource)
	for _, node := range componentDTO.Compass.Component.MetricSources.Nodes {
		metricSources[node.MetricDefinition.Name] = &resources.MetricSource{
			ID:     node.ID,
			Metric: node.MetricDefinition.ID,
		}
	}
	component := resources.Component{
		ID:            componentDTO.Compass.Component.ID,
		MetricSources: metricSources,
	}

	return &component, nil
}

func (r *Repository) AddDocument(ctx context.Context, componentID string, document resources.Document) (resources.Document, error) {
	r.initDocumentCategories(ctx)

	documentDTO := dtos.CreateDocumentOutput{}
	query := documentDTO.GetQuery()
	variables := documentDTO.SetVariables(
		componentID,
		document.Title,
		r.DocumentCategories[document.Type],
		document.URL,
	)

	if err := r.compass.Run(ctx, query, variables, &documentDTO); err != nil {
		log.Printf("Failed to create document: %v", err)
		return resources.Document{}, err
	}

	if !documentDTO.IsSuccessful() {
		return resources.Document{}, errors.New("failed to create document")
	}

	return resources.Document{
		ID:                      documentDTO.Compass.AddDocument.Details.ID,
		Title:                   document.Title,
		Type:                    document.Type,
		URL:                     document.URL,
		DocumentationCategoryId: r.DocumentCategories[document.Type],
	}, nil
}

func (r *Repository) UpdateDocument(ctx context.Context, componentID string, document resources.Document) error {
	r.initDocumentCategories(ctx)

	documentDTO := dtos.UpdateDocumentOutput{}
	query := documentDTO.GetQuery()
	variables := documentDTO.SetVariables(
		document.ID,
		document.Title,
		r.DocumentCategories[document.Type],
		document.URL,
	)

	if err := r.compass.Run(ctx, query, variables, &documentDTO); err != nil {
		log.Printf("Failed to update document: %v", err)
		return err
	}

	if !documentDTO.IsSuccessful() {
		return errors.New("failed to update document")
	}

	return nil
}

func (r *Repository) RemoveDocument(ctx context.Context, componentID, docuemntID string) error {
	query := `
		mutation deleteComponentLink($id: ID!) {
			compass {
				deleteComponentLink(input: {id: $id}) {
					deletedMetricSourceId
					errors {
						message
					}
					success
				}
			}
		}`

	variables := map[string]interface{}{
		"componentId": componentID,
		"id":          docuemntID,
	}

	var response struct {
		Compass struct {
			DeleteComponentLink struct {
				Success bool `json:"success"`
			} `json:"deleteComponentLink"`
		} `json:"compass"`
	}

	if err := r.compass.Run(ctx, query, variables, &response); err != nil {
		log.Printf("Failed to delete metric source: %v", err)
		return err
	}

	if !response.Compass.DeleteComponentLink.Success {
		return errors.New("failed to delete metric source")
	}

	return nil
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

func (r *Repository) initDocumentCategories(ctx context.Context) error {
	fmt.Printf("Category ID: %s\n", r.DocumentCategories)
	if r.DocumentCategories != nil {
		return nil
	}

	query := `
		query documentationCategories {
			compass {
				documentationCategories(cloudId: "fca6a80f-888b-4079-82e6-3c2f61c788e2") @optIn(to: "compass-beta")  {
					... on CompassDocumentationCategoriesConnection {
						nodes {
							name
							id
							description
						}
					}
				}
			}
		}`

	var response struct {
		Compass struct {
			DocumentationCategories struct {
				Nodes []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"nodes"`
			} `json:"documentationCategories"`
		} `json:"compass"`
	}

	if err := r.compass.Run(ctx, query, nil, &response); err != nil {
		log.Printf("Failed to delete metric source: %v", err)
		return err
	}

	categories := make(map[string]string, len(response.Compass.DocumentationCategories.Nodes))
	for _, category := range response.Compass.DocumentationCategories.Nodes {
		categories[category.Name] = category.ID
	}
	r.DocumentCategories = categories

	return nil
}
