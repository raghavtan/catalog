package dtos_test

import (
	"testing"

	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/services/compassservice"
)

func TestDeleteDependencyInput_GetQuery(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.DeleteDependencyInput
		want string
	}{
		{
			name: "valid query",
			dto:  dtos.DeleteDependencyInput{},
			want: `
		mutation deleteRelationship($dependentId: ID!, $providerId: ID!) {
			compass {
				deleteRelationship(input: {
					type: DEPENDS_ON,
					startNodeId: $dependentId,
					endNodeId: $providerId
				}) {
					errors {
						message
					}
					success
				}
			}
		}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dto.GetQuery(); got != tt.want {
				t.Errorf("GetQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestDeleteDependencyInput_SetVariables(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.DeleteDependencyInput
		want map[string]interface{}
	}{
		{
			name: "valid variables",
			dto: dtos.DeleteDependencyInput{
				DependentId: "123",
				ProviderId:  "456",
			},
			want: map[string]interface{}{
				"dependentId": "123",
				"providerId":  "456",
			},
		},
		{
			name: "empty variables",
			dto:  dtos.DeleteDependencyInput{},
			want: map[string]interface{}{
				"dependentId": "",
				"providerId":  "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dto.SetVariables(); !equalMaps(got, tt.want) {
				t.Errorf("SetVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteDependencyOutput_IsSuccessful(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.DeleteDependencyOutput
		want bool
	}{
		{
			name: "successful operation",
			dto: dtos.DeleteDependencyOutput{
				Compass: struct {
					DeleteDependency struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"deleteRelationship"`
				}{
					DeleteDependency: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					}{
						Success: true,
					},
				},
			},
			want: true,
		},
		{
			name: "unsuccessful operation",
			dto: dtos.DeleteDependencyOutput{
				Compass: struct {
					DeleteDependency struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"deleteRelationship"`
				}{
					DeleteDependency: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					}{
						Success: false,
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dto.IsSuccessful(); got != tt.want {
				t.Errorf("IsSuccessful() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteDependencyOutput_GetErrors(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.DeleteDependencyOutput
		want []string
	}{
		{
			name: "no errors",
			dto: dtos.DeleteDependencyOutput{
				Compass: struct {
					DeleteDependency struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"deleteRelationship"`
				}{
					DeleteDependency: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					}{
						Errors: []compassservice.CompassError{},
					},
				},
			},
			want: []string{},
		},
		{
			name: "multiple errors",
			dto: dtos.DeleteDependencyOutput{
				Compass: struct {
					DeleteDependency struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"deleteRelationship"`
				}{
					DeleteDependency: struct {
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
			want: []string{"Error 1", "Error 2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dto.GetErrors(); !equalSlices(got, tt.want) {
				t.Errorf("GetErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper functions for comparing maps and slices
func equalMaps(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
