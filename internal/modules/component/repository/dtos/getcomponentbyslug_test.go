package dtos_test

import (
	"reflect"
	"testing"

	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
)

func TestComponentByReferenceInput_GetQuery(t *testing.T) {
	input := &dtos.ComponentByReferenceInput{}
	expectedQuery := `
		query getComponentBySlug($cloudId: ID!, $slug: String!) {
			compass {
				componentByReference(reference: {slug: {slug: $slug, cloudId: $cloudId}}) {
					... on CompassComponent {
						id
						metricSources {
							... on CompassComponentMetricSourcesConnection {
								nodes {
									id,
									metricDefinition {
										name
									}
								}
							}
						}
					}
				}
			}
		}`

	if query := input.GetQuery(); query != expectedQuery {
		t.Errorf("GetQuery() = %v, want %v", query, expectedQuery)
	}
}

func TestComponentByReferenceInput_SetVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    *dtos.ComponentByReferenceInput
		expected map[string]interface{}
	}{
		{
			name: "Valid input",
			input: &dtos.ComponentByReferenceInput{
				CompassCloudID: "cloud123",
				Slug:           "test-slug",
			},
			expected: map[string]interface{}{
				"cloudId": "cloud123",
				"slug":    "test-slug",
			},
		},
		{
			name: "Empty input",
			input: &dtos.ComponentByReferenceInput{
				CompassCloudID: "",
				Slug:           "",
			},
			expected: map[string]interface{}{
				"cloudId": "",
				"slug":    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if variables := tt.input.SetVariables(); !reflect.DeepEqual(variables, tt.expected) {
				t.Errorf("SetVariables() = %v, want %v", variables, tt.expected)
			}
		})
	}
}

func TestComponentByReferenceOutput_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		output   *dtos.ComponentByReferenceOutput
		expected bool
	}{
		{
			name: "Successful output",
			output: &dtos.ComponentByReferenceOutput{
				Compass: struct {
					Component dtos.Component `json:"componentByReference"`
				}{
					Component: dtos.Component{ID: "component123"},
				},
			},
			expected: true,
		},
		{
			name: "Unsuccessful output",
			output: &dtos.ComponentByReferenceOutput{
				Compass: struct {
					Component dtos.Component `json:"componentByReference"`
				}{
					Component: dtos.Component{ID: ""},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if success := tt.output.IsSuccessful(); success != tt.expected {
				t.Errorf("IsSuccessful() = %v, want %v", success, tt.expected)
			}
		})
	}
}
func TestComponentByReferenceOutput_GetErrors(t *testing.T) {
	tests := []struct {
		name     string
		output   *dtos.ComponentByReferenceOutput
		expected []string
	}{
		{
			name: "No errors",
			output: &dtos.ComponentByReferenceOutput{
				Compass: struct {
					Component dtos.Component `json:"componentByReference"`
				}{
					Component: dtos.Component{ID: "component123"},
				},
			},
			expected: nil,
		},
		{
			name: "Empty component ID, no errors",
			output: &dtos.ComponentByReferenceOutput{
				Compass: struct {
					Component dtos.Component `json:"componentByReference"`
				}{
					Component: dtos.Component{ID: ""},
				},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if errors := tt.output.GetErrors(); !reflect.DeepEqual(errors, tt.expected) {
				t.Errorf("GetErrors() = %v, want %v", errors, tt.expected)
			}
		})
	}
}
