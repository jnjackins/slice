package slice

import "math"

func round(v float64) int {
	if v > 0.0 {
		return int(v + 0.5)
	} else if v < 0.0 {
		return int(v - 0.5)
	} else {
		return 0
	}
}

func approxEquals(v1, v2, delta float64) bool {
	if v1 == v2 {
		return true
	}
	if math.Abs(v1-v2) < delta {
		return true
	}
	return false
}
