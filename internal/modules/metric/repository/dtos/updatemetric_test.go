package dtos_test

import (
	"reflect"
	"testing"

	"github.com/motain/of-catalog/internal/modules/metric/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/metric/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
)

func TestUpdateMetricInput_GetQuery(t *testing.T) {
	dto := &dtos.UpdateMetricInput{}
	expectedQuery := `
		mutation updateMetricDefinition ($cloudId: ID!, $id: ID!, $name: String!, $description: String!, $unit: String!) {
			compass {
				updateMetricDefinition(
					input: {
						id: $id
						cloudId: $cloudId
						name: $name
						description: $description
						format: {
							suffix: { suffix: $unit }
						}
					}
				) {
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

func TestUpdateMetricInput_SetVariables(t *testing.T) {
	dto := &dtos.UpdateMetricInput{
		CompassCloudID: "cloud123",
		Metric: resources.Metric{
			ID:          "metric123",
			Name:        "Test Metric",
			Description: "A test metric",
			Format: resources.MetricFormat{
				Unit: "ms",
			},
		},
	}

	expectedVariables := map[string]interface{}{
		"cloudId":     "cloud123",
		"id":          "metric123",
		"name":        "Test Metric",
		"description": "A test metric",
		"unit":        "ms",
	}

	if variables := dto.SetVariables(); !reflect.DeepEqual(variables, expectedVariables) {
		t.Errorf("SetVariables() = %v, want %v", variables, expectedVariables)
	}
}

func TestUpdateMetricOutput_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		dto      dtos.UpdateMetricOutput
		expected bool
	}{
		{
			name: "Success case",
			dto: dtos.UpdateMetricOutput{
				Compass: struct {
					UpdateMetric struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"updateMetricDefinition"`
				}{
					UpdateMetric: struct {
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
			name: "Failure case",
			dto: dtos.UpdateMetricOutput{
				Compass: struct {
					UpdateMetric struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"updateMetricDefinition"`
				}{
					UpdateMetric: struct {
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
			if result := tt.dto.IsSuccessful(); result != tt.expected {
				t.Errorf("IsSuccessful() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUpdateMetricOutput_GetErrors(t *testing.T) {
	dto := &dtos.UpdateMetricOutput{
		Compass: struct {
			UpdateMetric struct {
				Errors  []compassservice.CompassError `json:"errors"`
				Success bool                          `json:"success"`
			} `json:"updateMetricDefinition"`
		}{
			UpdateMetric: struct {
				Errors  []compassservice.CompassError `json:"errors"`
				Success bool                          `json:"success"`
			}{
				Errors: []compassservice.CompassError{
					{Message: "Error 1"},
					{Message: "Error 2"},
				},
			},
		},
	}

	expectedErrors := []string{"Error 1", "Error 2"}

	if errors := dto.GetErrors(); !reflect.DeepEqual(errors, expectedErrors) {
		t.Errorf("GetErrors() = %v, want %v", errors, expectedErrors)
	}
}
