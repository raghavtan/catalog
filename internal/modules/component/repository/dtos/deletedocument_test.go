package dtos_test

import (
	"reflect"
	"testing"

	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/services/compassservice"
)

func TestDeleteDocumentInput_GetQuery(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.DeleteDocumentInput
		want string
	}{
		{
			name: "Valid query",
			dto:  dtos.DeleteDocumentInput{},
			want: `
		mutation deleteDocument($input: CompassDeleteDocumentInput!) {
   		compass @optIn(to: "compass-beta") {
   			deleteDocument(input: $input) {
   				success
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

func TestDeleteDocumentInput_SetVariables(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.DeleteDocumentInput
		want map[string]interface{}
	}{
		{
			name: "Valid variables",
			dto: dtos.DeleteDocumentInput{
				ID: "comp123",
			},
			want: map[string]interface{}{
				"input": map[string]interface{}{
					"id": "comp123",
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

func TestDeleteDocumentOutput_IsSuccessful(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.DeleteDocumentOutput
		want bool
	}{
		{
			name: "Success true",
			dto: dtos.DeleteDocumentOutput{
				Compass: struct {
					DeleteDocument struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
					} `json:"deleteDocument"`
				}{
					DeleteDocument: struct {
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
			name: "Success false",
			dto: dtos.DeleteDocumentOutput{
				Compass: struct {
					DeleteDocument struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
					} `json:"deleteDocument"`
				}{
					DeleteDocument: struct {
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

func TestDeleteDocumentOutput_GetErrors(t *testing.T) {
	tests := []struct {
		name string
		dto  dtos.DeleteDocumentOutput
		want []string
	}{
		{
			name: "No errors",
			dto: dtos.DeleteDocumentOutput{
				Compass: struct {
					DeleteDocument struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
					} `json:"deleteDocument"`
				}{
					DeleteDocument: struct {
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
		{
			name: "Multiple errors",
			dto: dtos.DeleteDocumentOutput{
				Compass: struct {
					DeleteDocument struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
					} `json:"deleteDocument"`
				}{
					DeleteDocument: struct {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dto.GetErrors(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}
