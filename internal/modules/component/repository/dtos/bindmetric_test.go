package dtos_test

import (
	"testing"

	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/services/compassservice"
	"github.com/stretchr/testify/assert"
)

func TestBindMetricInput_GetQuery(t *testing.T) {
	tests := []struct {
		name     string
		dto      dtos.BindMetricInput
		expected string
	}{
		{
			name: "returns correct GraphQL query",
			dto:  dtos.BindMetricInput{},
			expected: `
		mutation createMetricSource($metricId: ID!, $componentId: ID!, $externalId: ID!) {
			compass {
				createMetricSource(input: {metricDefinitionId: $metricId, componentId: $componentId, externalMetricSourceId: $externalId}) {
					success
					createdMetricSource {
						id
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
			query := tt.dto.GetQuery()
			assert.Equal(t, tt.expected, query)
		})
	}
}
func TestBindMetricInput_SetVariables(t *testing.T) {
	tests := []struct {
		name     string
		dto      dtos.BindMetricInput
		expected map[string]interface{}
	}{
		{
			name: "sets variables correctly",
			dto: dtos.BindMetricInput{
				MetricID:    "metric-123",
				ComponentID: "component-456",
				Identifier:  "external-789",
			},
			expected: map[string]interface{}{
				"metricId":    "metric-123",
				"componentId": "component-456",
				"externalId":  "external-789",
			},
		},
		{
			name: "handles empty fields",
			dto: dtos.BindMetricInput{
				MetricID:    "",
				ComponentID: "",
				Identifier:  "",
			},
			expected: map[string]interface{}{
				"metricId":    "",
				"componentId": "",
				"externalId":  "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variables := tt.dto.SetVariables()
			assert.Equal(t, tt.expected, variables)
		})
	}
}
func TestBindMetricOutput_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		dto      dtos.BindMetricOutput
		expected bool
	}{
		{
			name: "returns true when success is true",
			dto: dtos.BindMetricOutput{
				Compass: struct {
					CreateMetricSource struct {
						Errors             []compassservice.CompassError `json:"errors"`
						Success            bool                          `json:"success"`
						CreateMetricSource struct {
							ID string `json:"id"`
						} `json:"createdMetricSource"`
					} `json:"createMetricSource"`
				}{
					CreateMetricSource: struct {
						Errors             []compassservice.CompassError `json:"errors"`
						Success            bool                          `json:"success"`
						CreateMetricSource struct {
							ID string `json:"id"`
						} `json:"createdMetricSource"`
					}{
						Success: true,
					},
				},
			},
			expected: true,
		},
		{
			name: "returns false when success is false",
			dto: dtos.BindMetricOutput{
				Compass: struct {
					CreateMetricSource struct {
						Errors             []compassservice.CompassError `json:"errors"`
						Success            bool                          `json:"success"`
						CreateMetricSource struct {
							ID string `json:"id"`
						} `json:"createdMetricSource"`
					} `json:"createMetricSource"`
				}{
					CreateMetricSource: struct {
						Errors             []compassservice.CompassError `json:"errors"`
						Success            bool                          `json:"success"`
						CreateMetricSource struct {
							ID string `json:"id"`
						} `json:"createdMetricSource"`
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
			result := tt.dto.IsSuccessful()
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestBindMetricOutput_GetErrors(t *testing.T) {
	tests := []struct {
		name     string
		dto      dtos.BindMetricOutput
		expected []string
	}{
		{
			name: "returns error messages when errors are present",
			dto: dtos.BindMetricOutput{
				Compass: struct {
					CreateMetricSource struct {
						Errors             []compassservice.CompassError `json:"errors"`
						Success            bool                          `json:"success"`
						CreateMetricSource struct {
							ID string `json:"id"`
						} `json:"createdMetricSource"`
					} `json:"createMetricSource"`
				}{
					CreateMetricSource: struct {
						Errors             []compassservice.CompassError `json:"errors"`
						Success            bool                          `json:"success"`
						CreateMetricSource struct {
							ID string `json:"id"`
						} `json:"createdMetricSource"`
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
			name: "returns empty slice when no errors are present",
			dto: dtos.BindMetricOutput{
				Compass: struct {
					CreateMetricSource struct {
						Errors             []compassservice.CompassError `json:"errors"`
						Success            bool                          `json:"success"`
						CreateMetricSource struct {
							ID string `json:"id"`
						} `json:"createdMetricSource"`
					} `json:"createMetricSource"`
				}{
					CreateMetricSource: struct {
						Errors             []compassservice.CompassError `json:"errors"`
						Success            bool                          `json:"success"`
						CreateMetricSource struct {
							ID string `json:"id"`
						} `json:"createdMetricSource"`
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
			errors := tt.dto.GetErrors()
			assert.Equal(t, tt.expected, errors)
		})
	}
}
