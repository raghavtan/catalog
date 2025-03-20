package repository

//go:generate mockgen -destination=./mocks/mock_repository.go -package=repository github.com/motain/of-catalog/internal/modules/component/repository RepositoryInterface

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

type RepositoryInterface interface {
	Create(ctx context.Context, component resources.Component) (resources.Component, error)
	Update(ctx context.Context, component resources.Component) (resources.Component, error)
	Delete(ctx context.Context, id string) error
	GetBySlug(ctx context.Context, slug string) (*resources.Component, error)
	// Dependency operations
	SetDependency(ctx context.Context, dependentId, providerId string) error
	UnsetDependency(ctx context.Context, dependentId, providerId string) error
	// Documents operations
	AddDocument(ctx context.Context, componentID string, document resources.Document) (resources.Document, error)
	UpdateDocument(ctx context.Context, componentID string, document resources.Document) error
	RemoveDocument(ctx context.Context, componentID, documentID string) error
	// MetricSource operations
	BindMetric(ctx context.Context, componentID string, metricID string, identifier string) (string, error)
	UnbindMetric(ctx context.Context, metricSourceID string) error
	// API Specications operations
	SetAPISpecifications(ctx context.Context, componentID, apiSpecs, apiSpecsFile string) error
	// Push metric value
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
		remoteComponent, err := r.GetBySlug(ctx, component.Slug)
		if err != nil {
			return component, err
		}

		component.ID = remoteComponent.ID
		component.MetricSources = remoteComponent.MetricSources
		component, updateError := r.Update(ctx, component)

		return component, updateError
	}

	if !componentDTO.IsSuccessful() {
		return component, fmt.Errorf("failed to create component: %v", componentDTO.Compass.CreateComponent.Errors)
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

func (r *Repository) Update(ctx context.Context, component resources.Component) (resources.Component, error) {
	componentDTO := dtos.UpdateComponentOutput{}
	query := componentDTO.GetQuery()
	variables := componentDTO.SetVariables(component)

	if updateErr := r.compass.Run(ctx, query, variables, &componentDTO); updateErr != nil {
		return resources.Component{}, fmt.Errorf("failed to update component %s: %v", component.ID, updateErr)
	}

	if compassservice.HasNotFoundError(componentDTO.Compass.UpdateComponent.Errors) {
		remoteComponent, getBySlugErr := r.GetBySlug(ctx, component.Slug)
		if getBySlugErr != nil {
			return resources.Component{}, getBySlugErr
		}

		component.ID = remoteComponent.ID
		component.MetricSources = remoteComponent.MetricSources
		updatedComponent, updateError := r.Update(ctx, component)

		return updatedComponent, updateError
	}

	if !componentDTO.IsSuccessful() {
		return resources.Component{}, fmt.Errorf("failed to update component %s: %v", component.ID, componentDTO.Compass.UpdateComponent.Errors)
	}

	return component, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	componentDTO := dtos.DeleteComponentOutput{}
	query := componentDTO.GetQuery()
	variables := componentDTO.SetVariables(id)

	err := r.run(ctx, query, variables, &componentDTO, componentDTO.IsSuccessful)
	if err != nil {
		return fmt.Errorf("failed to delete component %s: %v", id, err)
	}

	return nil
}

func (r *Repository) SetDependency(ctx context.Context, dependentId, providerId string) error {
	componentDTO := dtos.CreateDependencyOutput{}
	query := componentDTO.GetQuery()
	variables := componentDTO.SetVariables(dependentId, providerId)

	err := r.run(ctx, query, variables, &componentDTO, componentDTO.IsSuccessful)
	if err != nil {
		return fmt.Errorf("failed to set component dependency %s -> %s: %v", dependentId, providerId, err)
	}

	return nil
}

func (r *Repository) UnsetDependency(ctx context.Context, dependentId, providerId string) error {
	componentDTO := dtos.DeleteDependencyOutput{}
	query := componentDTO.GetQuery()
	variables := componentDTO.SetVariables(dependentId, providerId)

	err := r.run(ctx, query, variables, &componentDTO, componentDTO.IsSuccessful)
	if err != nil {
		return fmt.Errorf("failed to unset component dependency %s -> %s: %v", dependentId, providerId, err)
	}

	return nil
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (*resources.Component, error) {
	componentDTO := dtos.ComponentByReferenceOutput{}
	query := componentDTO.GetQuery()
	variables := componentDTO.SetVariables(r.compass.GetCompassCloudId(), slug)

	err := r.run(ctx, query, variables, &componentDTO, componentDTO.IsSuccessful)
	if err != nil {
		return nil, fmt.Errorf("failed to get component by slug %s: %v", slug, err)
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

	err := r.run(ctx, query, variables, &documentDTO, documentDTO.IsSuccessful)
	if err != nil {
		return resources.Document{}, fmt.Errorf("failed to create document \"%s\" for %s: %v", document.Title, componentID, err)
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

	err := r.run(ctx, query, variables, &documentDTO, documentDTO.IsSuccessful)
	if err != nil {
		return fmt.Errorf("failed to update document \"%s\" for %s: %v", document.Title, componentID, err)
	}

	return nil
}

// @TODO work on this
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
		log.Printf("failed to delete metric source: %v", err)
		return err
	}

	if !response.Compass.DeleteComponentLink.Success {
		return errors.New("failed to delete metric source")
	}

	return nil
}

func (r *Repository) BindMetric(ctx context.Context, componentID string, metricID string, identifier string) (string, error) {
	bindMetricDTO := dtos.BindMetricOutput{}
	query := bindMetricDTO.GetQuery()
	variables := bindMetricDTO.SetVariables(metricID, componentID, identifier)

	err := r.run(ctx, query, variables, &bindMetricDTO, bindMetricDTO.IsSuccessful)
	if err != nil {
		return "", fmt.Errorf("failed to bind component %s to metric %s: %v", componentID, metricID, err)
	}

	return bindMetricDTO.Compass.CreateMetricSource.CreateMetricSource.ID, nil
}

func (r *Repository) UnbindMetric(ctx context.Context, metricSourceID string) error {
	unbindMetricDTO := dtos.UnbindMetricOutput{}
	query := unbindMetricDTO.GetQuery()
	variables := unbindMetricDTO.SetVariables(metricSourceID)

	err := r.run(ctx, query, variables, &unbindMetricDTO, unbindMetricDTO.IsSuccessful)
	if err != nil {
		return fmt.Errorf("failed to unbind metric source %s: %v", metricSourceID, err)
	}

	return nil
}

func (r *Repository) Push(ctx context.Context, metricSourceID string, value float64, recordedAt time.Time) error {
	requestBody := map[string]string{
		"metricSourceId": metricSourceID,
		"value":          fmt.Sprintf("%f", value),
		"timestamp":      recordedAt.UTC().Format(time.RFC3339),
	}

	_, errSend := r.compass.SendMetric(ctx, requestBody)

	return errSend
}

func (r *Repository) SetAPISpecifications(ctx context.Context, componentID, apiSpecs, apiSpecsFile string) error {
	lastSlashIndex := strings.LastIndex(componentID, "/")
	if lastSlashIndex == -1 {
		return errors.New("invalid componentID format")
	}

	input := compassdtos.APISpecificationsInput{
		ComponentID: componentID[lastSlashIndex+1:],
		ApiSpecs:    apiSpecs,
		FileName:    apiSpecsFile,
	}
	_, errSend := r.compass.SendAPISpecifications(ctx, input)
	return errSend
}

func (r *Repository) initDocumentCategories(ctx context.Context) error {
	if r.DocumentCategories != nil {
		return nil
	}

	documentationCategoriesDTO := dtos.DocumentationCategoriesOutput{}
	query := documentationCategoriesDTO.GetQuery()

	if err := r.compass.Run(ctx, query, nil, &documentationCategoriesDTO); err != nil {
		log.Printf("Failed to fetch document categories: %v", err)
		return err
	}

	categories := make(map[string]string, len(documentationCategoriesDTO.Compass.DocumentationCategories.Nodes))
	for _, category := range documentationCategoriesDTO.Compass.DocumentationCategories.Nodes {
		categories[category.Name] = category.ID
	}
	r.DocumentCategories = categories

	return nil
}

func (r *Repository) run(ctx context.Context, query string, variables map[string]interface{}, output interface{}, isSuccessful func() bool) error {
	if err := r.compass.Run(ctx, query, variables, output); err != nil {
		log.Printf("failed to create metric source: %v", err)
		return err
	}

	if !isSuccessful() {
		return fmt.Errorf("failed to run operation")
	}

	return nil
}
