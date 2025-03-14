package dtos

type DeleteDocumentOutput struct {
	Compass struct {
		DeleteDocument struct {
			Success bool     `json:"success"`
			Details Document `json:"documentDetails"`
		} `json:"addDocument"`
	} `json:"compass"`
}

func (c *DeleteDocumentOutput) GetQuery() string {
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

func (c *DeleteDocumentOutput) SetVariables(componentID, title, categoryID, url string) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"componentId":             componentID,
			"title":                   title,
			"documentationCategoryId": categoryID,
			"url":                     url,
		},
	}
}

func (c *DeleteDocumentOutput) IsSuccessful() bool {
	return c.Compass.DeleteDocument.Success
}
