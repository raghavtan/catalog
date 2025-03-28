package dtos_test

import (
	"testing"

	"github.com/motain/of-catalog/internal/modules/scorecard/repository/dtos"
	"github.com/motain/of-catalog/internal/services/compassservice"
	"github.com/stretchr/testify/assert"
)

func TestDeleteScorecardInput_GetQuery(t *testing.T) {
	input := &dtos.DeleteScorecardInput{}
	expectedQuery := `
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

	assert.Equal(t, expectedQuery, input.GetQuery())
}

func TestDeleteScorecardInput_SetVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    *dtos.DeleteScorecardInput
		expected map[string]interface{}
	}{
		{
			name: "Valid ScorecardID",
			input: &dtos.DeleteScorecardInput{
				ScorecardID: "12345",
			},
			expected: map[string]interface{}{
				"scorecardId": "12345",
			},
		},
		{
			name: "Empty ScorecardID",
			input: &dtos.DeleteScorecardInput{
				ScorecardID: "",
			},
			expected: map[string]interface{}{
				"scorecardId": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.SetVariables())
		})
	}
}

func TestDeleteScorecardOutput_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		output   *dtos.DeleteScorecardOutput
		expected bool
	}{
		{
			name: "Success is true",
			output: &dtos.DeleteScorecardOutput{
				Compass: struct {
					DeleteScorecardOutput struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"DeleteScorecardOutput"`
				}{
					DeleteScorecardOutput: struct {
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
			output: &dtos.DeleteScorecardOutput{
				Compass: struct {
					DeleteScorecardOutput struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"DeleteScorecardOutput"`
				}{
					DeleteScorecardOutput: struct {
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

func TestDeleteScorecardOutput_GetErrors(t *testing.T) {
	tests := []struct {
		name     string
		output   *dtos.DeleteScorecardOutput
		expected []string
	}{
		{
			name: "Multiple errors",
			output: &dtos.DeleteScorecardOutput{
				Compass: struct {
					DeleteScorecardOutput struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"DeleteScorecardOutput"`
				}{
					DeleteScorecardOutput: struct {
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
			output: &dtos.DeleteScorecardOutput{
				Compass: struct {
					DeleteScorecardOutput struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"DeleteScorecardOutput"`
				}{
					DeleteScorecardOutput: struct {
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
