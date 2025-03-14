package dtos

type Document struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Type  string `json:"type"`
}

type CreateDocumentOutput struct {
	Compass struct {
		AddDocument struct {
			Success bool     `json:"success"`
			Details Document `json:"documentDetails"`
		} `json:"addDocument"`
	} `json:"compass"`
}

func (c *CreateDocumentOutput) GetQuery() string {
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

func (c *CreateDocumentOutput) SetVariables(componentID, title, categoryID, url string) map[string]interface{} {
	return map[string]interface{}{
		"input": map[string]interface{}{
			"componentId":             componentID,
			"title":                   title,
			"documentationCategoryId": categoryID,
			"url":                     url,
		},
	}
}

func (c *CreateDocumentOutput) IsSuccessful() bool {
	return c.Compass.AddDocument.Success
}
