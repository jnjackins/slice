package slice

func approxEquals(v1, v2 float64) bool {
	return abs(v1-v2) < 0.00001
}

// faster than math.Abs
func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
