package dtos

import (
	"fmt"

	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

/*************
 * INPUT DTO *
 *************/
type CreateComponentInput struct {
	compassdtos.InputDTO
	CompassCloudID string
	Component      resources.Component
}

func (dto *CreateComponentInput) GetQuery() string {
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

func (dto *CreateComponentInput) SetVariables() map[string]interface{} {
	links := make([]map[string]string, 0)
	for _, link := range dto.Component.Links {
		links = append(links, map[string]string{
			"type": link.Type,
			"name": link.Name,
			"url":  link.URL,
		})
	}

	fields := make([]map[string]interface{}, 0, len(dto.Component.Fields))
	for k, v := range dto.Component.Fields {
		var valueObj map[string]interface{}
		switch v := v.(type) {
		case bool:
			valueObj = map[string]interface{}{
				"boolean": map[string]bool{"booleanValue": v},
			}
		default:
			valueObj = map[string]interface{}{
				"enum": map[string][]string{"value": {fmt.Sprintf("%v", v)}},
			}
		}

		fields = append(fields, map[string]interface{}{
			"definition": "compass:" + k,
			"value":      valueObj,
		})
	}

	variables := map[string]interface{}{
		"cloudId": dto.CompassCloudID,
		"componentDetails": map[string]interface{}{
			"name":        dto.Component.Name,
			"slug":        dto.Component.Slug,
			"description": dto.Component.Description,
			"typeId":      dto.Component.TypeID,
			"fields":      fields,
			"links":       links,
			"labels":      dto.Component.Labels,
		},
	}

	if dto.Component.OwnerID != "" {
		variables["componentDetails"].(map[string]interface{})["ownerId"] = dto.Component.OwnerID
	}

	return variables
}

/**************
 * OUTPUT DTO *
 **************/

type MetricSources struct {
	Nodes []MetricSource `json:"nodes"`
}

type Component struct {
	ID            string        `json:"id"`
	Links         []Link        `json:"links"`
	MetricSources MetricSources `json:"metricSources"`
}

type Link struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type MetricDefinition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type MetricSource struct {
	ID               string           `json:"id"`
	MetricDefinition MetricDefinition `json:"metricDefinition"`
}

type CompassCreateComponentOutput struct {
	Success bool                          `json:"success"`
	Errors  []compassservice.CompassError `json:"errors"`
	Details Component                     `json:"componentDetails"`
}

type CompassCreatedComponentOutput struct {
	CreateComponent CompassCreateComponentOutput `json:"createComponent"`
}
type CreateComponentOutput struct {
	Compass CompassCreatedComponentOutput `json:"compass"`
}

func (dto *CreateComponentOutput) IsSuccessful() bool {
	return dto.Compass.CreateComponent.Success
}

func (dto *CreateComponentOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.CreateComponent.Errors))
	for i, err := range dto.Compass.CreateComponent.Errors {
		errors[i] = err.Message
	}
	return errors
}
