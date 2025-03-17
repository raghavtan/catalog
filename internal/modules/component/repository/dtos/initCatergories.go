package dtos

type BindMetricOutput struct {
	Compass struct {
		CreateMetricSource struct {
			Success            bool `json:"success"`
			CreateMetricSource struct {
				ID string `json:"id"`
			} `json:"createdMetricSource"`
		} `json:"createMetricSource"`
	} `json:"compass"`
}

func (c *BindMetricOutput) GetQuery() string {
	return `
		mutation createMetricSource($metricId: ID!, $componentId: ID!, $externalId: ID!) {
			compass {
				createMetricSource(input: {metricDefinitionId: $metricId, componentId: $componentId, externalMetricSourceId: $externalId}) {
					success
					createdMetricSource {
						id
					}
					errors {
						message
					}
				}
			}
		}`
}

func (c *BindMetricOutput) SetVariables(metricID, componentID, identifier string) map[string]interface{} {
	return map[string]interface{}{
		"metricId":    metricID,
		"componentId": componentID,
		"externalId":  identifier,
	}
}

func (c *BindMetricOutput) IsSuccessful() bool {
	return c.Compass.CreateMetricSource.Success
}
