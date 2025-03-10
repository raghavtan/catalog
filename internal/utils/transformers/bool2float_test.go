package transformers_test

import (
	"testing"

	"github.com/motain/of-catalog/internal/utils/transformers"
)

func TestBool2Float64(t *testing.T) {
	tests := []struct {
		name     string
		input    bool
		expected float64
	}{
		{
			name:     "true to 1.0",
			input:    true,
			expected: 1.0,
		},
		{
			name:     "false to 0.0",
			input:    false,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformers.Bool2Float64(tt.input)
			if result != tt.expected {
				t.Errorf("Bool2Float64(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
