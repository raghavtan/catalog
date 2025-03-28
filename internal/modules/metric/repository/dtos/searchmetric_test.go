package dtos_test

import (
	"reflect"
	"testing"

	"github.com/motain/of-catalog/internal/modules/metric/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/metric/resources"
)

func TestSearchMetricsInput_GetQuery(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.SearchMetricsInput
		want string
	}{
		{
			name: "valid query",
			dto:  dtos.SearchMetricsInput{},
			want: `
		query searchMetricDefinition($cloudId: ID!) {
			compass {
				metricDefinitions(query: {cloudId: $cloudId, first: 100}) {
					... on CompassMetricDefinitionsConnection {
						nodes{
							id
							name
						}
					}
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

func TestSearchMetricsInput_SetVariables(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.SearchMetricsInput
		want map[string]interface{}
	}{
		{
			name: "valid variables",
			dto: dtos.SearchMetricsInput{
				CompassCloudID: "cloud123",
				Metric:         resources.Metric{Name: "metricName"},
			},
			want: map[string]interface{}{
				"cloudId": "cloud123",
				"name":    "metricName",
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

func TestSearchMetricsOutput_IsSuccessful(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.SearchMetricsOutput
		want bool
	}{
		{
			name: "successful output",
			dto: dtos.SearchMetricsOutput{
				Compass: struct {
					Definitions struct {
						Nodes []dtos.Metric `json:"nodes"`
					} `json:"metricDefinitions"`
				}{
					Definitions: struct {
						Nodes []dtos.Metric `json:"nodes"`
					}{
						Nodes: []dtos.Metric{{}},
					},
				},
			},
			want: true,
		},
		{
			name: "unsuccessful output",
			dto: dtos.SearchMetricsOutput{
				Compass: struct {
					Definitions struct {
						Nodes []dtos.Metric `json:"nodes"`
					} `json:"metricDefinitions"`
				}{
					Definitions: struct {
						Nodes []dtos.Metric `json:"nodes"`
					}{
						Nodes: nil,
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

func TestSearchMetricsOutput_GetErrors(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.SearchMetricsOutput
		want []string
	}{
		{
			name: "no errors",
			dto:  dtos.SearchMetricsOutput{},
			want: nil,
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
