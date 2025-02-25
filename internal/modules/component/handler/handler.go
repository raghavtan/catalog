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

type Handler struct {
	github     githubservice.GitHubRepositoriesServiceInterface
	repository repository.RepositoryInterface
}

func NewHandler(
	gh githubservice.GitHubRepositoriesServiceInterface,
	repository repository.RepositoryInterface,
) *Handler {
	return &Handler{github: gh, repository: repository}
}

func (h *Handler) Apply() string {
	configComponents, errConfig := yaml.ParseConfig[dtos.ComponentDTO]()
	if errConfig != nil {
		log.Fatalf("error: %v", errConfig)
	}

	stateComponents, errState := yaml.ParseState[dtos.ComponentDTO]()
	if errState != nil {
		log.Fatalf("error: %v", errState)
	}

	getUniqueKey := func(c *dtos.ComponentDTO) string {
		return c.Spec.Name
	}
	setID := func(c *dtos.ComponentDTO, id string) {
		c.Spec.ID = &id
	}
	getID := func(c *dtos.ComponentDTO) string {
		return *c.Spec.ID
	}

	isEqualLinks := func(l1, l2 []dtos.Link) bool {
		for i, link := range l1 {
			if link.Name != l2[i].Name || link.Type != l2[i].Type || link.URL != l2[i].URL {
				return false
			}
		}
		return true
	}

	isEqualLabels := func(l1, l2 []string) bool {
		if len(l1) != len(l2) {
			return false
		}
		for i, label := range l1 {
			if label != l2[i] {
				return false
			}
		}
		return true
	}

	isEqual := func(c1, c2 *dtos.ComponentDTO) bool {
		return c1.Spec.Name == c2.Spec.Name &&
			c1.Spec.Description == c2.Spec.Description &&
			c1.Spec.ConfigVersion == c2.Spec.ConfigVersion &&
			c1.Spec.TypeID == c2.Spec.TypeID &&
			c1.Spec.OwnerID == c2.Spec.OwnerID &&
			isEqualLinks(c1.Spec.Links, c2.Spec.Links) &&
			isEqualLabels(c1.Spec.Labels, c2.Spec.Labels)
	}

	created, updated, deleted, unchanged := drift.Detect(
		stateComponents,
		configComponents,
		getUniqueKey,
		getID,
		setID,
		isEqual,
	)
	h.handleDeleted(deleted)

	var result = unchanged
	result = h.handleCreated(result, created)
	result = h.handleUpdated(result, updated)

	err := yaml.WriteState[dtos.ComponentDTO](result)
	if err != nil {
		log.Fatalf("error writing components to file: %v", err)
	}

	return ""
}

func (h *Handler) handleDeleted(components []*dtos.ComponentDTO) {
	for _, componentDTO := range components {
		errComponent := h.repository.Delete(context.Background(), *componentDTO.Spec.ID)
		if errComponent != nil {
			panic(errComponent)
		}
	}
}

func (h *Handler) handleCreated(result, components []*dtos.ComponentDTO) []*dtos.ComponentDTO {
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

func (h *Handler) handleUpdated(result, components []*dtos.ComponentDTO) []*dtos.ComponentDTO {
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
