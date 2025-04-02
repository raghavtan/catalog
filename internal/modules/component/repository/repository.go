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
	Delete(ctx context.Context, component resources.Component) error
	GetBySlug(ctx context.Context, component resources.Component) (*resources.Component, error)
	// Dependency operations
	SetDependency(ctx context.Context, dependent, provider resources.Component) error
	UnsetDependency(ctx context.Context, dependent, provider resources.Component) error
	// Documents operations
	AddDocument(ctx context.Context, component resources.Component, document resources.Document) (resources.Document, error)
	UpdateDocument(ctx context.Context, component resources.Component, document resources.Document) error
	RemoveDocument(ctx context.Context, component resources.Component, document resources.Document) error
	// MetricSource operations
	BindMetric(ctx context.Context, component resources.Component, metricID string, identifier string) (string, error)
	UnbindMetric(ctx context.Context, metricSource resources.MetricSource) error
	// API Specications operations
	SetAPISpecifications(ctx context.Context, component resources.Component, apiSpecs, apiSpecsFile string) error
	// Push metric value
	Push(ctx context.Context, metricSource resources.MetricSource, value float64, recordedAt time.Time) error
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
	input := &dtos.CreateComponentInput{CompassCloudID: r.compass.GetCompassCloudId(), Component: component}
	output := &dtos.CreateComponentOutput{}

	// This function is executed before the validation of the operation
	// That is before checking if the operation was successful
	// If the component already exists, it updates the component, sets the ID and metric sources
	// and clears the errors
	input.PreValidationFunc = func() error {
		if !compassservice.HasAlreadyExistsError(output.Compass.CreateComponent.Errors) {
			return nil
		}

		remoteComponent, runErr := r.GetBySlug(ctx, component)
		if runErr != nil {
			return runErr
		}

		component.ID = remoteComponent.ID
		component.MetricSources = remoteComponent.MetricSources
		_, updateError := r.Update(ctx, component)
		if updateError != nil {
			return updateError
		}

		output.Compass.CreateComponent.Details.ID = remoteComponent.ID
		output.Compass.CreateComponent.Errors = nil
		output.Compass.CreateComponent.Success = true

		return nil
	}

	runErr := r.compass.RunWithDTOs(ctx, input, output)
	if runErr != nil {
		return resources.Component{}, fmt.Errorf("Create component error for %s: %s", component.Name, runErr)
	}

	metricSources := make(map[string]*resources.MetricSource)
	for _, node := range output.Compass.CreateComponent.Details.MetricSources.Nodes {
		metricSources[node.MetricDefinition.Name] = &resources.MetricSource{
			ID:     node.ID,
			Metric: node.MetricDefinition.ID,
		}
	}

	createdLinks := make([]resources.Link, len(output.Compass.CreateComponent.Details.Links))
	for i, link := range output.Compass.CreateComponent.Details.Links {
		createdLinks[i] = resources.Link{
			ID:   link.ID,
			Type: link.Type,
			Name: link.Name,
			URL:  link.URL,
		}
	}
	component.ID = output.Compass.CreateComponent.Details.ID
	component.MetricSources = metricSources
	component.Links = createdLinks

	return component, nil
}

func (r *Repository) Update(ctx context.Context, component resources.Component) (resources.Component, error) {
	input := &dtos.UpdateComponentInput{Component: component}
	output := &dtos.UpdateComponentOutput{}

	// This function is executed before the validation of the operation
	// That is before checking if the operation was successful
	// If the component does not exist, it searches the component by slug, sets the ID and metric sources
	// and clears the errors
	input.PreValidationFunc = func() error {
		if !compassservice.HasNotFoundError(output.Compass.UpdateComponent.Errors) {
			return nil
		}

		remoteComponent, getBySlugErr := r.GetBySlug(ctx, component)
		if getBySlugErr != nil {
			return getBySlugErr
		}

		component.ID = remoteComponent.ID
		component.MetricSources = remoteComponent.MetricSources
		_, updateError := r.Update(ctx, component)
		if updateError != nil {
			return updateError
		}

		output.Compass.UpdateComponent.Errors = nil
		output.Compass.UpdateComponent.Success = true

		return nil
	}

	runErr := r.compass.RunWithDTOs(ctx, input, output)
	if runErr != nil {
		return resources.Component{}, fmt.Errorf("Update component error for %s: %s", component.Name, runErr)
	}

	return component, nil
}

func (r *Repository) Delete(ctx context.Context, component resources.Component) error {
	input := &dtos.DeleteComponentInput{ComponentID: component.ID}
	output := &dtos.DeleteComponentOutput{}
	if runErr := r.compass.RunWithDTOs(ctx, input, output); runErr != nil {
		return fmt.Errorf("Delete component error for %s: %s", component.ID, runErr)
	}
	return nil
}

func (r *Repository) SetDependency(ctx context.Context, dependent, provider resources.Component) error {
	input := &dtos.CreateDependencyInput{DependentId: dependent.ID, ProviderId: provider.ID}
	output := &dtos.CreateDependencyOutput{}
	if runErr := r.compass.RunWithDTOs(ctx, input, output); runErr != nil {
		return fmt.Errorf("SetDependency error for %s: %s", dependent.ID, runErr)
	}
	return nil
}

func (r *Repository) UnsetDependency(ctx context.Context, dependent, provider resources.Component) error {
	input := &dtos.DeleteDependencyInput{DependentId: dependent.ID, ProviderId: provider.ID}
	output := &dtos.DeleteDependencyOutput{}
	if runErr := r.compass.RunWithDTOs(ctx, input, output); runErr != nil {
		return fmt.Errorf("UnsetDependency dependency error for %s: %s", dependent.ID, runErr)
	}
	return nil
}

func (r *Repository) GetBySlug(ctx context.Context, component resources.Component) (*resources.Component, error) {
	input := &dtos.ComponentByReferenceInput{CompassCloudID: r.compass.GetCompassCloudId(), Slug: component.Slug}
	output := &dtos.ComponentByReferenceOutput{}
	runErr := r.compass.RunWithDTOs(ctx, input, output)
	if runErr != nil {
		return nil, fmt.Errorf("GetBySlug error for %s: %s", component.Slug, runErr)
	}

	metricSources := make(map[string]*resources.MetricSource)
	for _, node := range output.Compass.Component.MetricSources.Nodes {
		metricSources[node.MetricDefinition.Name] = &resources.MetricSource{
			ID:     node.ID,
			Metric: node.MetricDefinition.ID,
		}
	}

	found := resources.Component{
		ID:            output.Compass.Component.ID,
		MetricSources: metricSources,
	}

	return &found, nil
}

func (r *Repository) AddDocument(ctx context.Context, component resources.Component, document resources.Document) (resources.Document, error) {
	r.initDocumentCategories(ctx)

	input := &dtos.CreateDocumentInput{
		ComponentID: component.ID,
		Document:    resources.Document{Title: document.Title, URL: document.URL},
		CategoryID:  r.DocumentCategories[document.Type],
	}
	output := &dtos.CreateDocumentOutput{}
	runErr := r.compass.RunWithDTOs(ctx, input, output)
	if runErr != nil {
		return resources.Document{}, fmt.Errorf("AddDocument error for %s/%s: %s", component.ID, document.Title, runErr)
	}

	doc := resources.Document{
		ID:                      output.Compass.AddDocument.Details.ID,
		Title:                   document.Title,
		Type:                    document.Type,
		URL:                     document.URL,
		DocumentationCategoryId: r.DocumentCategories[document.Type],
	}
	return doc, nil
}

func (r *Repository) UpdateDocument(ctx context.Context, component resources.Component, document resources.Document) error {
	r.initDocumentCategories(ctx)

	input := &dtos.UpdateDocumentInput{
		Document:   document,
		CategoryID: r.DocumentCategories[document.Type],
	}
	output := &dtos.UpdateDocumentOutput{}
	if runErr := r.compass.RunWithDTOs(ctx, input, output); runErr != nil {
		return fmt.Errorf("UpdateDocument error for %s/%s: %s", component.ID, document.Title, runErr)
	}
	return nil
}

// @TODO work on this
func (r *Repository) RemoveDocument(ctx context.Context, component resources.Component, document resources.Document) error {
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
		"componentId": component.ID,
		"id":          document.ID,
	}

	var response struct {
		Compass struct {
			DeleteComponentLink struct {
				Success bool `json:"success"`
			} `json:"deleteComponentLink"`
		} `json:"compass"`
	}

	if runErr := r.compass.Run(ctx, query, variables, &response); runErr != nil {
		log.Printf("failed to delete metric source: %v", runErr)
		return runErr
	}

	if !response.Compass.DeleteComponentLink.Success {
		return errors.New("failed to delete metric source")
	}
	return nil
}

func (r *Repository) BindMetric(ctx context.Context, component resources.Component, metricID string, identifier string) (string, error) {
	input := &dtos.BindMetricInput{
		ComponentID: component.ID,
		MetricID:    metricID,
		Identifier:  identifier,
	}
	output := &dtos.BindMetricOutput{}
	runErr := r.compass.RunWithDTOs(ctx, input, output)
	if runErr != nil {
		return "", fmt.Errorf("BindMetric error for %s/%s: %s", component.ID, metricID, runErr)
	}

	return output.Compass.CreateMetricSource.CreateMetricSource.ID, nil
}

func (r *Repository) UnbindMetric(ctx context.Context, metricSource resources.MetricSource) error {
	input := &dtos.UnbindMetricInput{MetricID: metricSource.ID}
	output := &dtos.UnbindMetricOutput{}
	if runErr := r.compass.RunWithDTOs(ctx, input, output); runErr != nil {
		return fmt.Errorf("UnbindMetric error for %s: %s", metricSource.ID, runErr)
	}
	return nil
}

func (r *Repository) Push(ctx context.Context, metricSource resources.MetricSource, value float64, recordedAt time.Time) error {
	requestBody := map[string]string{
		"metricSourceId": metricSource.ID,
		"value":          fmt.Sprintf("%f", value),
		"timestamp":      recordedAt.UTC().Format(time.RFC3339),
	}

	_, errSend := r.compass.SendMetric(ctx, requestBody)

	return errSend
}

func (r *Repository) SetAPISpecifications(ctx context.Context, component resources.Component, apiSpecs, apiSpecsFile string) error {
	lastSlashIndex := strings.LastIndex(component.ID, "/")
	if lastSlashIndex == -1 {
		return errors.New("invalid component.ID format")
	}

	input := compassdtos.APISpecificationsInput{
		ComponentID: component.ID[lastSlashIndex+1:],
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

	input := &dtos.DocumentationCategoriesInput{CompassCloudID: r.compass.GetCompassCloudId()}
	output := &dtos.DocumentationCategoriesOutput{}
	runErr := r.compass.RunWithDTOs(ctx, input, output)
	if runErr != nil {
		return runErr
	}

	categories := make(map[string]string, len(output.Compass.DocumentationCategories.Nodes))
	for _, category := range output.Compass.DocumentationCategories.Nodes {
		categories[category.Name] = category.ID
	}
	r.DocumentCategories = categories

	return nil
}
