package transformers

import "fmt"

func Interface2Float64(result interface{}) (float64, error) {
	if result == nil {
		return 0, nil
	}

	switch v := result.(type) {
	case float64:
		return v, nil
	case bool:
		return Bool2Float64(v), nil
	case string:
		return String2Float64(v)
	case []uint8:
		return String2Float64(string(v))
	default:
		fmt.Printf("unexpected result type: %T\n", v)
		fmt.Printf("result: %s\n", v)
		return 0, fmt.Errorf("unexpected result type: %T", v)
	}
}
