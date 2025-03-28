package transformers_test

import (
	"errors"
	"testing"

	"github.com/motain/of-catalog/internal/utils/transformers"
	"github.com/stretchr/testify/assert"
)

func TestInterface2Float64(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float64
		err      error
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: 0,
			err:      nil,
		},
		{
			name:     "float64 input",
			input:    42.5,
			expected: 42.5,
			err:      nil,
		},
		{
			name:     "bool true input",
			input:    true,
			expected: 1,
			err:      nil,
		},
		{
			name:     "bool false input",
			input:    false,
			expected: 0,
			err:      nil,
		},
		{
			name:     "string valid float input",
			input:    "123.45",
			expected: 123.45,
			err:      nil,
		},
		{
			name:     "string invalid float input",
			input:    "not_a_number",
			expected: 0,
			err:      errors.New("failed to parse string to float64: strconv.ParseFloat: parsing \"not_a_number\": invalid syntax"),
		},
		{
			name:     "unexpected type input",
			input:    []int{1, 2, 3},
			expected: 0,
			err:      errors.New("unexpected result type: []int"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := transformers.Interface2Float64(tt.input)
			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}
