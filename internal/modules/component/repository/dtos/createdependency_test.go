package dtos_test

import (
	"reflect"
	"testing"

	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/services/compassservice"
)

func TestCreateDependencyInput_GetQuery(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.CreateDependencyInput
		want string
	}{
		{
			name: "valid query",
			dto:  dtos.CreateDependencyInput{},
			want: `
		mutation createRelationship($dependentId: ID!, $providerId: ID!) {
			compass {
				createRelationship(input: {
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

func TestCreateDependencyInput_SetVariables(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.CreateDependencyInput
		want map[string]interface{}
	}{
		{
			name: "valid variables",
			dto:  dtos.CreateDependencyInput{DependentId: "123", ProviderId: "456"},
			want: map[string]interface{}{
				"dependentId": "123",
				"providerId":  "456",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dto.SetVariables(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateDependencyOutput_IsSuccessful(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.CreateDependencyOutput
		want bool
	}{
		{
			name: "success true",
			dto: dtos.CreateDependencyOutput{Compass: struct {
				CreateDependency struct {
					Errors  []compassservice.CompassError `json:"errors"`
					Success bool                          `json:"success"`
				} `json:"createRelationship"`
			}{
				CreateDependency: struct {
					Errors  []compassservice.CompassError `json:"errors"`
					Success bool                          `json:"success"`
				}{Success: true},
			}},
			want: true,
		},
		{
			name: "success false",
			dto: dtos.CreateDependencyOutput{Compass: struct {
				CreateDependency struct {
					Errors  []compassservice.CompassError `json:"errors"`
					Success bool                          `json:"success"`
				} `json:"createRelationship"`
			}{
				CreateDependency: struct {
					Errors  []compassservice.CompassError `json:"errors"`
					Success bool                          `json:"success"`
				}{Success: false},
			}},
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
func TestCreateDependencyOutput_GetErrors(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.CreateDependencyOutput
		want []string
	}{
		{
			name: "no errors",
			dto: dtos.CreateDependencyOutput{Compass: struct {
				CreateDependency struct {
					Errors  []compassservice.CompassError `json:"errors"`
					Success bool                          `json:"success"`
				} `json:"createRelationship"`
			}{
				CreateDependency: struct {
					Errors  []compassservice.CompassError `json:"errors"`
					Success bool                          `json:"success"`
				}{Errors: []compassservice.CompassError{}},
			}},
			want: []string{},
		},
		{
			name: "single error",
			dto: dtos.CreateDependencyOutput{Compass: struct {
				CreateDependency struct {
					Errors  []compassservice.CompassError `json:"errors"`
					Success bool                          `json:"success"`
				} `json:"createRelationship"`
			}{
				CreateDependency: struct {
					Errors  []compassservice.CompassError `json:"errors"`
					Success bool                          `json:"success"`
				}{Errors: []compassservice.CompassError{
					{Message: "Error 1"},
				}},
			}},
			want: []string{"Error 1"},
		},
		{
			name: "multiple errors",
			dto: dtos.CreateDependencyOutput{Compass: struct {
				CreateDependency struct {
					Errors  []compassservice.CompassError `json:"errors"`
					Success bool                          `json:"success"`
				} `json:"createRelationship"`
			}{
				CreateDependency: struct {
					Errors  []compassservice.CompassError `json:"errors"`
					Success bool                          `json:"success"`
				}{Errors: []compassservice.CompassError{
					{Message: "Error 1"},
					{Message: "Error 2"},
				}},
			}},
			want: []string{"Error 1", "Error 2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dto.GetErrors(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}
