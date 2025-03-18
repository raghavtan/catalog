package dtos_test

import (
	"reflect"
	"testing"

	"github.com/motain/of-catalog/internal/modules/scorecard/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/scorecard/resources"
)

func TestSetVariables(t *testing.T) {
	tests := []struct {
		name           string
		cloudId        string
		scorecard      resources.Scorecard
		expectedResult map[string]interface{}
	}{
		{
			name:    "Valid input with owner ID",
			cloudId: "cloud-123",
			scorecard: resources.Scorecard{
				Name:                "Test Scorecard",
				Description:         "A test scorecard",
				State:               "ACTIVE",
				ComponentTypeIDs:    []string{"type-1", "type-2"},
				Importance:          "HIGH",
				ScoringStrategyType: "WEIGHTED",
				OwnerID:             "owner-123",
				Criteria: []*resources.Criterion{
					{
						HasMetricValue: resources.MetricValue{
							Weight:             10,
							Name:               "Criterion 1",
							MetricDefinitionId: "metric-1",
							ComparatorValue:    5,
							Comparator:         "GREATER_THAN",
						},
					},
				},
			},
			expectedResult: map[string]interface{}{
				"cloudId": "cloud-123",
				"scorecardDetails": map[string]interface{}{
					"name":                "Test Scorecard",
					"description":         "A test scorecard",
					"state":               "ACTIVE",
					"componentTypeIds":    []string{"type-1", "type-2"},
					"importance":          "HIGH",
					"scoringStrategyType": "WEIGHTED",
					"ownerId":             "owner-123",
					"criterias": []map[string]map[string]string{
						{
							"hasMetricValue": {
								"weight":             "10",
								"name":               "Criterion 1",
								"metricDefinitionId": "metric-1",
								"comparatorValue":    "5",
								"comparator":         "GREATER_THAN",
							},
						},
					},
				},
			},
		},
		{
			name:    "Valid input without owner ID",
			cloudId: "cloud-456",
			scorecard: resources.Scorecard{
				Name:                "Another Scorecard",
				Description:         "Another test scorecard",
				State:               "INACTIVE",
				ComponentTypeIDs:    []string{"type-3"},
				Importance:          "LOW",
				ScoringStrategyType: "SIMPLE",
				Criteria: []*resources.Criterion{
					{
						HasMetricValue: resources.MetricValue{
							Weight:             20,
							Name:               "Criterion 2",
							MetricDefinitionId: "metric-2",
							ComparatorValue:    10,
							Comparator:         "LESS_THAN",
						},
					},
				},
			},
			expectedResult: map[string]interface{}{
				"cloudId": "cloud-456",
				"scorecardDetails": map[string]interface{}{
					"name":                "Another Scorecard",
					"description":         "Another test scorecard",
					"state":               "INACTIVE",
					"componentTypeIds":    []string{"type-3"},
					"importance":          "LOW",
					"scoringStrategyType": "SIMPLE",
					"criterias": []map[string]map[string]string{
						{
							"hasMetricValue": {
								"weight":             "20",
								"name":               "Criterion 2",
								"metricDefinitionId": "metric-2",
								"comparatorValue":    "10",
								"comparator":         "LESS_THAN",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := dtos.CreateScorecardOutput{}
			result := output.SetVariables(tt.cloudId, tt.scorecard)

			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("SetVariables() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}
func TestGetQuery(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name: "Valid query generation",
			expected: `
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
		}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := dtos.CreateScorecardOutput{}
			query := output.GetQuery()

			if query != tt.expected {
				t.Errorf("GetQuery() = %v, want %v", query, tt.expected)
			}
		})
	}
}
func TestIsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		output   dtos.CreateScorecardOutput
		expected bool
	}{
		{
			name: "Success is true",
			output: dtos.CreateScorecardOutput{
				Compass: struct {
					CreateScorecard struct {
						Success   bool                  `json:"success"`
						Scorecard dtos.ScorecardDetails `json:"scorecardDetails"`
					} `json:"createScorecard"`
				}{
					CreateScorecard: struct {
						Success   bool                  `json:"success"`
						Scorecard dtos.ScorecardDetails `json:"scorecardDetails"`
					}{
						Success: true,
					},
				},
			},
			expected: true,
		},
		{
			name: "Success is false",
			output: dtos.CreateScorecardOutput{
				Compass: struct {
					CreateScorecard struct {
						Success   bool                  `json:"success"`
						Scorecard dtos.ScorecardDetails `json:"scorecardDetails"`
					} `json:"createScorecard"`
				}{
					CreateScorecard: struct {
						Success   bool                  `json:"success"`
						Scorecard dtos.ScorecardDetails `json:"scorecardDetails"`
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
			result := tt.output.IsSuccessful()
			if result != tt.expected {
				t.Errorf("IsSuccessful() = %v, want %v", result, tt.expected)
			}
		})
	}
}
