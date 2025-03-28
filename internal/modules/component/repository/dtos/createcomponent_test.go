package dtos_test

import (
	"reflect"
	"testing"

	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
)

func TestCreateComponentInput_GetQuery(t *testing.T) {
	dto := &dtos.CreateComponentInput{}
	expectedQuery := `
		mutation createComponent ($cloudId: ID!, $componentDetails: CreateCompassComponentInput!) {
			compass {
				createComponent(cloudId: $cloudId, input: $componentDetails) {
					success
					componentDetails {
						id
						links {
							id
							type
							name
							url
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

func TestCreateComponentInput_SetVariables(t *testing.T) {
	tests := []struct {
		name     string
		dto      *dtos.CreateComponentInput
		expected map[string]interface{}
	}{
		{
			name: "valid input",
			dto: &dtos.CreateComponentInput{
				CompassCloudID: "cloud-123",
				Component: resources.Component{
					Name:        "Test Component",
					Slug:        "test-component",
					Description: "A test component",
					TypeID:      "type-123",
					OwnerID:     "owner-123",
					Links: []resources.Link{
						{Type: "github", Name: "Repo", URL: "https://github.com"},
					},
					Fields: map[string]interface{}{
						"field1": true,
					},
					Labels: []string{"label1"},
				},
			},
			expected: map[string]interface{}{
				"cloudId": "cloud-123",
				"componentDetails": map[string]interface{}{
					"name":        "Test Component",
					"slug":        "test-component",
					"description": "A test component",
					"typeId":      "type-123",
					"ownerId":     "owner-123",
					"fields": []map[string]interface{}{
						{
							"definition": "compass:field1",
							"value": map[string]interface{}{
								"boolean": map[string]bool{"booleanValue": true},
							},
						},
					},
					"links": []map[string]string{
						{"type": "github", "name": "Repo", "url": "https://github.com"},
					},
					"labels": []string{"label1"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dto.SetVariables(); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("SetVariables() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCreateComponentOutput_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		dto      *dtos.CreateComponentOutput
		expected bool
	}{
		{
			name: "success true",
			dto: &dtos.CreateComponentOutput{
				Compass: dtos.CompassCreatedComponentOutput{
					CreateComponent: dtos.CompassCreateComponentOutput{
						Success: true,
					},
				},
			},
			expected: true,
		},
		{
			name: "success false",
			dto: &dtos.CreateComponentOutput{
				Compass: dtos.CompassCreatedComponentOutput{
					CreateComponent: dtos.CompassCreateComponentOutput{
						Success: false,
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dto.IsSuccessful(); got != tt.expected {
				t.Errorf("IsSuccessful() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCreateComponentOutput_GetErrors(t *testing.T) {
	tests := []struct {
		name     string
		dto      *dtos.CreateComponentOutput
		expected []string
	}{
		{
			name: "no errors",
			dto: &dtos.CreateComponentOutput{
				Compass: dtos.CompassCreatedComponentOutput{
					CreateComponent: dtos.CompassCreateComponentOutput{
						Errors: []compassservice.CompassError{},
					},
				},
			},
			expected: []string{},
		},
		{
			name: "with errors",
			dto: &dtos.CreateComponentOutput{
				Compass: dtos.CompassCreatedComponentOutput{
					CreateComponent: dtos.CompassCreateComponentOutput{
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
			if got := tt.dto.GetErrors(); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("GetErrors() = %v, want %v", got, tt.expected)
			}
		})
	}
}
