package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"github.com/motain/of-catalog/internal/modules/component/dtos"
	"github.com/motain/of-catalog/internal/modules/component/repository"
	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/modules/component/utils"
	"github.com/motain/of-catalog/internal/services/documentservice"
	"github.com/motain/of-catalog/internal/services/githubservice"
	"github.com/motain/of-catalog/internal/services/ownerservice"
	"github.com/motain/of-catalog/internal/utils/drift"
	listutils "github.com/motain/of-catalog/internal/utils/list"
	"github.com/motain/of-catalog/internal/utils/yaml"
)

type ApplyHandler struct {
	github     githubservice.GitHubServiceInterface
	repository repository.RepositoryInterface
	owner      ownerservice.OwnerServiceInterface
	document   documentservice.DocumentServiceInterface
}

func NewApplyHandler(
	gh githubservice.GitHubServiceInterface,
	repository repository.RepositoryInterface,
	owner ownerservice.OwnerServiceInterface,
	document documentservice.DocumentServiceInterface,
) *ApplyHandler {
	return &ApplyHandler{github: gh, repository: repository, owner: owner, document: document}
}

func (h *ApplyHandler) Apply(ctx context.Context, configRootLocation string, stateRootLocation string, recursive bool, componentName string) {
	parseInput := yaml.ParseInput{
		RootLocation: configRootLocation,
		Recursive:    recursive,
	}
	configComponents, errConfig := yaml.Parse(parseInput, dtos.GetComponentUniqueKey)
	if errConfig != nil {
		log.Fatalf("error: %v", errConfig)
	}

	stateComponents, errState := yaml.Parse(yaml.GetStateInput(stateRootLocation), dtos.GetComponentUniqueKey)
	if errState != nil {
		log.Fatalf("error: %v", errState)
	}

	if componentName == "" {
		h.handleAll(ctx, stateComponents, configComponents)
		return
	}

	_, existsInState := stateComponents[componentName]
	_, existsInConfig := configComponents[componentName]
	if !existsInConfig && !existsInState {
		log.Fatalf("component %s not found", componentName)
	}

	h.handleOne(ctx, stateComponents, configComponents, componentName)
}

func (h *ApplyHandler) handleAll(ctx context.Context, stateComponents, configComponents map[string]*dtos.ComponentDTO) {
	created, updated, deleted, unchanged := drift.Detect(
		stateComponents,
		configComponents,
		dtos.FromStateToConfig,
		dtos.IsEqualComponent,
	)

	var result []*dtos.ComponentDTO
	h.handleDeleted(ctx, deleted)
	result = h.handleUnchanged(ctx, result, unchanged, stateComponents)
	result = h.handleCreated(ctx, result, created, stateComponents)
	result = h.handleUpdated(ctx, result, updated, stateComponents)

	err := yaml.WriteState(yaml.SortResults(result, dtos.GetComponentUniqueKey))
	if err != nil {
		log.Fatalf("error writing components to file: %v", err)
	}
}

func (h *ApplyHandler) handleOne(ctx context.Context, stateComponents, configComponents map[string]*dtos.ComponentDTO, componentName string) {
	configComponent := configComponents[componentName]
	result := make([]*dtos.ComponentDTO, 0)
	for stateComponentName, stateComponent := range stateComponents {
		if stateComponentName != componentName {
			result = append(result, stateComponent)
			continue
		}
	}

	stateMap := make(map[string]*dtos.ComponentDTO)
	if stateComponents[componentName] != nil {
		stateMap[componentName] = stateComponents[componentName]
	}

	created, updated, deleted, unchanged := drift.Detect(
		stateMap,
		map[string]*dtos.ComponentDTO{componentName: configComponent},
		dtos.FromStateToConfig,
		dtos.IsEqualComponent,
	)

	h.handleDeleted(ctx, deleted)
	result = h.handleUnchanged(ctx, result, unchanged, stateComponents)
	result = h.handleCreated(ctx, result, created, stateComponents)
	result = h.handleUpdated(ctx, result, updated, stateComponents)

	err := yaml.WriteState(yaml.SortResults(result, dtos.GetComponentUniqueKey))
	if err != nil {
		log.Fatalf("error writing components to file: %v", err)
	}
}

func (h *ApplyHandler) handleDeleted(ctx context.Context, components map[string]*dtos.ComponentDTO) {
	for _, componentDTO := range components {
		errComponent := h.repository.Delete(ctx, componentDTOToResource(componentDTO))
		if errComponent != nil {
			panic(errComponent)
		}
	}
}

func (h *ApplyHandler) handleUnchanged(
	ctx context.Context,
	result []*dtos.ComponentDTO,
	components map[string]*dtos.ComponentDTO,
	stateComponents map[string]*dtos.ComponentDTO,
) []*dtos.ComponentDTO {
	for _, componentDTO := range components {
		componentDTO = h.handleOwner(componentDTO)
		componentDTO = h.handleDocumenation(ctx, componentDTO, stateComponents)

		result = append(result, componentDTO)

		h.handleDependencies(ctx, componentDTO, stateComponents)

		h.handleAPISpecification(ctx, componentDTO)
	}
	return result
}

func (h *ApplyHandler) handleCreated(
	ctx context.Context,
	result []*dtos.ComponentDTO,
	components map[string]*dtos.ComponentDTO,
	stateComponents map[string]*dtos.ComponentDTO,
) []*dtos.ComponentDTO {
	for _, componentDTO := range components {
		componentDTO = h.handleOwner(componentDTO)

		// Should we call this at creation time?
		// componentDTO = h.handleDocumenation(componentDTO)

		component := componentDTOToResource(componentDTO)

		component, errComponent := h.repository.Create(ctx, component)
		if errComponent != nil {
			panic(errComponent)
		}

		for _, providerName := range componentDTO.Spec.DependsOn {
			if provider, exists := stateComponents[providerName]; exists {
				h.repository.SetDependency(ctx, component, componentDTOToResource(provider))
			} else {
				log.Printf("Provider %s not found for component %s", providerName, componentDTO.Spec.Name)
			}
		}

		componentDTO.Spec.ID = component.ID
		componentDTO.Spec.Slug = component.Slug

		createdLinks := make([]dtos.Link, len(component.Links))
		for i, link := range component.Links {
			createdLinks[i] = dtos.Link{
				ID:   link.ID,
				Name: link.Name,
				Type: link.Type,
				URL:  link.URL,
			}
		}
		componentDTO.Spec.Links = createdLinks

		if componentDTO.Spec.MetricSources == nil {
			componentDTO.Spec.MetricSources = make(map[string]*dtos.MetricSourceDTO)
		}
		for metricName, metricSource := range component.MetricSources {
			componentDTO.Spec.MetricSources[metricName] = &dtos.MetricSourceDTO{
				ID:     metricSource.ID,
				Name:   metricSource.Name,
				Metric: metricSource.Metric,
			}
		}
		// At this point dependecies may not be set because we set dependecies after creating the component
		// But the dependency for this component may not have been create yet.
		// We set nil and we will update it later when running aplpy again.
		// Eventually to make it more clear we can create a specific command to set dependencies
		// We need to think about the best way to handle this
		componentDTO.Spec.DependsOn = nil
		result = append(result, componentDTO)

		h.handleDependencies(ctx, componentDTO, stateComponents)

		h.handleAPISpecification(ctx, componentDTO)
	}

	return result
}

func (h *ApplyHandler) handleUpdated(
	ctx context.Context,
	result []*dtos.ComponentDTO,
	components map[string]*dtos.ComponentDTO,
	stateComponents map[string]*dtos.ComponentDTO,
) []*dtos.ComponentDTO {
	for _, componentDTO := range components {
		componentDTO = h.handleOwner(componentDTO)
		componentDTO = h.handleDocumenation(ctx, componentDTO, stateComponents)

		component := componentDTOToResource(componentDTO)
		component, errComponent := h.repository.Update(ctx, component)
		if errComponent != nil {
			panic(errComponent)
		}

		componentDTO.Spec.ID = component.ID
		for _, link := range component.Links {
			componentDTO.Spec.Links = append(componentDTO.Spec.Links, dtos.Link{
				ID:   link.ID,
				Name: link.Name,
				Type: link.Type,
				URL:  link.URL,
			})
		}

		if componentDTO.Spec.MetricSources == nil {
			componentDTO.Spec.MetricSources = make(map[string]*dtos.MetricSourceDTO)
		}
		for metricName, metricSource := range component.MetricSources {
			componentDTO.Spec.MetricSources[metricName] = &dtos.MetricSourceDTO{
				ID:     metricSource.ID,
				Name:   metricSource.Name,
				Metric: metricSource.Metric,
			}
		}

		h.handleDependencies(ctx, componentDTO, stateComponents)

		result = append(result, componentDTO)

		h.handleAPISpecification(ctx, componentDTO)
	}

	return result
}

func componentDTOToResource(componentDTO *dtos.ComponentDTO) resources.Component {
	return resources.Component{
		ID:            componentDTO.Spec.ID,
		Name:          componentDTO.Spec.Name,
		Slug:          utils.GetSlug(componentDTO.Spec.Name, componentDTO.Spec.TypeID),
		Description:   componentDTO.Spec.Description,
		ConfigVersion: componentDTO.Spec.ConfigVersion,
		TypeID:        componentDTO.Spec.TypeID,
		OwnerID:       componentDTO.Spec.OwnerID,
		Fields:        componentDTO.Spec.Fields,
		Links:         linksDTOToResource(componentDTO.Spec.Links),
		Labels:        componentDTO.Spec.Labels,
		MetricSources: metricSourcesDTOToResource(componentDTO.Spec.MetricSources),
	}
}

func linksDTOToResource(linksDTO []dtos.Link) []resources.Link {
	links := make([]resources.Link, 0)
	for _, link := range linksDTO {
		links = append(links, resources.Link{
			Name: link.Name,
			Type: link.Type,
			URL:  link.URL,
		})
	}
	return links
}

func metricSourcesDTOToResource(metricSourcesDTO map[string]*dtos.MetricSourceDTO) map[string]*resources.MetricSource {
	metricSources := make(map[string]*resources.MetricSource)
	for metricName, metricSourceDTO := range metricSourcesDTO {
		metricSources[metricName] = &resources.MetricSource{
			ID:     metricSourceDTO.ID,
			Name:   metricSourceDTO.Name,
			Metric: metricSourceDTO.Metric,
		}
	}
	return metricSources
}

func (h *ApplyHandler) handleOwner(componentDTO *dtos.ComponentDTO) *dtos.ComponentDTO {
	owner, ownerErr := h.owner.GetOwnerByTribeAndSquad(componentDTO.Spec.Tribe, componentDTO.Spec.Squad)
	if ownerErr != nil {
		// If no owner is found, we do not update the component
		return componentDTO
	}

	computedLinks := make(map[string]dtos.Link, 0)
	for _, link := range componentDTO.Spec.Links {
		computedLinks[link.Type+link.Name] = link
	}

	for slackChannelName, slackChannelURL := range owner.SlackChannels {
		computedLinks["CHAT_CHANNEL"+slackChannelName] = dtos.Link{
			Name: slackChannelName,
			Type: "CHAT_CHANNEL",
			URL:  slackChannelURL,
		}
	}

	for projectName, projectURL := range owner.Projects {
		computedLinks["PROJECT"+projectName] = dtos.Link{
			Name: projectName,
			Type: "PROJECT",
			URL:  projectURL,
		}
	}

	links := make([]dtos.Link, 0)
	for _, link := range computedLinks {
		links = append(links, link)
	}
	componentDTO.Spec.Links = links
	componentDTO.Spec.OwnerID = owner.OwnerID

	return componentDTO
}

func (h *ApplyHandler) handleDocuments(
	ctx context.Context,
	componentDTO *dtos.ComponentDTO,
	stateComponents map[string]*dtos.ComponentDTO,
) *dtos.ComponentDTO {
	resultDocuments := make(map[string]*dtos.Document, 0)
	componentInState := stateComponents[componentDTO.Metadata.Name]

	mappedStateDocuments := make(map[string]*dtos.Document, len(componentInState.Spec.Documents))
	for _, document := range componentInState.Spec.Documents {
		mappedStateDocuments[document.Title] = document
	}

	mappedComponentDocuments := make(map[string]*dtos.Document, len(componentDTO.Spec.Documents))
	for _, document := range componentDTO.Spec.Documents {
		mappedComponentDocuments[document.Title] = document
	}

	for _, document := range mappedStateDocuments {
		if _, exists := mappedComponentDocuments[document.Title]; !exists {
			documentResource := resources.Document{
				Title: document.Title,
				Type:  document.Type,
				URL:   document.URL,
			}
			h.repository.RemoveDocument(ctx, componentDTOToResource(componentDTO), documentResource)
			continue
		}

		resultDocuments[document.Title] = document
	}

	for _, document := range mappedComponentDocuments {
		if _, exists := mappedStateDocuments[document.Title]; !exists {
			documentResource := resources.Document{
				Title: document.Title,
				Type:  document.Type,
				URL:   document.URL,
			}

			newDocument, addDocumentErr := h.repository.AddDocument(ctx, componentDTOToResource(componentDTO), documentResource)
			if addDocumentErr != nil {
				fmt.Printf("apply documents %s", addDocumentErr)
			}

			document.ID = newDocument.ID
			document.DocumentationCategoryId = newDocument.DocumentationCategoryId
			resultDocuments[document.Title] = document

			continue
		}

		if document.URL != mappedStateDocuments[document.Title].URL {
			documentResource := resources.Document{
				ID:    mappedStateDocuments[document.Title].ID,
				Title: document.Title,
				Type:  document.Type,
				URL:   document.URL,
			}

			updateDocumentErr := h.repository.UpdateDocument(ctx, componentDTOToResource(componentDTO), documentResource)
			if updateDocumentErr != nil {
				fmt.Printf("apply documents %s", updateDocumentErr)
			}

			document.ID = mappedStateDocuments[document.Title].ID
			document.DocumentationCategoryId = mappedStateDocuments[document.Title].DocumentationCategoryId
			resultDocuments[document.Title] = document

			continue
		}
	}

	componentDTO.Spec.Documents = make([]*dtos.Document, 0)
	for _, document := range resultDocuments {
		componentDTO.Spec.Documents = append(componentDTO.Spec.Documents, document)
	}

	return componentDTO
}

func (h *ApplyHandler) handleDependencies(
	ctx context.Context,
	componentDTO *dtos.ComponentDTO,
	stateComponents map[string]*dtos.ComponentDTO,
) {
	componentInState := stateComponents[componentDTO.Metadata.Name]
	for _, providerName := range componentInState.Spec.DependsOn {
		if !listutils.Contains(componentDTO.Spec.DependsOn, providerName) {
			err := h.repository.UnsetDependency(ctx, componentDTOToResource(componentDTO), componentDTOToResource(stateComponents[providerName]))
			if err != nil {
				fmt.Printf("apply dependencies %s", err)
			}
		}
	}

	for _, providerName := range componentDTO.Spec.DependsOn {
		if !listutils.Contains(componentInState.Spec.DependsOn, providerName) {
			stateProvider, exists := stateComponents[providerName]
			if !exists {
				log.Printf("Provider %s not found for component %s", providerName, componentDTO.Spec.Name)
				continue
			}

			err := h.repository.SetDependency(ctx, componentDTOToResource(componentDTO), componentDTOToResource(stateProvider))
			if err != nil {
				fmt.Printf("apply dependencies %s", err)
			}
		}
	}
}

func (h *ApplyHandler) handleDocumenation(
	ctx context.Context,
	componentDTO *dtos.ComponentDTO,
	stateComponents map[string]*dtos.ComponentDTO,
) *dtos.ComponentDTO {
	documents, documentErr := h.document.GetDocuments(componentDTO.Spec.Name)
	if documentErr != nil {
		return componentDTO
	}

	mappedDocuments := make(map[string]*dtos.Document)
	for _, doc := range componentDTO.Spec.Documents {
		mappedDocuments[doc.Title] = doc
	}

	for documentTitle, documentURL := range documents {
		mappedDocuments[documentTitle] = &dtos.Document{
			Title: documentTitle,
			Type:  "Other",
			URL:   documentURL,
		}
	}

	processedDocuments := make([]*dtos.Document, len(mappedDocuments))
	i := 0
	for _, document := range mappedDocuments {
		processedDocuments[i] = document
		i++
	}
	componentDTO.Spec.Documents = processedDocuments

	return h.handleDocuments(ctx, componentDTO, stateComponents)
}

func (h *ApplyHandler) handleAPISpecification(ctx context.Context, componentDTO *dtos.ComponentDTO) {
	apiSpecs, apiSpecsFile, documentErr := h.getRemoteAPISpecifications(componentDTO.Spec.Name)
	if documentErr != nil {
		return
	}

	err := h.repository.SetAPISpecifications(ctx, componentDTOToResource(componentDTO), apiSpecs, apiSpecsFile)
	if err != nil {
		fmt.Printf("apply api specifications error: %s", err)
	}
}

func (h *ApplyHandler) getRemoteAPISpecifications(repo string) (string, string, error) {
	possibleLocations := []string{
		"",        // Let's assume the standard is to use the root folder
		"docs",    // Fallback to the docs folder
		"doc",     // Fallback to the doc folder
		".of",     // Fallback to the .of folder
		"openapi", // Fallback to the openapi folder
	}
	possibleFileNames := []string{
		"openapi.yaml",
		"openapi.yml",
		"openapi.json",
		"swagger.yaml",
		"swagger.yml",
		"swagger.json",
	}

	for _, folder := range possibleLocations {
		for _, fileName := range possibleFileNames {
			location := filepath.Join(folder, fileName)
			fileContent, fileErr := h.github.GetFileContent(repo, location)
			if fileErr == nil {
				return fileContent, location, nil
			}
		}
	}

	return "", "", errors.New("no API specification found")
}
