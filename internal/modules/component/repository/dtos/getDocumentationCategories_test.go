package dtos_test

import (
	"reflect"
	"testing"

	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
)

func TestDocumentationCategoriesInput_GetQuery(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.DocumentationCategoriesInput
		want string
	}{
		{
			name: "returns correct query string",
			dto:  dtos.DocumentationCategoriesInput{},
			want: `
		query documentationCategories {
			compass {
				documentationCategories(cloudId: "fca6a80f-888b-4079-82e6-3c2f61c788e2") @optIn(to: "compass-beta")  {
					... on CompassDocumentationCategoriesConnection {
						nodes {
							name
							id
							description
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

func TestDocumentationCategoriesInput_SetVariables(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.DocumentationCategoriesInput
		want map[string]interface{}
	}{
		{
			name: "returns correct variables map",
			dto:  dtos.DocumentationCategoriesInput{CompassCloudID: "test-cloud-id"},
			want: map[string]interface{}{"cloudId": "test-cloud-id"},
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

func TestDocumentationCategoriesOutput_IsSuccessful(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.DocumentationCategoriesOutput
		want bool
	}{
		{
			name: "returns true when nodes are not nil",
			dto: dtos.DocumentationCategoriesOutput{
				Compass: struct {
					DocumentationCategories struct {
						Nodes []struct {
							ID   string `json:"id"`
							Name string `json:"name"`
						} `json:"nodes"`
					} `json:"documentationCategories"`
				}{
					DocumentationCategories: struct {
						Nodes []struct {
							ID   string `json:"id"`
							Name string `json:"name"`
						} `json:"nodes"`
					}{
						Nodes: []struct {
							ID   string `json:"id"`
							Name string `json:"name"`
						}{
							{ID: "1", Name: "Category1"},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "returns false when nodes are nil",
			dto:  dtos.DocumentationCategoriesOutput{},
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

func TestDocumentationCategoriesOutput_GetErrors(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.DocumentationCategoriesOutput
		want []string
	}{
		{
			name: "returns nil as there are no errors",
			dto:  dtos.DocumentationCategoriesOutput{},
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
