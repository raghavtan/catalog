package dtos_test

import (
	"reflect"
	"testing"

	"github.com/motain/of-catalog/internal/modules/metric/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/metric/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
)

func TestCreateMetricInput_GetQuery(t *testing.T) {
	dto := &dtos.CreateMetricInput{}
	expectedQuery := `
		mutation createMetricDefinition ($cloudId: ID!, $name: String!, $description: String!, $unit: String!) {
			compass {
				createMetricDefinition(
					input: {
						cloudId: $cloudId
						name: $name
						description: $description
						format: {
							suffix: { suffix: $unit }
						}
					}
				) {
					success
					createdMetricDefinition {
						id
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

func TestCreateMetricInput_SetVariables(t *testing.T) {
	dto := &dtos.CreateMetricInput{
		CompassCloudID: "cloud-123",
		Metric: resources.Metric{
			Name:        "CPU Usage",
			Description: "Tracks CPU usage",
			Format: resources.MetricFormat{
				Unit: "percentage",
			},
		},
	}

	expectedVariables := map[string]interface{}{
		"cloudId":     "cloud-123",
		"name":        "CPU Usage",
		"description": "Tracks CPU usage",
		"unit":        "percentage",
	}

	if variables := dto.SetVariables(); !reflect.DeepEqual(variables, expectedVariables) {
		t.Errorf("SetVariables() = %v, want %v", variables, expectedVariables)
	}
}

func TestCreateMetricOutput_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		dto      dtos.CreateMetricOutput
		expected bool
	}{
		{
			name: "Success case",
			dto: dtos.CreateMetricOutput{
				Compass: struct {
					CreateMetric struct {
						Success    bool                          `json:"success"`
						Errors     []compassservice.CompassError `json:"errors"`
						Definition dtos.Metric                   `json:"createdMetricDefinition"`
					} `json:"createMetricDefinition"`
				}{
					CreateMetric: struct {
						Success    bool                          `json:"success"`
						Errors     []compassservice.CompassError `json:"errors"`
						Definition dtos.Metric                   `json:"createdMetricDefinition"`
					}{
						Success: true,
					},
				},
			},
			expected: true,
		},
		{
			name: "Failure case",
			dto: dtos.CreateMetricOutput{
				Compass: struct {
					CreateMetric struct {
						Success    bool                          `json:"success"`
						Errors     []compassservice.CompassError `json:"errors"`
						Definition dtos.Metric                   `json:"createdMetricDefinition"`
					} `json:"createMetricDefinition"`
				}{
					CreateMetric: struct {
						Success    bool                          `json:"success"`
						Errors     []compassservice.CompassError `json:"errors"`
						Definition dtos.Metric                   `json:"createdMetricDefinition"`
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

func TestCreateMetricOutput_GetErrors(t *testing.T) {
	dto := &dtos.CreateMetricOutput{
		Compass: struct {
			CreateMetric struct {
				Success    bool                          `json:"success"`
				Errors     []compassservice.CompassError `json:"errors"`
				Definition dtos.Metric                   `json:"createdMetricDefinition"`
			} `json:"createMetricDefinition"`
		}{
			CreateMetric: struct {
				Success    bool                          `json:"success"`
				Errors     []compassservice.CompassError `json:"errors"`
				Definition dtos.Metric                   `json:"createdMetricDefinition"`
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
