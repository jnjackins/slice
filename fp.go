package slice

import "math"

func approxEquals(v1, v2, delta float64) bool {
	if v1 == v2 {
		return true
	}
	if math.Abs(v1-v2) < delta {
		return true
	}
	return false
}
