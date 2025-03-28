package list_test

import (
	"testing"

	"github.com/motain/of-catalog/internal/utils/list"
	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		list     []string
		element  string
		expected bool
	}{
		{
			name:     "element exists in list",
			list:     []string{"apple", "banana", "cherry"},
			element:  "banana",
			expected: true,
		},
		{
			name:     "element does not exist in list",
			list:     []string{"apple", "banana", "cherry"},
			element:  "grape",
			expected: false,
		},
		{
			name:     "empty list",
			list:     []string{},
			element:  "apple",
			expected: false,
		},
		{
			name:     "list with one matching element",
			list:     []string{"apple"},
			element:  "apple",
			expected: true,
		},
		{
			name:     "list with one non-matching element",
			list:     []string{"banana"},
			element:  "apple",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := list.Contains(tt.list, tt.element)
			assert.Equal(t, tt.expected, result)
		})
	}
}
