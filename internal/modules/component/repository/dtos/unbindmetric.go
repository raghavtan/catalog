package dtos

import (
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

/*************
 * INPUT DTO *
 *************/
type UnbindMetricInput struct {
	compassdtos.InputDTO
	MetricID string
}

func (dto *UnbindMetricInput) GetQuery() string {
	return `
		mutation deleteMetricSource($id: ID!) {
			compass {
				deleteMetricSource(input: {id: $id}) {
					deletedMetricSourceId
					errors {
						message
					}
					success
				}
			}
		}`
}

func (dto *UnbindMetricInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"id": dto.MetricID,
	}
}

/**************
 * OUTPUT DTO *
 **************/

type UnbindMetricOutput struct {
	Compass struct {
		DeleteMetricSource struct {
			Errors  []compassservice.CompassError `json:"errors"`
			Success bool                          `json:"success"`
		} `json:"deleteMetricSource"`
	} `json:"compass"`
}

func (dto *UnbindMetricOutput) IsSuccessful() bool {
	return dto.Compass.DeleteMetricSource.Success
}

func (dto *UnbindMetricOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.DeleteMetricSource.Errors))
	for i, err := range dto.Compass.DeleteMetricSource.Errors {
		errors[i] = err.Message
	}
	return errors
}
