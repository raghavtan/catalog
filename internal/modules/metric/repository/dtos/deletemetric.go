package dtos

type DeleteMetricOutput struct {
	Compass struct {
		DeleteMetric struct {
			Success bool `json:"success"`
		} `json:"deleteMetricDefinition"`
	} `json:"compass"`
}

func (d *DeleteMetricOutput) GetQuery() string {
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

func (d *DeleteMetricOutput) SetVariables(id string) map[string]interface{} {
	return map[string]interface{}{
		"id": id,
	}
}

func (c *DeleteMetricOutput) IsSuccessful() bool {
	return c.Compass.DeleteMetric.Success
}
