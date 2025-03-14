package dtos

type UpdateDocumentOutput struct {
	Compass struct {
		UpdateDocument struct {
			Success bool     `json:"success"`
			Details Document `json:"documentDetails"`
		} `json:"updateDocument"`
	} `json:"compass"`
}

func (c *UpdateDocumentOutput) GetQuery() string {
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

func (c *UpdateDocumentOutput) SetVariables(documentID, title, categoryID, url string) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"id":                      documentID,
			"title":                   title,
			"documentationCategoryId": categoryID,
			"url":                     url,
		},
	}
}

func (c *UpdateDocumentOutput) IsSuccessful() bool {
	return c.Compass.UpdateDocument.Success
}
