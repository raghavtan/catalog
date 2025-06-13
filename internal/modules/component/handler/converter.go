package handler

import (
	"fmt"
	"log"

	"github.com/motain/of-catalog/internal/modules/component/dtos"
	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/modules/component/utils"
	"github.com/motain/of-catalog/internal/services/githubservice"
)

// ComponentConverter handles conversion from DTO to Resource with optional GitHub integration
type ComponentConverter struct {
	github githubservice.GitHubServiceInterface
}

// NewComponentConverter creates a new converter with GitHub service
func NewComponentConverter(github githubservice.GitHubServiceInterface) *ComponentConverter {
	return &ComponentConverter{github: github}
}

// ToResource converts ComponentDTO to Resource with description handling
func (c *ComponentConverter) ToResource(componentDTO *dtos.ComponentDTO) resources.Component {
	description := c.getDescription(componentDTO)

	return resources.Component{
		ID:            componentDTO.Spec.ID,
		Name:          componentDTO.Spec.Name,
		Slug:          utils.GetSlug(componentDTO.Spec.Name, componentDTO.Spec.TypeID),
		Description:   description,
		ConfigVersion: componentDTO.Spec.ConfigVersion,
		TypeID:        componentDTO.Spec.TypeID,
		OwnerID:       componentDTO.Spec.OwnerID,
		Fields:        componentDTO.Spec.Fields,
		Links:         linksDTOToResource(componentDTO.Spec.Links),
		Labels:        componentDTO.Spec.Labels,
		MetricSources: metricSourcesDTOToResource(componentDTO.Spec.MetricSources),
	}
}

// getDescription handles description logic with GitHub fallback
func (c *ComponentConverter) getDescription(componentDTO *dtos.ComponentDTO) string {
	// If description is already set, use it
	if componentDTO.Spec.Description != "" {
		return componentDTO.Spec.Description
	}

	// If no GitHub service available, use default
	if c.github == nil {
		return fmt.Sprintf("Component %s", componentDTO.Spec.Name)
	}

	// Try to get description from GitHub
	description, err := c.github.GetRepoDescription(componentDTO.Metadata.Name)
	if err != nil {
		log.Printf("Warning: Could not get repository description for %s: %v", componentDTO.Metadata.Name, err)
		return fmt.Sprintf("Component %s", componentDTO.Spec.Name)
	}

	log.Printf("INFO: Using GitHub description for %s: %s", componentDTO.Spec.Name, description)
	return description
}

// Backward compatibility functions
func componentDTOToResource(componentDTO *dtos.ComponentDTO) resources.Component {
	// Create converter without GitHub service for backward compatibility
	converter := &ComponentConverter{github: nil}
	return converter.ToResource(componentDTO)
}

func componentDTOToResourceWithGitHub(componentDTO *dtos.ComponentDTO, github githubservice.GitHubServiceInterface) resources.Component {
	converter := &ComponentConverter{github: github}
	return converter.ToResource(componentDTO)
}

// MetricSourceDTOToResource converts MetricSourceDTO to Resource
// Used by BindHandler and ComputeHandler
func MetricSourceDTOToResource(metricSourceDTO *dtos.MetricSourceDTO) resources.MetricSource {
	return resources.MetricSource{
		ID:     metricSourceDTO.ID,
		Name:   metricSourceDTO.Name,
		Metric: metricSourceDTO.Metric,
		Facts:  metricSourceDTO.Facts, // Include Facts for ComputeHandler
	}
}

// Helper functions remain the same
func linksDTOToResource(linksDTO []dtos.Link) []resources.Link {
	uniqueLinks := make(map[string]resources.Link)

	for _, link := range linksDTO {
		uniqueKey := fmt.Sprintf("%s-%s-%s", link.Name, link.Type, link.URL)
		if _, exists := uniqueLinks[uniqueKey]; !exists {
			uniqueLinks[uniqueKey] = resources.Link{
				Name: link.Name,
				Type: link.Type,
				URL:  link.URL,
			}
		}
	}

	links := make([]resources.Link, 0, len(uniqueLinks))
	for _, link := range uniqueLinks {
		links = append(links, link)
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
