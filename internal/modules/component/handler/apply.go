package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/motain/of-catalog/internal/modules/component/dtos"
	"github.com/motain/of-catalog/internal/modules/component/repository"
	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/services/documentservice"
	fsdtos "github.com/motain/of-catalog/internal/services/factsystem/dtos"
	"github.com/motain/of-catalog/internal/services/githubservice"
	"github.com/motain/of-catalog/internal/services/ownerservice"
	"github.com/motain/of-catalog/internal/utils/drift"
	listutils "github.com/motain/of-catalog/internal/utils/list"
	"github.com/motain/of-catalog/internal/utils/yaml"
	"log"
	"path/filepath"
)

type ApplyHandler struct {
	github     githubservice.GitHubServiceInterface
	repository repository.RepositoryInterface
	owner      ownerservice.OwnerServiceInterface
	document   documentservice.DocumentServiceInterface
	converter  *ComponentConverter
}

func NewApplyHandler(
	gh githubservice.GitHubServiceInterface,
	repository repository.RepositoryInterface,
	owner ownerservice.OwnerServiceInterface,
	document documentservice.DocumentServiceInterface,
) *ApplyHandler {
	return &ApplyHandler{
		github:     gh,
		repository: repository,
		owner:      owner,
		document:   document,
		converter:  NewComponentConverter(gh), // Initialize converter with GitHub service
	}
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

	stateComponents, errState := yaml.Parse(yaml.GetComponentStateInput(), dtos.GetComponentUniqueKey)
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
	// FIXED: Ensure owner is calculated correctly before drift detection
	// Apply handleOwner to all components before comparison
	correctedConfigComponents := make(map[string]*dtos.ComponentDTO)
	for name, component := range configComponents {
		correctedConfigComponents[name] = h.handleOwner(component)
	}

	correctedStateComponents := make(map[string]*dtos.ComponentDTO)
	for name, component := range stateComponents {
		correctedStateComponents[name] = h.handleOwner(component)
	}

	created, updated, deleted, unchanged := drift.Detect(
		correctedStateComponents,
		correctedConfigComponents,
		dtos.FromStateToConfig,
		dtos.IsEqualComponent,
	)

	var result []*dtos.ComponentDTO
	h.handleDeleted(ctx, deleted)
	result = h.handleUnchanged(ctx, result, unchanged, stateComponents)
	result = h.handleCreated(ctx, result, created, stateComponents)
	result = h.handleUpdated(ctx, result, updated, stateComponents)
	err := yaml.WriteComponentStates(yaml.SortResults(result, dtos.GetComponentUniqueKey), dtos.GetComponentUniqueKey)
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

	configComponentWithCorrectOwner := h.handleOwner(configComponent)

	var stateComponentWithCorrectOwner *dtos.ComponentDTO
	if stateComponents[componentName] != nil {
		stateComponentWithCorrectOwner = h.handleOwner(stateComponents[componentName])
		stateMap[componentName] = stateComponentWithCorrectOwner
	}

	created, updated, deleted, unchanged := drift.Detect(
		stateMap,
		map[string]*dtos.ComponentDTO{componentName: configComponentWithCorrectOwner},
		dtos.FromStateToConfig,
		dtos.IsEqualComponent,
	)
	fmt.Printf("DEBUG: created: %d, updated: %d, deleted: %d, unchanged: %d\n", len(created), len(updated), len(deleted), len(unchanged))

	h.handleDeleted(ctx, deleted)
	result = h.handleUnchanged(ctx, result, unchanged, stateComponents)
	result = h.handleCreated(ctx, result, created, stateComponents)
	result = h.handleUpdated(ctx, result, updated, stateComponents)

	err := yaml.WriteComponentStates(yaml.SortResults(result, dtos.GetComponentUniqueKey), dtos.GetComponentUniqueKey)
	if err != nil {
		log.Fatalf("error writing components to file: %v", err)
	}
}

func (h *ApplyHandler) handleDeleted(ctx context.Context, components map[string]*dtos.ComponentDTO) {
	for _, componentDTO := range components {
		errComponent := h.repository.Delete(ctx, h.converter.ToResource(componentDTO))
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
		componentDTO = h.handleDescription(componentDTO) // FIXED: Handle description for unchanged components too
		componentDTO = h.handleDocumenation(ctx, componentDTO, stateComponents)

		// FIXED: For unchanged components, ensure we preserve all state data including facts
		stateComponent := stateComponents[componentDTO.Metadata.Name]
		if stateComponent != nil {
			// Preserve MetricSources with facts from state
			if stateComponent.Spec.MetricSources != nil {
				componentDTO.Spec.MetricSources = make(map[string]*dtos.MetricSourceDTO)
				for metricName, stateMetricSource := range stateComponent.Spec.MetricSources {
					componentDTO.Spec.MetricSources[metricName] = &dtos.MetricSourceDTO{
						ID:     stateMetricSource.ID,
						Name:   stateMetricSource.Name,
						Metric: stateMetricSource.Metric,
						Facts:  stateMetricSource.Facts, // PRESERVE FACTS!
					}
				}
			}
		}

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
		componentDTO = h.handleDescription(componentDTO)

		// Should we call this at creation time?
		// componentDTO = h.handleDocumenation(componentDTO)

		component := h.converter.ToResource(componentDTO)

		component, errComponent := h.repository.Create(ctx, component)
		if errComponent != nil {
			panic(errComponent)
		}

		for _, providerName := range componentDTO.Spec.DependsOn {
			if provider, exists := stateComponents[providerName]; exists {
				h.repository.SetDependency(ctx, component, h.converter.ToResource(provider))
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
				Facts:  []*fsdtos.Task{}, // Empty facts for new sources
			}
		}

		// FIXED: Handle DependsOn duplicates - check if kubernetes already exists
		if len(componentDTO.Spec.DependsOn) == 0 {
			// Default mandatory element [kubernetes]
			componentDTO.Spec.DependsOn = []string{"kubernetes"}
		} else {
			// Check if kubernetes is already in the list to avoid duplicates
			hasKubernetes := false
			for _, dep := range componentDTO.Spec.DependsOn {
				if dep == "kubernetes" {
					hasKubernetes = true
					break
				}
			}
			if !hasKubernetes {
				componentDTO.Spec.DependsOn = append(componentDTO.Spec.DependsOn, "kubernetes")
			}
		}

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
		componentDTO = h.handleDescription(componentDTO)
		componentDTO = h.handleDocumenation(ctx, componentDTO, stateComponents)

		component := h.converter.ToResource(componentDTO)
		component, errComponent := h.repository.Update(ctx, component)
		if errComponent != nil {
			panic(errComponent)
		}

		componentDTO.Spec.ID = component.ID

		updatedLinks := make([]dtos.Link, len(component.Links))
		for i, link := range component.Links {
			updatedLinks[i] = dtos.Link{
				ID:   link.ID,
				Name: link.Name,
				Type: link.Type,
				URL:  link.URL,
			}
		}
		componentDTO.Spec.Links = updatedLinks

		// FIXED: Preserve existing MetricSources with facts, only update basic properties
		stateComponent := stateComponents[componentDTO.Metadata.Name]
		if stateComponent != nil && stateComponent.Spec.MetricSources != nil {
			if componentDTO.Spec.MetricSources == nil {
				componentDTO.Spec.MetricSources = make(map[string]*dtos.MetricSourceDTO)
			}

			// Copy existing metric sources with their facts
			for metricName, stateMetricSource := range stateComponent.Spec.MetricSources {
				componentDTO.Spec.MetricSources[metricName] = &dtos.MetricSourceDTO{
					ID:     stateMetricSource.ID,     // Keep existing ID
					Name:   stateMetricSource.Name,   // Keep existing name
					Metric: stateMetricSource.Metric, // Keep existing metric
					Facts:  stateMetricSource.Facts,  // PRESERVE FACTS!
				}
			}

			// Update any metric sources that changed in the repository response
			for metricName, metricSource := range component.MetricSources {
				if existingMetricSource, exists := componentDTO.Spec.MetricSources[metricName]; exists {
					// Update basic properties but keep facts
					existingMetricSource.ID = metricSource.ID
					existingMetricSource.Name = metricSource.Name
					existingMetricSource.Metric = metricSource.Metric
					// Facts are preserved from state
				} else {
					// New metric source
					componentDTO.Spec.MetricSources[metricName] = &dtos.MetricSourceDTO{
						ID:     metricSource.ID,
						Name:   metricSource.Name,
						Metric: metricSource.Metric,
						Facts:  []*fsdtos.Task{}, // Empty facts for new sources
					}
				}
			}
		} else {
			// Fallback if no state component (shouldn't happen in update)
			if componentDTO.Spec.MetricSources == nil {
				componentDTO.Spec.MetricSources = make(map[string]*dtos.MetricSourceDTO)
			}
			for metricName, metricSource := range component.MetricSources {
				componentDTO.Spec.MetricSources[metricName] = &dtos.MetricSourceDTO{
					ID:     metricSource.ID,
					Name:   metricSource.Name,
					Metric: metricSource.Metric,
					Facts:  []*fsdtos.Task{}, // Empty facts for new sources
				}
			}
		}

		h.handleDependencies(ctx, componentDTO, stateComponents)

		result = append(result, componentDTO)

		h.handleAPISpecification(ctx, componentDTO)
	}

	return result
}

func (h *ApplyHandler) handleOwner(componentDTO *dtos.ComponentDTO) *dtos.ComponentDTO {

	// FIXED: Always recalculate owner based on tribe/squad
	// This ensures correct OwnerID even if component already has one set

	if componentDTO.Spec.Tribe != "" && componentDTO.Spec.Squad != "" {

		owner, ownerErr := h.owner.GetOwnerByTribeAndSquad(componentDTO.Spec.Tribe, componentDTO.Spec.Squad)
		if ownerErr == nil {

			if componentDTO.Spec.OwnerID != "" && componentDTO.Spec.OwnerID != owner.OwnerID {
				fmt.Printf("INFO: Updating OwnerID for %s from %s to %s (squad: %s)\n",
					componentDTO.Spec.Name, componentDTO.Spec.OwnerID, owner.OwnerID, componentDTO.Spec.Squad)
			} else if componentDTO.Spec.OwnerID == "" {
				fmt.Printf("INFO: Setting OwnerID for %s to %s (squad: %s)\n",
					componentDTO.Spec.Name, owner.OwnerID, componentDTO.Spec.Squad)
			}

			componentDTO.Spec.OwnerID = owner.OwnerID

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
		} else {
			// If no owner is found, keep the existing OwnerID (don't clear it)
			fmt.Printf("WARNING: Owner lookup failed for tribe '%s', squad '%s': %v\n",
				componentDTO.Spec.Tribe, componentDTO.Spec.Squad, ownerErr)
		}
	} else {
		fmt.Printf("WARNING: Tribe or Squad not set for component %s (tribe: '%s', squad: '%s')\n",
			componentDTO.Spec.Name, componentDTO.Spec.Tribe, componentDTO.Spec.Squad)
	}

	return componentDTO
}

// handleDescription fetches and sets repository description from GitHub
func (h *ApplyHandler) handleDescription(componentDTO *dtos.ComponentDTO) *dtos.ComponentDTO {
	if componentDTO.Spec.Description != "" {
		return componentDTO
	}
	description, err := h.github.GetRepoDescription(componentDTO.Metadata.Name)
	if err != nil {
		log.Printf("Warning: Could not get repository description for %s: %v", componentDTO.Metadata.Name, err)
		componentDTO.Spec.Description = fmt.Sprintf("Component %s", componentDTO.Spec.Name)
	} else {
		fmt.Printf("INFO: Setting description for %s from repository: %s\n", componentDTO.Spec.Name, description)
		componentDTO.Spec.Description = description
	}

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
			h.repository.RemoveDocument(ctx, h.converter.ToResource(componentDTO), documentResource)
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

			newDocument, addDocumentErr := h.repository.AddDocument(ctx, h.converter.ToResource(componentDTO), documentResource)
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

			updateDocumentErr := h.repository.UpdateDocument(ctx, h.converter.ToResource(componentDTO), documentResource)
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
	if componentInState == nil {
		return // or handle this case appropriately
	}

	for _, providerName := range componentInState.Spec.DependsOn {
		if !listutils.Contains(componentDTO.Spec.DependsOn, providerName) {
			stateProvider := stateComponents[providerName]
			if stateProvider == nil {
				log.Printf("Provider %s not found in state for component %s", providerName, componentDTO.Spec.Name)
				continue
			}

			err := h.repository.UnsetDependency(ctx, h.converter.ToResource(componentDTO), h.converter.ToResource(stateProvider))
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

			err := h.repository.SetDependency(ctx, h.converter.ToResource(componentDTO), h.converter.ToResource(stateProvider))
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
		docType := h.determineDocumentType(documentTitle, documentURL)

		mappedDocuments[documentTitle] = &dtos.Document{
			Title: documentTitle,
			Type:  docType,
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

func (h *ApplyHandler) determineDocumentType(title, url string) string {
	return "OTHER"
}

func (h *ApplyHandler) handleAPISpecification(ctx context.Context, componentDTO *dtos.ComponentDTO) {
	apiSpecs, apiSpecsFile, documentErr := h.getRemoteAPISpecifications(componentDTO.Spec.Name)
	if documentErr != nil {
		return
	}

	err := h.repository.SetAPISpecifications(ctx, h.converter.ToResource(componentDTO), apiSpecs, apiSpecsFile)
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
