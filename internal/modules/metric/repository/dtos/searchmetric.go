package dtos

import (
	"github.com/motain/of-catalog/internal/modules/metric/resources"
)

type SearchMetricsOutput struct {
	Compass struct {
		Definitions struct {
			Nodes []Metric `json:"nodes"`
		} `json:"metricDefinitions"`
	} `json:"compass"`
}

func (c *SearchMetricsOutput) GetQuery() string {
	return `
		query searchMetricDefinition($cloudId: ID!) {
			compass {
				metricDefinitions(query: {cloudId: $cloudId, first: 100}) {
					... on CompassMetricDefinitionsConnection {
						nodes{
							id
							name
						}
					}
				}
			}
		}`
}

func (c *SearchMetricsOutput) SetVariables(compassCloudIdD string, metric resources.Metric) map[string]interface{} {
	return map[string]interface{}{
		"cloudId": compassCloudIdD,
		"name":    metric.Name,
	}
}

func (c *SearchMetricsOutput) IsSuccessful() bool {
	return c.Compass.Definitions.Nodes != nil
}
