package dtos_test

import (
	"reflect"
	"testing"

	"github.com/motain/of-catalog/internal/services/ownerservice/dtos"
)

func TestOwner(t *testing.T) {
	tests := []struct {
		name     string
		owner    dtos.Owner
		expected dtos.Owner
	}{
		{
			name: "Valid Owner with all fields",
			owner: dtos.Owner{
				OwnerID: "owner123",
				SlackChannels: map[string]string{
					"general": "C123456",
					"dev":     "C654321",
				},
				Projects: map[string]string{
					"project1": "P123",
					"project2": "P456",
				},
				DisplayName: "John Doe",
			},
			expected: dtos.Owner{
				OwnerID: "owner123",
				SlackChannels: map[string]string{
					"general": "C123456",
					"dev":     "C654321",
				},
				Projects: map[string]string{
					"project1": "P123",
					"project2": "P456",
				},
				DisplayName: "John Doe",
			},
		},
		{
			name: "Empty Owner",
			owner: dtos.Owner{
				OwnerID:       "",
				SlackChannels: nil,
				Projects:      nil,
				DisplayName:   "",
			},
			expected: dtos.Owner{
				OwnerID:       "",
				SlackChannels: nil,
				Projects:      nil,
				DisplayName:   "",
			},
		},
		{
			name: "Owner with partial fields",
			owner: dtos.Owner{
				OwnerID: "owner456",
				SlackChannels: map[string]string{
					"support": "C789123",
				},
				Projects:    nil,
				DisplayName: "Jane Smith",
			},
			expected: dtos.Owner{
				OwnerID: "owner456",
				SlackChannels: map[string]string{
					"support": "C789123",
				},
				Projects:    nil,
				DisplayName: "Jane Smith",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.owner, tt.expected) {
				t.Errorf("Owner mismatch. Got %+v, expected %+v", tt.owner, tt.expected)
			}
		})
	}
}
