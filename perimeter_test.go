package slice

import "testing"

func TestContains(t *testing.T) {
	perimeter := []*segment{
		{from: Vertex2{X: 0, Y: 0}, to: Vertex2{X: 0, Y: 10}},
		{from: Vertex2{X: 0, Y: 10}, to: Vertex2{X: 10, Y: 10}},
		{from: Vertex2{X: 10, Y: 10}, to: Vertex2{X: 10, Y: 0}},
		{from: Vertex2{X: 10, Y: 0}, to: Vertex2{X: 0, Y: 0}},
	}

	outside := []Vertex2{
		{X: 11, Y: 0},
		{X: -1, Y: 10},
		{X: -1, Y: 5},
		{X: 11, Y: 5},
		{X: 11, Y: -2},
		{X: 5, Y: 11},
		{X: 5, Y: -1},
	}

	inside := []Vertex2{
		{X: 1, Y: 0.1},
		{X: 5, Y: 3},
		{X: 5, Y: 0.9},
		{X: 9.9, Y: 9.9},
	}

	for _, p := range outside {
		if contains(perimeter, p) {
			t.Errorf("contains(%v, %v) == true, expected false", perimeter, p)
		}
	}

	for _, p := range inside {
		if !contains(perimeter, p) {
			t.Errorf("contains(%v, %v) == false, expected true", perimeter, p)
		}
	}
}
