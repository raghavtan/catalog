package dtos

import (
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

/*************
 * INPUT DTO *
 *************/
type DeleteMetricInput struct {
	compassdtos.InputDTO
	CompassCloudID string
	MetricID       string
}

func (dto *DeleteMetricInput) GetQuery() string {
	return `
		mutation deleteMetric($scorecardId: ID!) {
			compass {
				deleteMetric(scorecardId: $scorecardId) {
					scorecardId
					errors {
						message
					}
					success
				}
			}
		}`
}

func (dto *DeleteMetricInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"id": dto.MetricID,
	}
}

/**************
 * OUTPUT DTO *
 **************/

type DeleteMetricOutput struct {
	Compass struct {
		DeleteMetric struct {
			Errors  []compassservice.CompassError `json:"errors"`
			Success bool                          `json:"success"`
		} `json:"deleteMetricDefinition"`
	} `json:"compass"`
}

func (dto *DeleteMetricOutput) IsSuccessful() bool {
	return dto.Compass.DeleteMetric.Success
}

func (dto *DeleteMetricOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.DeleteMetric.Errors))
	for i, err := range dto.Compass.DeleteMetric.Errors {
		errors[i] = err.Message
	}
	return errors
}
