package dtos_test

import (
	"reflect"
	"testing"

	"github.com/motain/of-catalog/internal/modules/scorecard/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/scorecard/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
	"github.com/stretchr/testify/assert"
)

var (
	scorecard1 = "scorecard1"
)

func TestSetVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    dtos.UpdateScorecardInput
		expected map[string]interface{}
	}{
		{
			name: "Valid input with all fields",
			input: dtos.UpdateScorecardInput{
				Scorecard: resources.Scorecard{
					ID:                  &scorecard1,
					Name:                "Test Scorecard",
					Description:         "A test scorecard",
					State:               "ACTIVE",
					ComponentTypeIDs:    []string{"type1", "type2"},
					Importance:          "5",
					ScoringStrategyType: "strategy1",
					OwnerID:             "owner1",
				},
				CreateCriteria: []*resources.Criterion{
					{
						HasMetricValue: resources.MetricValue{
							Weight:             10,
							Name:               "Criterion1",
							MetricDefinitionId: "metric1",
							ComparatorValue:    100,
							Comparator:         "GREATER_THAN",
						},
					},
				},
				UpdateCriteria: []*resources.Criterion{
					{
						HasMetricValue: resources.MetricValue{
							ID:                 "criterion1",
							Weight:             20,
							Name:               "Updated Criterion",
							MetricDefinitionId: "metric2",
							ComparatorValue:    200,
							Comparator:         "LESS_THAN",
						},
					},
				},
				DeleteCriteria: []string{"criterion2"},
			},
			expected: map[string]interface{}{
				"scorecardId": &scorecard1,
				"scorecardDetails": map[string]interface{}{
					"name":                "Test Scorecard",
					"description":         "A test scorecard",
					"state":               "ACTIVE",
					"componentTypeIds":    []string{"type1", "type2"},
					"importance":          "5",
					"scoringStrategyType": "strategy1",
					"ownerId":             "owner1",
					"createCriteria": []map[string]map[string]string{
						{
							"hasMetricValue": {
								"weight":             "10",
								"name":               "Criterion1",
								"metricDefinitionId": "metric1",
								"comparatorValue":    "100",
								"comparator":         "GREATER_THAN",
							},
						},
					},
					"updateCriteria": []map[string]map[string]string{
						{
							"hasMetricValue": {
								"id":                 "criterion1",
								"weight":             "20",
								"name":               "Updated Criterion",
								"metricDefinitionId": "metric2",
								"comparatorValue":    "200",
								"comparator":         "LESS_THAN",
							},
						},
					},
					"deleteCriteria": []string{"criterion2"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.SetVariables()
			assert.Equal(t, tt.expected, result, "SetVariables() result mismatch")
		})
	}
}

func TestGetQuery(t *testing.T) {
	dto := dtos.UpdateScorecardInput{}
	expectedQuery := `
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
	if query := dto.GetQuery(); query != expectedQuery {
		t.Errorf("GetQuery() = %v, want %v", query, expectedQuery)
	}
}

func TestIsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		output   dtos.UpdateScorecardOutput
		expected bool
	}{
		{
			name: "Success is true",
			output: dtos.UpdateScorecardOutput{
				Compass: struct {
					UpdateScorecardOutput struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"updateScorecard"`
				}{
					UpdateScorecardOutput: struct {
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
			output: dtos.UpdateScorecardOutput{
				Compass: struct {
					UpdateScorecardOutput struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"updateScorecard"`
				}{
					UpdateScorecardOutput: struct {
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
			if result := tt.output.IsSuccessful(); result != tt.expected {
				t.Errorf("IsSuccessful() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetErrors(t *testing.T) {
	tests := []struct {
		name     string
		output   dtos.UpdateScorecardOutput
		expected []string
	}{
		{
			name: "With errors",
			output: dtos.UpdateScorecardOutput{
				Compass: struct {
					UpdateScorecardOutput struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"updateScorecard"`
				}{
					UpdateScorecardOutput: struct {
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
			output: dtos.UpdateScorecardOutput{
				Compass: struct {
					UpdateScorecardOutput struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"updateScorecard"`
				}{
					UpdateScorecardOutput: struct {
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
			if result := tt.output.GetErrors(); !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetErrors() = %v, want %v", result, tt.expected)
			}
		})
	}
}
