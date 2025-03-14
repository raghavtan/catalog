package dtos

import (
	"github.com/motain/of-catalog/internal/modules/metric/resources"
)

type UpdateMetricOutput struct {
	Compass struct {
		UpdateMetric struct {
			Success bool `json:"success"`
		} `json:"updateMetricDefinition"`
	} `json:"compass"`
}

func (u *UpdateMetricOutput) GetQuery() string {
	return `
		mutation updateMetricDefinition ($cloudId: ID!, $id: ID!, $name: String!, $description: String!, $unit: String!) {
			compass {
				updateMetricDefinition(
					input: {
						id: $id
						cloudId: $cloudId
						name: $name
						description: $description
						format: {
							suffix: { suffix: $unit }
						}
					}
				) {
					success
					errors {
						message
					}
				}
			}
		}`
}

func (u *UpdateMetricOutput) SetVariables(compassCloudIdD string, metric resources.Metric) map[string]interface{} {
	return map[string]interface{}{
		"cloudId":     compassCloudIdD,
		"id":          metric.ID,
		"name":        metric.Name,
		"description": metric.Description,
		"unit":        metric.Format.Unit,
	}
}

func (c *UpdateMetricOutput) IsSuccessful() bool {
	return c.Compass.UpdateMetric.Success
}
