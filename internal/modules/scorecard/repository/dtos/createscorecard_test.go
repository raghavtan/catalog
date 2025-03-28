package dtos_test

import (
	"reflect"
	"testing"

	"github.com/motain/of-catalog/internal/modules/scorecard/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/scorecard/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
)

func TestCreateScorecardInputSetVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    dtos.CreateScorecardInput
		expected map[string]interface{}
	}{
		{
			name: "valid input with owner ID",
			input: dtos.CreateScorecardInput{
				CompassCloudID: "cloud-123",
				Scorecard: resources.Scorecard{
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
								Name:               "Metric 1",
								MetricDefinitionId: "metric-1",
								ComparatorValue:    5,
								Comparator:         "GREATER_THAN",
							},
						},
					},
				},
			},
			expected: map[string]interface{}{
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
								"name":               "Metric 1",
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
			name: "valid input without owner ID",
			input: dtos.CreateScorecardInput{
				CompassCloudID: "cloud-456",
				Scorecard: resources.Scorecard{
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
								Name:               "Metric 2",
								MetricDefinitionId: "metric-2",
								ComparatorValue:    10,
								Comparator:         "LESS_THAN",
							},
						},
					},
				},
			},
			expected: map[string]interface{}{
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
								"name":               "Metric 2",
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
			result := tt.input.SetVariables()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SetVariables() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCreateScorecardInputGetQuery(t *testing.T) {
	dto := dtos.CreateScorecardInput{}
	expectedQuery := `
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
	if query := dto.GetQuery(); query != expectedQuery {
		t.Errorf("GetQuery() = %v, want %v", query, expectedQuery)
	}
}

func TestCreateScorecardOuputIsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		output   dtos.CreateScorecardOutput
		expected bool
	}{
		{
			name: "success is true",
			output: dtos.CreateScorecardOutput{
				Compass: struct {
					CreateScorecard struct {
						Errors    []compassservice.CompassError `json:"errors"`
						Success   bool                          `json:"success"`
						Scorecard dtos.ScorecardDetails         `json:"scorecardDetails"`
					} `json:"createScorecard"`
				}{
					CreateScorecard: struct {
						Errors    []compassservice.CompassError `json:"errors"`
						Success   bool                          `json:"success"`
						Scorecard dtos.ScorecardDetails         `json:"scorecardDetails"`
					}{
						Success: true,
					},
				},
			},
			expected: true,
		},
		{
			name: "success is false",
			output: dtos.CreateScorecardOutput{
				Compass: struct {
					CreateScorecard struct {
						Errors    []compassservice.CompassError `json:"errors"`
						Success   bool                          `json:"success"`
						Scorecard dtos.ScorecardDetails         `json:"scorecardDetails"`
					} `json:"createScorecard"`
				}{
					CreateScorecard: struct {
						Errors    []compassservice.CompassError `json:"errors"`
						Success   bool                          `json:"success"`
						Scorecard dtos.ScorecardDetails         `json:"scorecardDetails"`
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

func TestCreateScorecardOutputGetErrors(t *testing.T) {
	tests := []struct {
		name     string
		output   dtos.CreateScorecardOutput
		expected []string
	}{
		{
			name: "with errors",
			output: dtos.CreateScorecardOutput{
				Compass: struct {
					CreateScorecard struct {
						Errors    []compassservice.CompassError `json:"errors"`
						Success   bool                          `json:"success"`
						Scorecard dtos.ScorecardDetails         `json:"scorecardDetails"`
					} `json:"createScorecard"`
				}{
					CreateScorecard: struct {
						Errors    []compassservice.CompassError `json:"errors"`
						Success   bool                          `json:"success"`
						Scorecard dtos.ScorecardDetails         `json:"scorecardDetails"`
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
			name: "no errors",
			output: dtos.CreateScorecardOutput{
				Compass: struct {
					CreateScorecard struct {
						Errors    []compassservice.CompassError `json:"errors"`
						Success   bool                          `json:"success"`
						Scorecard dtos.ScorecardDetails         `json:"scorecardDetails"`
					} `json:"createScorecard"`
				}{
					CreateScorecard: struct {
						Errors    []compassservice.CompassError `json:"errors"`
						Success   bool                          `json:"success"`
						Scorecard dtos.ScorecardDetails         `json:"scorecardDetails"`
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
			result := tt.output.GetErrors()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetErrors() = %v, want %v", result, tt.expected)
			}
		})
	}
}
