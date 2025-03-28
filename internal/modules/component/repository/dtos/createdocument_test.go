package dtos_test

import (
	"reflect"
	"testing"

	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
	"github.com/stretchr/testify/assert"
)

func TestCreateDocumentInput_GetQuery(t *testing.T) {
	tests := []struct {
		name     string
		dto      dtos.CreateDocumentInput
		expected string
	}{
		{
			name: "returns correct GraphQL query",
			dto:  dtos.CreateDocumentInput{},
			expected: `
		mutation addDocument($input: CompassAddDocumentInput!) {
   		compass @optIn(to: "compass-beta") {
   			addDocument(input: $input) {
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
			query := tt.dto.GetQuery()
			assert.Equal(t, tt.expected, query)
		})
	}
}

func TestCreateDocumentInput_SetVariables(t *testing.T) {
	tests := []struct {
		name     string
		dto      *dtos.CreateDocumentInput
		expected map[string]interface{}
	}{
		{
			name: "valid input",
			dto: &dtos.CreateDocumentInput{
				ComponentID: "component-123",
				Document: resources.Document{
					Title: "Test Document",
					URL:   "http://example.com",
				},
				CategoryID: "category-456",
			},
			expected: map[string]interface{}{
				"input": map[string]interface{}{
					"componentId":             "component-123",
					"title":                   "Test Document",
					"documentationCategoryId": "category-456",
					"url":                     "http://example.com",
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

func TestCreateDocumentOutput_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		dto      *dtos.CreateDocumentOutput
		expected bool
	}{
		{
			name: "success true",
			dto: &dtos.CreateDocumentOutput{
				Compass: struct {
					AddDocument struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
					} `json:"addDocument"`
				}{
					AddDocument: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
					}{
						Success: true,
					},
				},
			},
			expected: true,
		},
		{
			name: "success false",
			dto: &dtos.CreateDocumentOutput{
				Compass: struct {
					AddDocument struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
					} `json:"addDocument"`
				}{
					AddDocument: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
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
			if got := tt.dto.IsSuccessful(); got != tt.expected {
				t.Errorf("IsSuccessful() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCreateDocumentOutput_GetErrors(t *testing.T) {
	tests := []struct {
		name     string
		dto      *dtos.CreateDocumentOutput
		expected []string
	}{
		{
			name: "with errors",
			dto: &dtos.CreateDocumentOutput{
				Compass: struct {
					AddDocument struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
					} `json:"addDocument"`
				}{
					AddDocument: struct {
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
			expected: []string{"Error 1", "Error 2"},
		},
		{
			name: "no errors",
			dto: &dtos.CreateDocumentOutput{
				Compass: struct {
					AddDocument struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
					} `json:"addDocument"`
				}{
					AddDocument: struct {
						Errors  []compassservice.CompassError `json:"errors"`
						Success bool                          `json:"success"`
						Details dtos.Document                 `json:"documentDetails"`
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
			if got := tt.dto.GetErrors(); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("GetErrors() = %v, want %v", got, tt.expected)
			}
		})
	}
}
