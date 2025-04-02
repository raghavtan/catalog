package dtos

import (
	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

/*************
 * INPUT DTO *
 *************/
type UpdateDocumentInput struct {
	compassdtos.InputDTO
	Document   resources.Document
	CategoryID string
}

func (dto *UpdateDocumentInput) GetQuery() string {
	return `
		mutation updateDocument($input: CompassUpdateDocumentInput!) {
		compass @optIn(to: "compass-beta") {
			updateDocument(input: $input) {
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

func (dto *UpdateDocumentInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"id":                      dto.Document.ID,
			"title":                   dto.Document.Title,
			"documentationCategoryId": dto.CategoryID,
			"url":                     dto.Document.URL,
		},
	}
}

/**************
 * OUTPUT DTO *
 **************/

type UpdateDocumentOutput struct {
	Compass struct {
		UpdateDocument struct {
			Errors  []compassservice.CompassError `json:"errors"`
			Success bool                          `json:"success"`
			Details Document                      `json:"documentDetails"`
		} `json:"updateDocument"`
	} `json:"compass"`
}

func (dto *UpdateDocumentOutput) IsSuccessful() bool {
	return dto.Compass.UpdateDocument.Success
}

func (dto *UpdateDocumentOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.UpdateDocument.Errors))
	for i, err := range dto.Compass.UpdateDocument.Errors {
		errors[i] = err.Message
	}
	return errors
}
