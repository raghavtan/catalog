package eval

import (
	"fmt"
	"strconv"
	"strings"
)

func Expression(expr string) (bool, error) {
	expr = strings.TrimSpace(expr)
	op, err := findOperator(expr)
	if err != nil {
		return false, err
	}

	left, right, err := splitOperands(expr, op)
	if err != nil {
		return false, err
	}

	return evaluate(left, right, op)
}

func findOperator(expr string) (string, error) {
	operators := []string{">=", "<=", ">", "<", "==", "!="}
	for _, o := range operators {
		if strings.Contains(expr, o) {
			return o, nil
		}
	}
	return "", fmt.Errorf("no valid operator found in expression")
}

func splitOperands(expr, op string) (float64, float64, error) {
	parts := strings.Split(expr, op)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid expression format")
	}

	leftStr, rightStr := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	left, err := strconv.ParseFloat(leftStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid left operand: %v", err)
	}

	right, err := strconv.ParseFloat(rightStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid right operand: %v", err)
	}

	return left, right, nil
}

func evaluate(left, right float64, op string) (bool, error) {
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
