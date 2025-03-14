package dtos

import (
	"github.com/motain/of-catalog/internal/services/compassservice"
)

type CreateDependencyOutput struct {
	Compass struct {
		CreateDependency struct {
			Errors  []compassservice.CompassError `json:"errors"`
			Success bool                          `json:"success"`
		} `json:"createRelationship"`
	} `json:"compass"`
}

func (c *CreateDependencyOutput) GetQuery() string {
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

func (c *CreateDependencyOutput) SetVariables(dependentId, providerId string) map[string]interface{} {
	return map[string]interface{}{
		"dependentId": dependentId,
		"providerId":  providerId,
	}
}

func (c *CreateDependencyOutput) IsSuccessful() bool {
	return c.Compass.CreateDependency.Success
}
