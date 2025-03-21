package dtos

import (
	"github.com/motain/of-catalog/internal/modules/metric/resources"
)

/*************
 * INPUT DTO *
 *************/
type SearchMetricsInput struct {
	CompassCloudID string
	Metric         resources.Metric
}

func (dto *SearchMetricsInput) GetQuery() string {
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

func (dto *SearchMetricsInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"cloudId": dto.CompassCloudID,
		"name":    dto.Metric.Name,
	}
}

/**************
 * OUTPUT DTO *
 **************/

type SearchMetricsOutput struct {
	Compass struct {
		Definitions struct {
			Nodes []Metric `json:"nodes"`
		} `json:"metricDefinitions"`
	} `json:"compass"`
}

func (c *SearchMetricsOutput) IsSuccessful() bool {
	return c.Compass.Definitions.Nodes != nil
}

func (dto *SearchMetricsOutput) GetErrors() []string {
	return nil
}
