package dtos

import (
	"fmt"

	"github.com/motain/of-catalog/internal/modules/scorecard/resources"
)

type ScorecardDetails struct {
	ID       string      `json:"id"`
	Criteria []Criterion `json:"criterias"`
}

type Criterion struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateScorecardOutput struct {
	Compass struct {
		CreateScorecard struct {
			Success   bool             `json:"success"`
			Scorecard ScorecardDetails `json:"scorecardDetails"`
		} `json:"createScorecard"`
	} `json:"compass"`
}

func (c *CreateScorecardOutput) GetQuery() string {
	return `
		mutation createScorecard ($cloudId: ID!, $scorecardDetails: CreateCompassScorecardInput!) {
			compass {
				createScorecard(cloudId: $cloudId, input: $scorecardDetails) {
					success
					scorecardDetails {
						id
						criterias {
							id
							name
						}
					}
					errors {
						message
					}
				}
			}
		}`
}

func (c *CreateScorecardOutput) SetVariables(compassCloudIdD string, scorecard resources.Scorecard) map[string]interface{} {
	criteria := make([]map[string]map[string]string, len(scorecard.Criteria))
	for i, criterion := range scorecard.Criteria {
		criteria[i] = make(map[string]map[string]string)
		criteria[i]["hasMetricValue"] = make(map[string]string)
		criteria[i]["hasMetricValue"] = map[string]string{
			"weight":             fmt.Sprintf("%d", criterion.HasMetricValue.Weight),
			"name":               criterion.HasMetricValue.Name,
			"metricDefinitionId": criterion.HasMetricValue.MetricDefinitionId,
			"comparatorValue":    fmt.Sprintf("%d", criterion.HasMetricValue.ComparatorValue),
			"comparator":         criterion.HasMetricValue.Comparator,
		}
	}

	variables := map[string]interface{}{
		"cloudId": compassCloudIdD,
		"scorecardDetails": map[string]interface{}{
			"name":                scorecard.Name,
			"description":         scorecard.Description,
			"state":               scorecard.State,
			"componentTypeIds":    scorecard.ComponentTypeIDs,
			"importance":          scorecard.Importance,
			"scoringStrategyType": scorecard.ScoringStrategyType,
			"criterias":           criteria,
		},
	}

	if scorecard.OwnerID != "" {
		variables["scorecardDetails"].(map[string]interface{})["ownerId"] = scorecard.OwnerID
	}

	return variables
}

func (c *CreateScorecardOutput) IsSuccessful() bool {
	return c.Compass.CreateScorecard.Success
}
