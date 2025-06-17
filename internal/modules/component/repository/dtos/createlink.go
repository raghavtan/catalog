package dtos

import (
	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

/*************
 * INPUT DTO *
 *************/
type CreateLinkInput struct {
	compassdtos.InputDTO
	ComponentID string
	Link        resources.Link
}

func (dto *CreateLinkInput) GetQuery() string {
	return `
	mutation createComponentLink($input: CreateCompassComponentLinkInput!) {
   		compass {
   			createComponentLink(input: $input) {
   				success
				errors {
					message
				}
				createdComponentLink {
					id
				}
			}
		}
	}`
}

func (dto *CreateLinkInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"componentId": dto.ComponentID,
			"link": map[string]interface{}{
				"name": dto.Link.Name,
				"type": dto.Link.Type,
				"url":  dto.Link.URL,
			},
		},
	}
}

/**************
 * OUTPUT DTO *
 **************/
type LinkDetails struct {
	ID string `json:"id"`
}

type CreateLinkOutput struct {
	Compass struct {
		CreateComponentLink struct {
			Errors  []compassservice.CompassError `json:"errors"`
			Success bool                          `json:"success"`
			Details LinkDetails                   `json:"createdComponentLink"`
		} `json:"createComponentLink"`
	} `json:"compass"`
}

func (dto *CreateLinkOutput) IsSuccessful() bool {
	return dto.Compass.CreateComponentLink.Success
}

func (dto *CreateLinkOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.CreateComponentLink.Errors))
	for i, err := range dto.Compass.CreateComponentLink.Errors {
		errors[i] = err.Message
	}
	return errors
}
