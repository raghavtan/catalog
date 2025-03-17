package dtos

import "github.com/motain/of-catalog/internal/services/compassservice"

type DeleteComponentOutput struct {
	Compass struct {
		DeleteComponent struct {
			Errors  []compassservice.CompassError `json:"errors"`
			Success bool                          `json:"success"`
		} `json:"deleteComponent"`
	} `json:"compass"`
}

func (d *DeleteComponentOutput) GetQuery() string {
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

func (d *DeleteComponentOutput) SetVariables(id string) map[string]interface{} {
	return map[string]interface{}{
		"id": id,
	}
}

func (c *DeleteComponentOutput) IsSuccessful() bool {
	// Ignoring the error if the component is not found
	if compassservice.HasNotFoundError(c.Compass.DeleteComponent.Errors) {
		return true
	}

	return c.Compass.DeleteComponent.Success
}
