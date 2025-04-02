package dtos

import (
	"github.com/motain/of-catalog/internal/modules/metric/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

/*************
 * INPUT DTO *
 *************/
type UpdateMetricInput struct {
	compassdtos.InputDTO
	CompassCloudID string
	Metric         resources.Metric
}

func (dto *UpdateMetricInput) GetQuery() string {
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

func (dto *UpdateMetricInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"cloudId":     dto.CompassCloudID,
		"id":          dto.Metric.ID,
		"name":        dto.Metric.Name,
		"description": dto.Metric.Description,
		"unit":        dto.Metric.Format.Unit,
	}
}

/**************
 * OUTPUT DTO *
 **************/

type UpdateMetricOutput struct {
	Compass struct {
		UpdateMetric struct {
			Errors  []compassservice.CompassError `json:"errors"`
			Success bool                          `json:"success"`
		} `json:"updateMetricDefinition"`
	} `json:"compass"`
}

func (dto *UpdateMetricOutput) IsSuccessful() bool {
	return dto.Compass.UpdateMetric.Success
}

func (dto *UpdateMetricOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.UpdateMetric.Errors))
	for i, err := range dto.Compass.UpdateMetric.Errors {
		errors[i] = err.Message
	}
	return errors
}
