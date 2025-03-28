package dtos_test

import (
	"reflect"
	"testing"

	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
)

func TestSetVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    dtos.UpdateComponentInput
		expected map[string]interface{}
	}{
		{
			name: "With OwnerID",
			input: dtos.UpdateComponentInput{
				Component: resources.Component{
					ID:          "123",
					Name:        "Test Component",
					Slug:        "test-component",
					Description: "A test component",
					OwnerID:     "owner-456",
				},
			},
			expected: map[string]interface{}{
				"componentDetails": map[string]interface{}{
					"id":          "123",
					"name":        "Test Component",
					"slug":        "test-component",
					"description": "A test component",
					"ownerId":     "owner-456",
				},
			},
		},
		{
			name: "Without OwnerID",
			input: dtos.UpdateComponentInput{
				Component: resources.Component{
					ID:          "123",
					Name:        "Test Component",
					Slug:        "test-component",
					Description: "A test component",
				},
			},
			expected: map[string]interface{}{
				"componentDetails": map[string]interface{}{
					"id":          "123",
					"name":        "Test Component",
					"slug":        "test-component",
					"description": "A test component",
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

func TestGetQuery(t *testing.T) {
	input := dtos.UpdateComponentInput{}
	expectedQuery := `
		mutation updateComponent ($componentDetails: UpdateCompassComponentInput!) {
			compass {
				updateComponent(input: $componentDetails) {
					success
					errors {
						message
					}
				}
			}
		}`
	if query := input.GetQuery(); query != expectedQuery {
		t.Errorf("GetQuery() = %v, want %v", query, expectedQuery)
	}
}

func TestIsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		output   dtos.UpdateComponentOutput
		expected bool
	}{
		{
			name: "Success is true",
			output: dtos.UpdateComponentOutput{
				Compass: struct {
					UpdateComponent struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"updateComponent"`
				}{
					UpdateComponent: struct {
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
			output: dtos.UpdateComponentOutput{
				Compass: struct {
					UpdateComponent struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"updateComponent"`
				}{
					UpdateComponent: struct {
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
		output   dtos.UpdateComponentOutput
		expected []string
	}{
		{
			name: "With errors",
			output: dtos.UpdateComponentOutput{
				Compass: struct {
					UpdateComponent struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"updateComponent"`
				}{
					UpdateComponent: struct {
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
			output: dtos.UpdateComponentOutput{
				Compass: struct {
					UpdateComponent struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"updateComponent"`
				}{
					UpdateComponent: struct {
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
