package dtos

type UnbindMetricOutput struct {
	Compass struct {
		DeleteMetricSource struct {
			Success bool `json:"success"`
		} `json:"deleteMetricSource"`
	} `json:"compass"`
}

func (c *UnbindMetricOutput) GetQuery() string {
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

func (c *UnbindMetricOutput) SetVariables(id string) map[string]interface{} {
	return map[string]interface{}{
		"id": id,
	}
}

func (c *UnbindMetricOutput) IsSuccessful() bool {
	return c.Compass.DeleteMetricSource.Success
}
