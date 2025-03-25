package transformers

import (
	"fmt"
	"strconv"
)

func String2Float64(s string) (float64, error) {
	parsed, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse string to float64: %v", err)
	}
	return parsed, nil
}
