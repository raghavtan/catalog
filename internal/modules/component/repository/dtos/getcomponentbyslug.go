package dtos

type ComponentByReferenceOutput struct {
	Compass struct {
		Component Component `json:"componentByReference"`
	} `json:"compass"`
}

func (c *ComponentByReferenceOutput) GetQuery() string {
	return `
		query getComponentBySlug($cloudId: ID!, $slug: String!) {
			compass {
				componentByReference(reference: {slug: {slug: $slug, cloudId: $cloudId}}) {
					... on CompassComponent {
						id
						metricSources {
							... on CompassComponentMetricSourcesConnection {
								nodes {
									id,
									metricDefinition {
										name
									}
								}
							}
						}
					}
				}
			}
		}`
}

func (c *ComponentByReferenceOutput) SetVariables(compassCloudIdD, slug string) map[string]interface{} {
	return map[string]interface{}{
		"cloudId": compassCloudIdD,
		"slug":    slug,
	}
}

func (c *ComponentByReferenceOutput) IsSuccessful() bool {
	return c.Compass.Component.ID != ""
}
