package dtos

import compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"

/*************
 * INPUT DTO *
 *************/
type GetDocumentsInput struct {
	compassdtos.InputDTO
	ComponentID string
}

func (dto *GetDocumentsInput) GetQuery() string {
	return `
		query documents($componentId: ID!) {
			compass {
				documents(componentId: $componentId, first: 100) @optIn(to: "compass-beta")  {
					... on CompassDocumentConnection {
						nodes {
							title
							id
							url
							documentationCategoryId
						}
					}
				}
			}
		}`
}

func (dto *GetDocumentsInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"componentId": dto.ComponentID,
	}
}

/**************
 * OUTPUT DTO *
 **************/

type GetDocumentsOutput struct {
	Compass struct {
		Documents struct {
			Nodes []struct {
				ID                      string `json:"id"`
				Title                   string `json:"title"`
				URL                     string `json:"url"`
				DocumentationCategoryID string `json:"documentationCategoryId,omitempty"`
			} `json:"nodes"`
		} `json:"documents"`
	} `json:"compass"`
}

func (dto *GetDocumentsOutput) IsSuccessful() bool {
	return true
}

func (dto *GetDocumentsOutput) GetErrors() []string {
	return nil
}
