package repository

//go:generate mockgen -destination=./mocks/mock_repository.go -package=repository github.com/motain/of-catalog/internal/modules/component/repository RepositoryInterface

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/motain/of-catalog/internal/interfaces/repositoryinterfaces"
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
	input := &dtos.CreateComponentInput{CompassCloudID: r.compass.GetCompassCloudId(), Component: component}
	output := &dtos.CreateComponentOutput{}

	// This function is executed before the validation of the operation
	// That is before checking if the operation was successful
	// If the component already exists, it updates the component, sets the ID and metric sources
	// and clears the errors
	preValidationFunc := func() error {
		if !compassservice.HasAlreadyExistsError(output.Compass.CreateComponent.Errors) {
			return nil
		}

		remoteComponent, runErr := r.GetBySlug(ctx, component.Slug)
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

	runErr := r.run(ctx, input, output, preValidationFunc)
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
	preValidationFunc := func() error {
		if !compassservice.HasNotFoundError(output.Compass.UpdateComponent.Errors) {
			return nil
		}

		remoteComponent, getBySlugErr := r.GetBySlug(ctx, component.Slug)
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

	runErr := r.run(ctx, input, output, preValidationFunc)
	if runErr != nil {
		return resources.Component{}, fmt.Errorf("Update component error for %s: %s", component.Name, runErr)
	}

	return component, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	input := &dtos.DeleteComponentInput{ComponentID: id}
	output := &dtos.DeleteComponentOutput{}
	if runErr := r.run(ctx, input, output, nil); runErr != nil {
		return fmt.Errorf("Delete component error for %s: %s", id, runErr)
	}
	return nil
}

func (r *Repository) SetDependency(ctx context.Context, dependentId, providerId string) error {
	input := &dtos.CreateDependencyInput{DependentId: dependentId, ProviderId: providerId}
	output := &dtos.CreateDependencyOutput{}
	if runErr := r.run(ctx, input, output, nil); runErr != nil {
		return fmt.Errorf("SetDependency error for %s: %s", dependentId, runErr)
	}
	return nil
}

func (r *Repository) UnsetDependency(ctx context.Context, dependentId, providerId string) error {
	input := &dtos.DeleteDependencyInput{DependentId: dependentId, ProviderId: providerId}
	output := &dtos.DeleteDependencyOutput{}
	if runErr := r.run(ctx, input, output, nil); runErr != nil {
		return fmt.Errorf("UnsetDependency dependency error for %s: %s", dependentId, runErr)
	}
	return nil
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (*resources.Component, error) {
	input := &dtos.ComponentByReferenceInput{CompassCloudID: r.compass.GetCompassCloudId(), Slug: slug}
	output := &dtos.ComponentByReferenceOutput{}
	runErr := r.run(ctx, input, output, nil)
	if runErr != nil {
		return nil, fmt.Errorf("GetBySlug error for %s: %s", slug, runErr)
	}

	metricSources := make(map[string]*resources.MetricSource)
	for _, node := range output.Compass.Component.MetricSources.Nodes {
		metricSources[node.MetricDefinition.Name] = &resources.MetricSource{
			ID:     node.ID,
			Metric: node.MetricDefinition.ID,
		}
	}

	component := resources.Component{
		ID:            output.Compass.Component.ID,
		MetricSources: metricSources,
	}

	return &component, nil
}

func (r *Repository) AddDocument(ctx context.Context, componentID string, document resources.Document) (resources.Document, error) {
	r.initDocumentCategories(ctx)

	input := &dtos.CreateDocumentInput{
		ComponentID: componentID,
		Document:    resources.Document{Title: document.Title, URL: document.URL},
		CategoryID:  r.DocumentCategories[document.Type],
	}
	output := &dtos.CreateDocumentOutput{}
	runErr := r.run(ctx, input, output, nil)
	if runErr != nil {
		return resources.Document{}, fmt.Errorf("AddDocument error for %s/%s: %s", componentID, document.Title, runErr)
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

func (r *Repository) UpdateDocument(ctx context.Context, componentID string, document resources.Document) error {
	r.initDocumentCategories(ctx)

	input := &dtos.UpdateDocumentInput{
		Document:   document,
		CategoryID: r.DocumentCategories[document.Type],
	}
	output := &dtos.UpdateDocumentOutput{}
	if runErr := r.run(ctx, input, output, nil); runErr != nil {
		return fmt.Errorf("UpdateDocument error for %s/%s: %s", componentID, document.Title, runErr)
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

	if runErr := r.compass.Run(ctx, query, variables, &response); runErr != nil {
		log.Printf("failed to delete metric source: %v", runErr)
		return runErr
	}

	if !response.Compass.DeleteComponentLink.Success {
		return errors.New("failed to delete metric source")
	}
	return nil
}

func (r *Repository) BindMetric(ctx context.Context, componentID string, metricID string, identifier string) (string, error) {
	input := &dtos.BindMetricInput{
		ComponentID: componentID,
		MetricID:    metricID,
		Identifier:  identifier,
	}
	output := &dtos.BindMetricOutput{}
	runErr := r.run(ctx, input, output, nil)
	if runErr != nil {
		return "", fmt.Errorf("BindMetric error for %s/%s: %s", componentID, metricID, runErr)
	}

	return output.Compass.CreateMetricSource.CreateMetricSource.ID, nil
}

func (r *Repository) UnbindMetric(ctx context.Context, metricSourceID string) error {
	input := &dtos.UnbindMetricInput{MetricID: metricSourceID}
	output := &dtos.UnbindMetricOutput{}
	if runErr := r.run(ctx, input, output, nil); runErr != nil {
		return fmt.Errorf("UnbindMetric error for %s: %s", metricSourceID, runErr)
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

	input := &dtos.DocumentationCategoriesInput{CompassCloudID: r.compass.GetCompassCloudId()}
	output := &dtos.DocumentationCategoriesOutput{}
	runErr := r.run(ctx, input, output, nil)
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
