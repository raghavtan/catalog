package dtos

import (
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

/*************
 * INPUT DTO *
 *************/
type BindMetricInput struct {
	compassdtos.InputDTO
	MetricID    string
	ComponentID string
	Identifier  string
}

func (dto *BindMetricInput) GetQuery() string {
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

func (dto *BindMetricInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"metricId":    dto.MetricID,
		"componentId": dto.ComponentID,
		"externalId":  dto.Identifier,
	}
}

/**************
 * OUTPUT DTO *
 **************/

type BindMetricOutput struct {
	Compass struct {
		CreateMetricSource struct {
			Errors             []compassservice.CompassError `json:"errors"`
			Success            bool                          `json:"success"`
			CreateMetricSource struct {
				ID string `json:"id"`
			} `json:"createdMetricSource"`
		} `json:"createMetricSource"`
	} `json:"compass"`
}

func (dto *BindMetricOutput) IsSuccessful() bool {
	return dto.Compass.CreateMetricSource.Success
}

func (dto *BindMetricOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.CreateMetricSource.Errors))
	for i, err := range dto.Compass.CreateMetricSource.Errors {
		errors[i] = err.Message
	}
	return errors
}
