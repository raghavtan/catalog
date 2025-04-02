package dtos

import (
	"fmt"

	"github.com/motain/of-catalog/internal/modules/scorecard/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

/*************
 * INPUT DTO *
 *************/

type UpdateScorecardInput struct {
	compassdtos.InputDTO
	Scorecard      resources.Scorecard
	CreateCriteria []*resources.Criterion
	UpdateCriteria []*resources.Criterion
	DeleteCriteria []string
}

func (dto *UpdateScorecardInput) GetQuery() string {
	return `
		mutation updateScorecard ($scorecardId: ID! $scorecardDetails: UpdateCompassScorecardInput!) {
			compass {
				updateScorecard(scorecardId: $scorecardId, input: $scorecardDetails) {
					success
					errors {
						message
					}
				}
			}
		}`
}

func (dto *UpdateScorecardInput) SetVariables() map[string]interface{} {
	criteriaToAdd := make([]map[string]map[string]string, len(dto.CreateCriteria))
	for i, criterion := range dto.CreateCriteria {
		criteriaToAdd[i] = make(map[string]map[string]string)
		criteriaToAdd[i]["hasMetricValue"] = make(map[string]string)
		criteriaToAdd[i]["hasMetricValue"] = map[string]string{
			"weight":             fmt.Sprintf("%d", criterion.HasMetricValue.Weight),
			"name":               criterion.HasMetricValue.Name,
			"metricDefinitionId": criterion.HasMetricValue.MetricDefinitionId,
			"comparatorValue":    fmt.Sprintf("%d", criterion.HasMetricValue.ComparatorValue),
			"comparator":         criterion.HasMetricValue.Comparator,
		}
	}

	criteriaToUpdate := make([]map[string]map[string]string, len(dto.UpdateCriteria))
	for i, criterion := range dto.UpdateCriteria {
		criteriaToUpdate[i] = make(map[string]map[string]string)
		criteriaToUpdate[i]["hasMetricValue"] = make(map[string]string)
		criteriaToUpdate[i]["hasMetricValue"] = map[string]string{
			"id":                 criterion.HasMetricValue.ID,
			"weight":             fmt.Sprintf("%d", criterion.HasMetricValue.Weight),
			"name":               criterion.HasMetricValue.Name,
			"metricDefinitionId": criterion.HasMetricValue.MetricDefinitionId,
			"comparatorValue":    fmt.Sprintf("%d", criterion.HasMetricValue.ComparatorValue),
			"comparator":         criterion.HasMetricValue.Comparator,
		}
	}

	variables := map[string]interface{}{
		"scorecardId": dto.Scorecard.ID,
		"scorecardDetails": map[string]interface{}{
			"name":                dto.Scorecard.Name,
			"description":         dto.Scorecard.Description,
			"state":               dto.Scorecard.State,
			"componentTypeIds":    dto.Scorecard.ComponentTypeIDs,
			"importance":          dto.Scorecard.Importance,
			"scoringStrategyType": dto.Scorecard.ScoringStrategyType,
			"createCriteria":      criteriaToAdd,
			"updateCriteria":      criteriaToUpdate,
			"deleteCriteria":      dto.DeleteCriteria,
		},
	}

	if dto.Scorecard.OwnerID != "" {
		variables["scorecardDetails"].(map[string]interface{})["ownerId"] = dto.Scorecard.OwnerID
	}

	return variables
}

/**************
 * OUTPUT DTO *
 **************/
type UpdateScorecardOutput struct {
	Compass struct {
		UpdateScorecardOutput struct {
			Errors  []compassservice.CompassError `json:"errors"`
			Success bool                          `json:"success"`
		} `json:"updateScorecard"`
	} `json:"compass"`
}

func (dto *UpdateScorecardOutput) IsSuccessful() bool {
	return dto.Compass.UpdateScorecardOutput.Success
}

func (dto *UpdateScorecardOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.UpdateScorecardOutput.Errors))
	for i, err := range dto.Compass.UpdateScorecardOutput.Errors {
		errors[i] = err.Message
	}
	return errors
}
