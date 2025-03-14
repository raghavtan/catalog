package dtos

import (
	"github.com/motain/of-catalog/internal/services/compassservice"
)

type DeleteDependencyOutput struct {
	Compass struct {
		DeleteDependency struct {
			Errors  []compassservice.CompassError `json:"errors"`
			Success bool                          `json:"success"`
		} `json:"deleteRelationship"`
	} `json:"compass"`
}

func (c *DeleteDependencyOutput) GetQuery() string {
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

func (c *DeleteDependencyOutput) SetVariables(dependentId, providerId string) map[string]interface{} {
	return map[string]interface{}{
		"dependentId": dependentId,
		"providerId":  providerId,
	}
}

func (c *DeleteDependencyOutput) IsSuccessful() bool {
	return c.Compass.DeleteDependency.Success
}
