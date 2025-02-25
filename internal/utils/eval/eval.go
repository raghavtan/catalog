package eval

import (
	"fmt"
	"strconv"
	"strings"
)

func Expression(expr string) (bool, error) {
	// Trim spaces
	expr = strings.TrimSpace(expr)

	// Supported operators
	operators := []string{">=", "<=", ">", "<", "==", "!="}

	// Find the operator in the expression
	var op string
	for _, o := range operators {
		if strings.Contains(expr, o) {
			op = o
			break
		}
	}

	if op == "" {
		return false, fmt.Errorf("no valid operator found in expression")
	}

	// Split the expression into left and right operands
	parts := strings.Split(expr, op)
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid expression format")
	}

	leftStr, rightStr := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

	// Convert to float (supports both int and float inputs)
	left, err := strconv.ParseFloat(leftStr, 64)
	if err != nil {
		return false, fmt.Errorf("invalid left operand: %v", err)
	}

	right, err := strconv.ParseFloat(rightStr, 64)
	if err != nil {
		return false, fmt.Errorf("invalid right operand: %v", err)
	}

	// Evaluate based on the operator
	switch op {
	case ">":
		return left > right, nil
	case "<":
		return left < right, nil
	case ">=":
		return left >= right, nil
	case "<=":
		return left <= right, nil
	case "==":
		return left == right, nil
	case "!=":
		return left != right, nil
	default:
		return false, fmt.Errorf("unsupported operator")
	}
}
