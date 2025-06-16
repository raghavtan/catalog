package dtos

import (
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

/*************
 * INPUT DTO *
 *************/
type DeleteDocumentInput struct {
	compassdtos.InputDTO
	ID string
}

func (dto *DeleteDocumentInput) GetQuery() string {
	return `
		mutation deleteDocument($input: CompassDeleteDocumentInput!) {
   		compass @optIn(to: "compass-beta") {
   			deleteDocument(input: $input) {
   				success
 					errors {
						message
					}
					success
				}
			}
		}`
}

func (dto *DeleteDocumentInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"id": dto.ID,
		},
	}
}

/**************
 * OUTPUT DTO *
 **************/

type DeleteDocumentOutput struct {
	Compass struct {
		DeleteDocument struct {
			Errors  []compassservice.CompassError `json:"errors"`
			Success bool                          `json:"success"`
			Details Document                      `json:"documentDetails"`
		} `json:"deleteDocument"`
	} `json:"compass"`
}

func (dto *DeleteDocumentOutput) IsSuccessful() bool {
	return dto.Compass.DeleteDocument.Success
}

func (dto *DeleteDocumentOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.DeleteDocument.Errors))
	for i, err := range dto.Compass.DeleteDocument.Errors {
		errors[i] = err.Message
	}
	return errors
}
