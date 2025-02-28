package handler

import (
	"context"
	"log"

	"github.com/motain/fact-collector/internal/modules/component/dtos"
	"github.com/motain/fact-collector/internal/modules/component/repository"
	"github.com/motain/fact-collector/internal/modules/component/resources"
	"github.com/motain/fact-collector/internal/modules/component/utils"
	"github.com/motain/fact-collector/internal/services/githubservice"
	"github.com/motain/fact-collector/internal/utils/drift"
	"github.com/motain/fact-collector/internal/utils/yaml"
)

type ApplyHandler struct {
	github     githubservice.GitHubServiceInterface
	repository repository.RepositoryInterface
}

func NewApplyHandler(
	gh githubservice.GitHubServiceInterface,
	repository repository.RepositoryInterface,
) *ApplyHandler {
	return &ApplyHandler{github: gh, repository: repository}
}

func (h *ApplyHandler) Apply(configRootLocation string, stateRootLocation string, recursive bool) {
	configComponents, errConfig := yaml.Parse[dtos.ComponentDTO](configRootLocation, recursive, dtos.GetComponentUniqueKey)
	if errConfig != nil {
		log.Fatalf("error: %v", errConfig)
	}

	stateComponents, errState := yaml.Parse[dtos.ComponentDTO](stateRootLocation, false, dtos.GetComponentUniqueKey)
	if errState != nil {
		log.Fatalf("error: %v", errState)
	}

	created, updated, deleted, unchanged := drift.Detect(
		stateComponents,
		configComponents,
		dtos.GetComponentID,
		dtos.SetComponentID,
		dtos.IsEqualComponent,
	)
	h.handleDeleted(deleted)

	var result []*dtos.ComponentDTO
	result = h.handleUnchanged(result, unchanged)
	result = h.handleCreated(result, created)
	result = h.handleUpdated(result, updated)

	err := yaml.WriteState[dtos.ComponentDTO](result)
	if err != nil {
		log.Fatalf("error writing components to file: %v", err)
	}
}

func (h *ApplyHandler) handleDeleted(components map[string]*dtos.ComponentDTO) {
	for _, componentDTO := range components {
		errComponent := h.repository.Delete(context.Background(), *componentDTO.Spec.ID)
		if errComponent != nil {
			panic(errComponent)
		}
	}
}

func (h *ApplyHandler) handleUnchanged(result []*dtos.ComponentDTO, components map[string]*dtos.ComponentDTO) []*dtos.ComponentDTO {
	for _, componentDTO := range components {
		result = append(result, componentDTO)
	}
	return result
}

func (h *ApplyHandler) handleCreated(result []*dtos.ComponentDTO, components map[string]*dtos.ComponentDTO) []*dtos.ComponentDTO {
	for _, componentDTO := range components {
		component := componentDTOToResource(componentDTO)

		id, errComponent := h.repository.Create(context.Background(), component)
		if errComponent != nil {
			panic(errComponent)
		}

		componentDTO.Spec.ID = &id
		componentDTO.Spec.Slug = component.Slug
		result = append(result, componentDTO)
	}

	return result
}

func (h *ApplyHandler) handleUpdated(result []*dtos.ComponentDTO, components map[string]*dtos.ComponentDTO) []*dtos.ComponentDTO {
	for _, componentDTO := range components {
		component := componentDTOToResource(componentDTO)
		errComponent := h.repository.Update(context.Background(), component)
		if errComponent != nil {
			panic(errComponent)
		}

		result = append(result, componentDTO)
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
		Links:         linksDTOToResource(componentDTO.Spec.Links),
		Labels:        componentDTO.Spec.Labels,
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
