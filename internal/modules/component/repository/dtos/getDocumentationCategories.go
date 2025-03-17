package dtos

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

func (c *DocumentationCategoriesOutput) GetQuery() string {
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

func (c *DocumentationCategoriesOutput) SetVariables(metricID, componentID, identifier string) map[string]interface{} {
	return map[string]interface{}{}
}

func (c *DocumentationCategoriesOutput) IsSuccessful() bool {
	return c.Compass.DocumentationCategories.Nodes != nil
}
