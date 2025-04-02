package dtos

import (
	"github.com/motain/of-catalog/internal/services/compassservice"
	compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"
)

/*************
 * INPUT DTO *
 *************/

type DeleteScorecardInput struct {
	compassdtos.InputDTO
	ScorecardID string
}

func (dto *DeleteScorecardInput) GetQuery() string {
	return `
		mutation deleteScorecard($scorecardId: ID!) {
			compass {
				deleteScorecard(scorecardId: $scorecardId) {
					scorecardId
					errors {
						message
					}
					success
				}
			}
		}`
}

func (dto *DeleteScorecardInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"scorecardId": dto.ScorecardID,
	}
}

/**************
 * OUTPUT DTO *
 **************/

type DeleteScorecardOutput struct {
	Compass struct {
		DeleteScorecardOutput struct {
			Errors  []compassservice.CompassError `json:"errors"`
			Success bool                          `json:"success"`
		} `json:"DeleteScorecardOutput"`
	} `json:"compass"`
}

func (dto *DeleteScorecardOutput) IsSuccessful() bool {
	return dto.Compass.DeleteScorecardOutput.Success
}

func (dto *DeleteScorecardOutput) GetErrors() []string {
	errors := make([]string, len(dto.Compass.DeleteScorecardOutput.Errors))
	for i, err := range dto.Compass.DeleteScorecardOutput.Errors {
		errors[i] = err.Message
	}
	return errors
}
