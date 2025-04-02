package dtos

import compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"

/*************
 * INPUT DTO *
 *************/
type DocumentationCategoriesInput struct {
	compassdtos.InputDTO
	CompassCloudID string
}

func (dto *DocumentationCategoriesInput) GetQuery() string {
	return `
		query documentationCategories {
			compass {
				documentationCategories(cloudId: "fca6a80f-888b-4079-82e6-3c2f61c788e2") @optIn(to: "compass-beta")  {
					... on CompassDocumentationCategoriesConnection {
						nodes {
							name
							id
							description
						}
					}
				}
			}
		}`
}

func (dto *DocumentationCategoriesInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"cloudId": dto.CompassCloudID,
	}
}

/**************
 * OUTPUT DTO *
 **************/

type DocumentationCategoriesOutput struct {
	Compass struct {
		DocumentationCategories struct {
			Nodes []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"nodes"`
		} `json:"documentationCategories"`
	} `json:"compass"`
}

func (dto *DocumentationCategoriesOutput) IsSuccessful() bool {
	return dto.Compass.DocumentationCategories.Nodes != nil
}

func (dto *DocumentationCategoriesOutput) GetErrors() []string {
	return nil
}
