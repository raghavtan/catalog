package transformers

func Bool2Float64(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}
