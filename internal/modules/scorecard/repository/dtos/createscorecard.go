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
type CreateScorecardInput struct {
	compassdtos.InputDTO
	CompassCloudID string
	Scorecard      resources.Scorecard
}

func (dto *CreateScorecardInput) GetQuery() string {
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

func (dto *CreateScorecardInput) SetVariables() map[string]interface{} {
	criteria := make([]map[string]map[string]string, len(dto.Scorecard.Criteria))
	for i, criterion := range dto.Scorecard.Criteria {
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
		"cloudId": dto.CompassCloudID,
		"scorecardDetails": map[string]interface{}{
			"name":                dto.Scorecard.Name,
			"description":         dto.Scorecard.Description,
			"state":               dto.Scorecard.State,
			"componentTypeIds":    dto.Scorecard.ComponentTypeIDs,
			"importance":          dto.Scorecard.Importance,
			"scoringStrategyType": dto.Scorecard.ScoringStrategyType,
			"criterias":           criteria,
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
			Errors    []compassservice.CompassError `json:"errors"`
			Success   bool                          `json:"success"`
			Scorecard ScorecardDetails              `json:"scorecardDetails"`
		} `json:"createScorecard"`
	} `json:"compass"`
}

func (dto *CreateScorecardOutput) IsSuccessful() bool {
	return dto.Compass.CreateScorecard.Success
}

func (dto *CreateScorecardOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.CreateScorecard.Errors))
	for i, err := range dto.Compass.CreateScorecard.Errors {
		errors[i] = err.Message
	}
	return errors
}
