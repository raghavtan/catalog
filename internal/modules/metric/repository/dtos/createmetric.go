package dtos

import (
	"github.com/motain/of-catalog/internal/modules/metric/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

/*************
 * INPUT DTO *
 *************/
type CreateMetricInput struct {
	compassdtos.InputDTO
	CompassCloudID string
	Metric         resources.Metric
}

func (dto *CreateMetricInput) GetQuery() string {
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

func (dto *CreateMetricInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"cloudId":     dto.CompassCloudID,
		"name":        dto.Metric.Name,
		"description": dto.Metric.Description,
		"unit":        dto.Metric.Format.Unit,
	}
}

/**************
 * OUTPUT DTO *
 **************/

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

func (dto *CreateMetricOutput) IsSuccessful() bool {
	return dto.Compass.CreateMetric.Success
}

func (dto *CreateMetricOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.CreateMetric.Errors))
	for i, err := range dto.Compass.CreateMetric.Errors {
		errors[i] = err.Message
	}
	return errors
}
