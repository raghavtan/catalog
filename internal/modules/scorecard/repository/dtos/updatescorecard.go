package dtos

import (
	"fmt"

	"github.com/motain/of-catalog/internal/modules/scorecard/resources"
)

type UpdateScorecard struct {
	Compass struct {
		UpdateScorecard struct {
			Success bool `json:"success"`
		} `json:"updateScorecard"`
	} `json:"compass"`
}

func (u *UpdateScorecard) GetQuery() string {
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

func (u *UpdateScorecard) SetVariables(
	scorecard resources.Scorecard,
	createCriteria []*resources.Criterion,
	updateCriteria []*resources.Criterion,
	deleteCriteria []string,
) map[string]interface{} {
	criteriaToAdd := make([]map[string]map[string]string, len(createCriteria))
	for i, criterion := range createCriteria {
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

	criteriaToUpdate := make([]map[string]map[string]string, len(updateCriteria))
	for i, criterion := range updateCriteria {
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
		"scorecardId": scorecard.ID,
		"scorecardDetails": map[string]interface{}{
			"name":                scorecard.Name,
			"description":         scorecard.Description,
			"state":               scorecard.State,
			"componentTypeIds":    scorecard.ComponentTypeIDs,
			"importance":          scorecard.Importance,
			"scoringStrategyType": scorecard.ScoringStrategyType,
			"createCriteria":      criteriaToAdd,
			"updateCriteria":      criteriaToUpdate,
			"deleteCriteria":      deleteCriteria,
		},
	}

	if scorecard.OwnerID != "" {
		variables["scorecardDetails"].(map[string]interface{})["ownerId"] = scorecard.OwnerID
	}

	return variables
}

func (c *UpdateScorecard) IsSuccessful() bool {
	return c.Compass.UpdateScorecard.Success
}
