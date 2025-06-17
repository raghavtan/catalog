package dtos

import (
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

/*************
 * INPUT DTO *
 *************/
type RemoveLinkInput struct {
	compassdtos.InputDTO
	ComponentID string
	LinkID      string
}

func (dto *RemoveLinkInput) GetQuery() string {
	return `
	mutation deleteComponentLink($input: DeleteCompassComponentLinkInput!) {
   		compass {
   			deleteComponentLink(input: $input) {
   				success
				errors {
					message
				}
			}
		}
	}`
}

func (dto *RemoveLinkInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"componentId": dto.ComponentID,
			"link":        dto.LinkID,
		},
	}
}

type RemoveLinkOutput struct {
	Compass struct {
		DeleteComponentLink struct {
			Errors  []compassservice.CompassError `json:"errors"`
			Success bool                          `json:"success"`
		} `json:"deletedComponentLink"`
	} `json:"compass"`
}

func (dto *RemoveLinkOutput) IsSuccessful() bool {
	return dto.Compass.DeleteComponentLink.Success
}

func (dto *RemoveLinkOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.DeleteComponentLink.Errors))
	for i, err := range dto.Compass.DeleteComponentLink.Errors {
		errors[i] = err.Message
	}
	return errors
}
