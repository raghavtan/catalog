package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/motain/of-catalog/internal/modules/component/dtos"
	"github.com/motain/of-catalog/internal/modules/component/repository"
	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/modules/component/utils"
	"github.com/motain/of-catalog/internal/services/githubservice"
	"github.com/motain/of-catalog/internal/services/ownerservice"
	ownerservicedtos "github.com/motain/of-catalog/internal/services/ownerservice/dtos"
	"github.com/motain/of-catalog/internal/utils/drift"
	listutils "github.com/motain/of-catalog/internal/utils/list"
	"github.com/motain/of-catalog/internal/utils/yaml"
)

type ApplyHandler struct {
	github     githubservice.GitHubServiceInterface
	repository repository.RepositoryInterface
	owner      ownerservice.OwnerServiceInterface
}

func NewApplyHandler(
	gh githubservice.GitHubServiceInterface,
	repository repository.RepositoryInterface,
	owner ownerservice.OwnerServiceInterface,
) *ApplyHandler {
	return &ApplyHandler{github: gh, repository: repository, owner: owner}
}

func (h *ApplyHandler) Apply(configRootLocation string, stateRootLocation string, recursive bool) {
	configComponents, errConfig := yaml.Parse[dtos.ComponentDTO](configRootLocation, recursive, dtos.GetComponentUniqueKey)
	if errConfig != nil {
		log.Fatalf("error: %v", errConfig)
	}

	stateComponents, errState := yaml.Parse(stateRootLocation, false, dtos.GetComponentUniqueKey)
	if errState != nil {
		log.Fatalf("error: %v", errState)
	}

	created, updated, deleted, unchanged := drift.Detect(
		stateComponents,
		configComponents,
		dtos.FromStateToConfig,
		dtos.IsEqualComponent,
	)
	h.handleDeleted(deleted)

	var result []*dtos.ComponentDTO
	result = h.handleUnchanged(result, unchanged)
	result = h.handleCreated(result, created, stateComponents)
	result = h.handleUpdated(result, updated, stateComponents)

	err := yaml.WriteState(result)
	if err != nil {
		log.Fatalf("error writing components to file: %v", err)
	}
}

func (h *ApplyHandler) handleDeleted(components map[string]*dtos.ComponentDTO) {
	for _, componentDTO := range components {
		errComponent := h.repository.Delete(context.Background(), componentDTO.Spec.ID)
		if errComponent != nil {
			panic(errComponent)
		}
	}
}

func (h *ApplyHandler) handleUnchanged(result []*dtos.ComponentDTO, components map[string]*dtos.ComponentDTO) []*dtos.ComponentDTO {
	for _, componentDTO := range components {
		componentDTO = h.handleOwner(componentDTO)
		result = append(result, componentDTO)
	}
	return result
}

func (h *ApplyHandler) handleCreated(
	result []*dtos.ComponentDTO,
	components map[string]*dtos.ComponentDTO,
	stateComponents map[string]*dtos.ComponentDTO,
) []*dtos.ComponentDTO {
	for _, componentDTO := range components {
		componentDTO = h.handleOwner(componentDTO)
		component := componentDTOToResource(componentDTO)

		component, errComponent := h.repository.Create(context.Background(), component)
		if errComponent != nil {
			panic(errComponent)
		}

		fmt.Printf("Deps: %+v \n", componentDTO.Spec.DependsOn)
		for _, providerName := range componentDTO.Spec.DependsOn {
			if provider, exists := stateComponents[providerName]; exists {
				h.repository.SetDependency(context.Background(), component.ID, provider.Spec.ID)
			} else {
				log.Printf("Provider %s not found for component %s", providerName, componentDTO.Spec.Name)
			}
		}

		componentDTO.Spec.ID = component.ID
		componentDTO.Spec.Slug = component.Slug
		for metricName, metricSource := range component.MetricSources {
			componentDTO.Spec.MetricSources[metricName] = &dtos.MetricSourceDTO{
				ID:     metricSource.ID,
				Name:   metricSource.Name,
				Metric: metricSource.Metric,
			}
		}
		result = append(result, componentDTO)
	}

	return result
}

func (h *ApplyHandler) handleUpdated(
	result []*dtos.ComponentDTO,
	components map[string]*dtos.ComponentDTO,
	stateComponents map[string]*dtos.ComponentDTO,
) []*dtos.ComponentDTO {
	for _, componentDTO := range components {
		componentDTO = h.handleOwner(componentDTO)
		component := componentDTOToResource(componentDTO)
		errComponent := h.repository.Update(context.Background(), component)
		if errComponent != nil {
			panic(errComponent)
		}

		h.handleDependencies(componentDTO, stateComponents)

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

	computedLinks := make([]dtos.Link, 0)
	for _, link := range componentDTO.Spec.Links {
		if link.Type != "CHAT_CHANNEL" {
			computedLinks = append(computedLinks, link)
		}
	}
	if owner.SlackChannel != "" {
		computedLinks = append(computedLinks, h.handleChatLinkFromOwner(owner))
	}

	componentDTO.Spec.Links = computedLinks
	componentDTO.Spec.OwnerID = owner.CompassID

	return componentDTO
}

func (h *ApplyHandler) handleChatLinkFromOwner(owner *ownerservicedtos.Owner) dtos.Link {
	if owner.DisplayName == "" {
		owner.DisplayName = "Slack"
	}

	return dtos.Link{
		Name: owner.DisplayName,
		Type: "CHAT_CHANNEL",
		URL:  owner.SlackChannel,
	}
}

func (h *ApplyHandler) handleDependencies(
	componentDTO *dtos.ComponentDTO,
	stateComponents map[string]*dtos.ComponentDTO,
) {
	componentInState := stateComponents[componentDTO.Metadata.Name]
	for _, providerName := range componentInState.Spec.DependsOn {
		if !listutils.Contains(componentDTO.Spec.DependsOn, providerName) {
			err := h.repository.UnSetDependency(context.Background(), componentDTO.Spec.ID, stateComponents[providerName].Spec.ID)
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

			err := h.repository.SetDependency(context.Background(), componentDTO.Spec.ID, stateProvider.Spec.ID)
			if err != nil {
				fmt.Printf("apply dependencies %s", err)
			}
		}
	}
}
