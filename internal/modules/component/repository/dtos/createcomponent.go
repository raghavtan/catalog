package dtos

import (
	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
)

type Component struct {
	ID            string `json:"id"`
	Links         []Link `json:"links"`
	MetricSources struct {
		Nodes []MetricSource `json:"nodes"`
	} `json:"metricSources"`
}

type Link struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type MetricSource struct {
	ID               string `json:"id"`
	MetricDefinition struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"metricDefinition"`
}

type CreateComponentOutput struct {
	Compass struct {
		CreateComponent struct {
			Success bool                          `json:"success"`
			Errors  []compassservice.CompassError `json:"errors"`
			Details Component                     `json:"componentDetails"`
		} `json:"createComponent"`
	} `json:"compass"`
}

func (c *CreateComponentOutput) GetQuery() string {
	return `
		mutation createComponent ($cloudId: ID!, $componentDetails: CreateCompassComponentInput!) {
			compass {
				createComponent(cloudId: $cloudId, input: $componentDetails) {
					success
					componentDetails {
						id
						links {
							id
							type
							name
							url
						}
					}
					errors {
						message
					}
				}
			}
		}`
}

func (c *CreateComponentOutput) SetVariables(compassCloudIdD string, component resources.Component) map[string]interface{} {
	links := make([]map[string]string, 0)
	for _, link := range component.Links {
		links = append(links, map[string]string{
			"type": link.Type,
			"name": link.Name,
			"url":  link.URL,
		})
	}

	variables := map[string]interface{}{
		"cloudId": compassCloudIdD,
		"componentDetails": map[string]interface{}{
			"name":        component.Name,
			"slug":        component.Slug,
			"description": component.Description,
			"typeId":      component.TypeID,
			"links":       links,
			"labels":      component.Labels,
		},
	}

	if component.OwnerID != "" {
		variables["componentDetails"].(map[string]interface{})["ownerId"] = component.OwnerID
	}

	return variables
}

func (c *CreateComponentOutput) IsSuccessful() bool {
	return c.Compass.CreateComponent.Success
}
