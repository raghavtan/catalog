package dtos_test

import (
	"reflect"
	"testing"

	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
)

func TestUpdateDocumentInput_GetQuery(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.UpdateDocumentInput
		want string
	}{
		{
			name: "returns correct query",
			dto:  dtos.UpdateDocumentInput{},
			want: `
		mutation updateDocument($input: CompassUpdateDocumentInput!) {
		compass @optIn(to: "compass-beta") {
			updateDocument(input: $input) {
				success
				errors {
					message
				}
				documentDetails {
					id
					title
					url
					componentId
					documentationCategoryId
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

func TestUpdateDocumentInput_SetVariables(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.UpdateDocumentInput
		want map[string]interface{}
	}{
		{
			name: "returns correct variables",
			dto: dtos.UpdateDocumentInput{
				Document: resources.Document{
					ID:    "doc123",
					Title: "Test Document",
					URL:   "http://example.com",
				},
				CategoryID: "cat456",
			},
			want: map[string]interface{}{
				"input": map[string]interface{}{
					"id":                      "doc123",
					"title":                   "Test Document",
					"documentationCategoryId": "cat456",
					"url":                     "http://example.com",
				},
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

func TestUpdateDocumentOutput_IsSuccessful(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.UpdateDocumentOutput
		want bool
	}{
		{
			name: "returns true when successful",
			dto: dtos.UpdateDocumentOutput{
				Compass: struct {
					UpdateDocument struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
					} `json:"updateDocument"`
				}{
					UpdateDocument: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
					}{
						Success: true,
					},
				},
			},
			want: true,
		},
		{
			name: "returns false when not successful",
			dto: dtos.UpdateDocumentOutput{
				Compass: struct {
					UpdateDocument struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
					} `json:"updateDocument"`
				}{
					UpdateDocument: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
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

func TestUpdateDocumentOutput_GetErrors(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.UpdateDocumentOutput
		want []string
	}{
		{
			name: "returns error messages",
			dto: dtos.UpdateDocumentOutput{
				Compass: struct {
					UpdateDocument struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
					} `json:"updateDocument"`
				}{
					UpdateDocument: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
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
		{
			name: "returns empty slice when no errors",
			dto: dtos.UpdateDocumentOutput{
				Compass: struct {
					UpdateDocument struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
					} `json:"updateDocument"`
				}{
					UpdateDocument: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
					}{
						Errors: []compassservice.CompassError{},
					},
				},
			},
			want: []string{},
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
