package dtos

import (
	"github.com/motain/of-catalog/internal/modules/metric/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
)

type Metric struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateMetricOutput struct {
	Compass struct {
		CreateMetric struct {
			Success    bool                          `json:"success"`
			Errors     []compassservice.CompassError `json:"errors"`
			Definition Metric                        `json:"createdMetricDefinition"`
		} `json:"createMetricDefinition"`
	} `json:"compass"`
}

func (c *CreateMetricOutput) GetQuery() string {
	return `
		mutation createMetricDefinition ($cloudId: ID!, $name: String!, $description: String!, $unit: String!) {
			compass {
				createMetricDefinition(
					input: {
						cloudId: $cloudId
						name: $name
						description: $description
						format: {
							suffix: { suffix: $unit }
						}
					}
				) {
					success
					createdMetricDefinition {
						id
					}
					errors {
						message
					}
				}
			}
		}`
}

func (c *CreateMetricOutput) SetVariables(compassCloudIdD string, metric resources.Metric) map[string]interface{} {
	return map[string]interface{}{
		"cloudId":     compassCloudIdD,
		"name":        metric.Name,
		"description": metric.Description,
		"unit":        metric.Format.Unit,
	}
}

func (c *CreateMetricOutput) IsSuccessful() bool {
	return c.Compass.CreateMetric.Success
}
