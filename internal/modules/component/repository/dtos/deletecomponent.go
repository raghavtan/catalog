package dtos

import (
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

/*************
 * INPUT DTO *
 *************/
type DeleteComponentInput struct {
	compassdtos.InputDTO
	ComponentID string
}

func (dto *DeleteComponentInput) GetQuery() string {
	return `
		mutation deleteComponent($id: ID!) {
			compass {
				deleteComponent(input: {id: $id}) {
					deletedComponentId
					errors {
						message
					}
					success
				}
			}
		}`
}

func (dto *DeleteComponentInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"id": dto.ComponentID,
	}
}

/**************
 * OUTPUT DTO *
 **************/

type DeleteComponentOutput struct {
	Compass struct {
		DeleteComponent struct {
			Errors  []compassservice.CompassError `json:"errors"`
			Success bool                          `json:"success"`
		} `json:"deleteComponent"`
	} `json:"compass"`
}

func (dto *DeleteComponentOutput) IsSuccessful() bool {
	// Ignoring the error if the component is not found
	if compassservice.HasNotFoundError(dto.Compass.DeleteComponent.Errors) {
		return true
	}

	return dto.Compass.DeleteComponent.Success
}

func (dto *DeleteComponentOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.DeleteComponent.Errors))
	for i, err := range dto.Compass.DeleteComponent.Errors {
		errors[i] = err.Message
	}
	return errors
}
