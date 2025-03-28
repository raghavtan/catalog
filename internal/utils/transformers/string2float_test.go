package transformers_test

import (
	"testing"

	"github.com/motain/of-catalog/internal/utils/transformers"
	"github.com/stretchr/testify/assert"
)

func TestString2Float64(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  float64
		expectErr bool
	}{
		{
			name:      "valid float string",
			input:     "123.456",
			expected:  123.456,
			expectErr: false,
		},
		{
			name:      "valid integer string",
			input:     "789",
			expected:  789,
			expectErr: false,
		},
		{
			name:      "invalid string",
			input:     "abc",
			expected:  0,
			expectErr: true,
		},
		{
			name:      "empty string",
			input:     "",
			expected:  0,
			expectErr: true,
		},
		{
			name:      "string with spaces",
			input:     " 123.45 ",
			expected:  123.45,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := transformers.String2Float64(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
