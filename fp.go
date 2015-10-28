package slice

func approxEquals(v1, v2 float64) bool {
	return abs(v1-v2) < 0.00001
}

func abs(v float64) float64 {
	if v < 0 {
		return -1 * v
	}
	return v
}

func min(v1, v2 float64) float64 {
	if v1 < v2 {
		return v1
	}
	return v2
}

func max(v1, v2 float64) float64 {
	if v1 > v2 {
		return v1
	}
	return v2
}
