package eval_test

import (
	"testing"

	"github.com/motain/of-catalog/internal/utils/eval"

	"github.com/stretchr/testify/assert"
)

func TestExpression(t *testing.T) {
	tests := []struct {
		expr     string
		expected bool
		err      bool
	}{
		{"5 >3", true, false},
		{"2< 4", true, false},
		{"7 >= 7", true, false},
		{"8 <= 10", true, false},
		{"5 == 5", true, false},
		{"6 != 7", true, false},
		{"5 > 5", false, false},
		{"3<2", false, false},
		{"7 >= 8", false, false},
		{"10 <= 9", false, false},
		{"5 == 6", false, false},
		{"7 != 7", false, false},
		{"invalid expression", false, true},
		{"5 > ", false, true},
		{" > 5", false, true},
		{"5 >> 3", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result, err := eval.Expression(tt.expr)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
