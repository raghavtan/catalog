package dtos

import (
	"github.com/motain/of-catalog/internal/modules/component/resources"
)

type UpdateComponentOutput struct {
	Compass struct {
		UpdateComponent struct {
			Success bool `json:"success"`
		} `json:"updateComponent"`
	} `json:"compass"`
}

func (u *UpdateComponentOutput) GetQuery() string {
	return `
		mutation updateComponent ($componentDetails: UpdateCompassComponentInput!) {
			compass {
				updateComponent(input: $componentDetails) {
					success
					errors {
						message
					}
				}
			}
		}`
}

func (u *UpdateComponentOutput) SetVariables(component resources.Component) map[string]interface{} {
	variables := map[string]interface{}{
		"componentDetails": map[string]interface{}{
			"id":          component.ID,
			"name":        component.Name,
			"slug":        component.Slug,
			"description": component.Description,
		},
	}

	if component.OwnerID != "" {
		variables["componentDetails"].(map[string]interface{})["ownerId"] = component.OwnerID
	}

	return variables
}

func (c *UpdateComponentOutput) IsSuccessful() bool {
	return c.Compass.UpdateComponent.Success
}
