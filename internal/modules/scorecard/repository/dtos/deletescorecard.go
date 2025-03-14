package dtos

type DeleteScorecard struct {
	Compass struct {
		DeleteScorecard struct {
			Success bool `json:"success"`
		} `json:"deleteScorecard"`
	} `json:"compass"`
}

func (d *DeleteScorecard) GetQuery() string {
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

func (d *DeleteScorecard) SetVariables(scorecardId string) map[string]interface{} {
	return map[string]interface{}{
		"scorecardId": scorecardId,
	}
}

func (c *DeleteScorecard) IsSuccessful() bool {
	return c.Compass.DeleteScorecard.Success
}
