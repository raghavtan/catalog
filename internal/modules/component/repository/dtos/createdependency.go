package dtos

import (
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

/*************
 * INPUT DTO *
 *************/
type CreateDependencyInput struct {
	compassdtos.InputDTO
	DependentId string
	ProviderId  string
}

func (dto *CreateDependencyInput) GetQuery() string {
	return `
		mutation createRelationship($dependentId: ID!, $providerId: ID!) {
			compass {
				createRelationship(input: {
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

func (dto *CreateDependencyInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"dependentId": dto.DependentId,
		"providerId":  dto.ProviderId,
	}
}

/**************
 * OUTPUT DTO *
 **************/

type CreateDependencyOutput struct {
	Compass struct {
		CreateDependency struct {
			Errors  []compassservice.CompassError `json:"errors"`
			Success bool                          `json:"success"`
		} `json:"createRelationship"`
	} `json:"compass"`
}

func (dto *CreateDependencyOutput) IsSuccessful() bool {
	return dto.Compass.CreateDependency.Success
}

func (dto *CreateDependencyOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.CreateDependency.Errors))
	for i, err := range dto.Compass.CreateDependency.Errors {
		errors[i] = err.Message
	}
	return errors
}
