package dtos_test

import (
	"testing"

	"github.com/motain/of-catalog/internal/modules/metric/repository/dtos"
	"github.com/motain/of-catalog/internal/services/compassservice"
	"github.com/stretchr/testify/assert"
)

func TestDeleteMetricInput_GetQuery(t *testing.T) {
	input := &dtos.DeleteMetricInput{}
	expectedQuery := `
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

	assert.Equal(t, expectedQuery, input.GetQuery())
}

func TestDeleteMetricInput_SetVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    *dtos.DeleteMetricInput
		expected map[string]interface{}
	}{
		{
			name: "Valid MetricID",
			input: &dtos.DeleteMetricInput{
				MetricID: "12345",
			},
			expected: map[string]interface{}{
				"id": "12345",
			},
		},
		{
			name: "Empty MetricID",
			input: &dtos.DeleteMetricInput{
				MetricID: "",
			},
			expected: map[string]interface{}{
				"id": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.SetVariables())
		})
	}
}

func TestDeleteMetricOutput_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		output   *dtos.DeleteMetricOutput
		expected bool
	}{
		{
			name: "Success is true",
			output: &dtos.DeleteMetricOutput{
				Compass: struct {
					DeleteMetric struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"deleteMetricDefinition"`
				}{
					DeleteMetric: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					}{
						Success: true,
					},
				},
			},
			expected: true,
		},
		{
			name: "Success is false",
			output: &dtos.DeleteMetricOutput{
				Compass: struct {
					DeleteMetric struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"deleteMetricDefinition"`
				}{
					DeleteMetric: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					}{
						Success: false,
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.output.IsSuccessful())
		})
	}
}

func TestDeleteMetricOutput_GetErrors(t *testing.T) {
	tests := []struct {
		name     string
		output   *dtos.DeleteMetricOutput
		expected []string
	}{
		{
			name: "Multiple errors",
			output: &dtos.DeleteMetricOutput{
				Compass: struct {
					DeleteMetric struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"deleteMetricDefinition"`
				}{
					DeleteMetric: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					}{
						Errors: []compassservice.CompassError{
							{Message: "Error 1"},
							{Message: "Error 2"},
						},
					},
				},
			},
			expected: []string{"Error 1", "Error 2"},
		},
		{
			name: "No errors",
			output: &dtos.DeleteMetricOutput{
				Compass: struct {
					DeleteMetric struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"deleteMetricDefinition"`
				}{
					DeleteMetric: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					}{
						Errors: []compassservice.CompassError{},
					},
				},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.output.GetErrors())
		})
	}
}
