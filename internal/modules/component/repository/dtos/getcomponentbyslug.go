package dtos

import compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"

/*************
 * INPUT DTO *
 *************/
type ComponentByReferenceInput struct {
	compassdtos.InputDTO
	CompassCloudID string
	Slug           string
}

func (dto *ComponentByReferenceInput) GetQuery() string {
	return `
		query getComponentBySlug($cloudId: ID!, $slug: String!) {
			compass {
				componentByReference(reference: {slug: {slug: $slug, cloudId: $cloudId}}) {
					... on CompassComponent {
						id
						name
						description
						ownerId
						labels
						type {
							id
						}
						links {
							id
							name
							type
							url
						}
						documents {
							nodes {
								id
								title
								url
							}
						}
						metricSources {
							nodes {
								id
								metricDefinition {
									id
									name
								}
							}
						}
					}
				}
			}
		}`
}

func (dto *ComponentByReferenceInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"cloudId": dto.CompassCloudID,
		"slug":    dto.Slug,
	}
}

/**************
 * OUTPUT DTO *
 **************/

// Define the structure for a Link based on the GraphQL query
type CompassLink struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

// Define the structure for a Document node based on the GraphQL query
type CompassDocumentNode struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	// Category CompassDocumentCategory `json:"category"` // If category is fetched
}

// Define the structure for Document connection (if paginated)
type CompassDocumentConnection struct {
	Nodes []CompassDocumentNode `json:"nodes"`
}

// Define the structure for Component Type based on the GraphQL query
type CompassComponentType struct {
	ID string `json:"id"`
	// Name string `json:"name"` // If name is fetched
}

// Define the structure for MetricDefinition based on the GraphQL query
type CompassMetricDefinition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Define the structure for a MetricSource node based on the GraphQL query
type CompassMetricSourceNode struct {
	ID               string                  `json:"id"`
	MetricDefinition CompassMetricDefinition `json:"metricDefinition"`
}

// Define the structure for MetricSource connection
type CompassMetricSourceConnection struct {
	Nodes []CompassMetricSourceNode `json:"nodes"`
}

// This is the main Component structure within the output DTO
type ComponentDetail struct {
	ID            string                        `json:"id"`
	Name          string                        `json:"name"`
	Description   string                        `json:"description"`
	OwnerID       string                        `json:"ownerId"` // Adjust if schema is different
	Labels        []string                      `json:"labels"`
	Type          CompassComponentType          `json:"type"`
	Links         []CompassLink                 `json:"links"`
	Documents     CompassDocumentConnection     `json:"documents"`
	MetricSources CompassMetricSourceConnection `json:"metricSources"`
}

type ComponentByReferenceOutput struct {
	Compass struct {
		Component ComponentDetail `json:"componentByReference"`
	} `json:"compass"`
}

func (dto *ComponentByReferenceOutput) IsSuccessful() bool {
	return dto.Compass.Component.ID != ""
}

func (dto *ComponentByReferenceOutput) GetErrors() []string {
	// Basic error check, can be expanded if GraphQL returns specific errors in a structured way
	// For now, assuming IsSuccessful is the primary check.
	return nil
}
