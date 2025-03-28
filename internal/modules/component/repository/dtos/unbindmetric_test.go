package dtos_test

import (
	"reflect"
	"testing"

	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/services/compassservice"
)

func TestUnbindMetricInput_GetQuery(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.UnbindMetricInput
		want string
	}{
		{
			name: "Valid query",
			dto:  dtos.UnbindMetricInput{},
			want: `
		mutation deleteMetricSource($id: ID!) {
			compass {
				deleteMetricSource(input: {id: $id}) {
					deletedMetricSourceId
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
				t.Errorf("UnbindMetricInput.GetQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnbindMetricInput_SetVariables(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.UnbindMetricInput
		want map[string]interface{}
	}{
		{
			name: "Valid variables",
			dto:  dtos.UnbindMetricInput{MetricID: "123"},
			want: map[string]interface{}{"id": "123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dto.SetVariables(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnbindMetricInput.SetVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnbindMetricOutput_IsSuccessful(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.UnbindMetricOutput
		want bool
	}{
		{
			name: "Success is true",
			dto: dtos.UnbindMetricOutput{
				Compass: struct {
					DeleteMetricSource struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"deleteMetricSource"`
				}{
					DeleteMetricSource: struct {
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
			name: "Success is false",
			dto: dtos.UnbindMetricOutput{
				Compass: struct {
					DeleteMetricSource struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"deleteMetricSource"`
				}{
					DeleteMetricSource: struct {
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
				t.Errorf("UnbindMetricOutput.IsSuccessful() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnbindMetricOutput_GetErrors(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.UnbindMetricOutput
		want []string
	}{
		{
			name: "No errors",
			dto: dtos.UnbindMetricOutput{
				Compass: struct {
					DeleteMetricSource struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"deleteMetricSource"`
				}{
					DeleteMetricSource: struct {
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
			name: "Multiple errors",
			dto: dtos.UnbindMetricOutput{
				Compass: struct {
					DeleteMetricSource struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
					} `json:"deleteMetricSource"`
				}{
					DeleteMetricSource: struct {
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
			if got := tt.dto.GetErrors(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnbindMetricOutput.GetErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}
