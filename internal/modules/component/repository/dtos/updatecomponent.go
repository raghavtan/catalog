package dtos

import (
	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

/*************
 * INPUT DTO *
 *************/
type UpdateComponentInput struct {
	compassdtos.InputDTO
	Component resources.Component
}

func (dto *UpdateComponentInput) GetQuery() string {
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

func (dto *UpdateComponentInput) SetVariables() map[string]interface{} {
	variables := map[string]interface{}{
		"componentDetails": map[string]interface{}{
			"id":          dto.Component.ID,
			"name":        dto.Component.Name,
			"slug":        dto.Component.Slug,
			"description": dto.Component.Description,
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

type UpdateComponentOutput struct {
	Compass struct {
		UpdateComponent struct {
			Errors  []compassservice.CompassError `json:"errors"`
			Success bool                          `json:"success"`
		} `json:"updateComponent"`
	} `json:"compass"`
}

func (dto *UpdateComponentOutput) IsSuccessful() bool {
	return dto.Compass.UpdateComponent.Success
}

func (dto *UpdateComponentOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.UpdateComponent.Errors))
	for i, err := range dto.Compass.UpdateComponent.Errors {
		errors[i] = err.Message
	}
	return errors
}
