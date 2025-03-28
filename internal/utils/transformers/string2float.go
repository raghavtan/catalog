package transformers

import (
	"fmt"
	"strconv"
	"strings"
)

func String2Float64(s string) (float64, error) {
	if s == "" {
		return 0, fmt.Errorf("input string is empty")
	}

	parsed, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse string to float64: %v", err)
	}
	return parsed, nil
}
