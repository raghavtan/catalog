package dtos

import (
	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

/*************
 * INPUT DTO *
 *************/
type CreateDocumentInput struct {
	compassdtos.InputDTO
	ComponentID string
	Document    resources.Document
	CategoryID  string
}

func (dto *CreateDocumentInput) GetQuery() string {
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

func (dto *CreateDocumentInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"componentId":             dto.ComponentID,
			"title":                   dto.Document.Title,
			"documentationCategoryId": dto.CategoryID,
			"url":                     dto.Document.URL,
		},
	}
}

/**************
 * OUTPUT DTO *
 **************/

type Document struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Type  string `json:"type"`
}

type CreateDocumentOutput struct {
	Compass struct {
		AddDocument struct {
			Errors  []compassservice.CompassError `json:"errors"`
			Success bool                          `json:"success"`
			Details Document                      `json:"documentDetails"`
		} `json:"addDocument"`
	} `json:"compass"`
}

func (dto *CreateDocumentOutput) IsSuccessful() bool {
	return dto.Compass.AddDocument.Success
}

func (dto *CreateDocumentOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.AddDocument.Errors))
	for i, err := range dto.Compass.AddDocument.Errors {
		errors[i] = err.Message
	}
	return errors
}
