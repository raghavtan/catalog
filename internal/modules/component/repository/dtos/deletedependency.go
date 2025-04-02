package dtos

import (
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

/*************
 * INPUT DTO *
 *************/
type DeleteDependencyInput struct {
	compassdtos.InputDTO
	DependentId string
	ProviderId  string
}

func (dto *DeleteDependencyInput) GetQuery() string {
	return `
		mutation deleteRelationship($dependentId: ID!, $providerId: ID!) {
			compass {
				deleteRelationship(input: {
					type: DEPENDS_ON,
					startNodeId: $dependentId,
					endNodeId: $providerId
				}) {
					errors {
						message
					}
					success
				}
			}
		}`
}

func (dto *DeleteDependencyInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"dependentId": dto.DependentId,
		"providerId":  dto.ProviderId,
	}
}

/**************
 * OUTPUT DTO *
 **************/

type DeleteDependencyOutput struct {
	Compass struct {
		DeleteDependency struct {
			Errors  []compassservice.CompassError `json:"errors"`
			Success bool                          `json:"success"`
		} `json:"deleteRelationship"`
	} `json:"compass"`
}

func (dto *DeleteDependencyOutput) IsSuccessful() bool {
	return dto.Compass.DeleteDependency.Success
}

func (dto *DeleteDependencyOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.DeleteDependency.Errors))
	for i, err := range dto.Compass.DeleteDependency.Errors {
		errors[i] = err.Message
	}
	return errors
}
