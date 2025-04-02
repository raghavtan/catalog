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
	ComponentID string
	Title       string
	CategoryID  string
	URL         string
}

func (dto *DeleteDocumentInput) GetQuery() string {
	return `
		mutation addDocument($input: CompassAddDocumentInput!) {
   		compass @optIn(to: "compass-beta") {
   			addDocument(input: $input) {
   				success
 					errors {
						message
					}
					documentDetails {
						id
						title
						url
						componentId
						documentationCategoryId
					}
				}
			}
		}`
}

func (dto *DeleteDocumentInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"componentId":             dto.ComponentID,
			"title":                   dto.Title,
			"documentationCategoryId": dto.CategoryID,
			"url":                     dto.URL,
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
		} `json:"addDocument"`
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
