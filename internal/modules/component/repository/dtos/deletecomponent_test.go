package dtos_test

import (
	"testing"

	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/services/compassservice"
	"github.com/stretchr/testify/assert"
)

func TestDeleteComponentInput_GetQuery(t *testing.T) {
	input := &dtos.DeleteComponentInput{}
	expectedQuery := `
		mutation deleteComponent($id: ID!) {
			compass {
				deleteComponent(input: {id: $id}) {
					deletedComponentId
					errors {
						message
					}
					success
				}
			}
		}`

	assert.Equal(t, expectedQuery, input.GetQuery())
}

func TestDeleteComponentInput_SetVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    *dtos.DeleteComponentInput
		expected map[string]interface{}
	}{
		{
			name: "Valid ComponentID",
			input: &dtos.DeleteComponentInput{
				ComponentID: "123",
			},
			expected: map[string]interface{}{
				"id": "123",
			},
		},
		{
			name: "Empty ComponentID",
			input: &dtos.DeleteComponentInput{
				ComponentID: "",
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

func TestDeleteComponentOutput_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		output   *dtos.DeleteComponentOutput
		expected bool
	}{
		{
			name: "Success without errors",
			output: &dtos.DeleteComponentOutput{
				Compass: struct {
					DeleteComponent struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"deleteComponent"`
				}{
					DeleteComponent: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					}{
						Errors:  nil,
						Success: true,
					},
				},
			},
			expected: true,
		},
		{
			name: "Not found error",
			output: &dtos.DeleteComponentOutput{
				Compass: struct {
					DeleteComponent struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"deleteComponent"`
				}{
					DeleteComponent: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					}{
						Errors: []compassservice.CompassError{
							{Message: "Not Found"},
						},
						Success: false,
					},
				},
			},
			expected: false,
		},
		{
			name: "Failure with errors",
			output: &dtos.DeleteComponentOutput{
				Compass: struct {
					DeleteComponent struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"deleteComponent"`
				}{
					DeleteComponent: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					}{
						Errors: []compassservice.CompassError{
							{Message: "Some error"},
						},
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

func TestDeleteComponentOutput_GetErrors(t *testing.T) {
	tests := []struct {
		name     string
		output   *dtos.DeleteComponentOutput
		expected []string
	}{
		{
			name: "No errors",
			output: &dtos.DeleteComponentOutput{
				Compass: struct {
					DeleteComponent struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"deleteComponent"`
				}{
					DeleteComponent: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					}{
						Errors: nil,
					},
				},
			},
			expected: []string{},
		},
		{
			name: "Multiple errors",
			output: &dtos.DeleteComponentOutput{
				Compass: struct {
					DeleteComponent struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"deleteComponent"`
				}{
					DeleteComponent: struct {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.output.GetErrors())
		})
	}
}
