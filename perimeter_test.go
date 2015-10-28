package slice

import "testing"

func TestContains(t *testing.T) {
	stl := STL{
		Min: Vertex3{X: -10, Y: -10, Z: 0},
		Max: Vertex3{X: 20, Y: 20, Z: 0},
	}
	l := Layer{stl: &stl}

	s := new(solid)
	s.perimeters = []*segment{
		{from: Vertex2{X: 0, Y: 0}, to: Vertex2{X: 5, Y: 10}},
		{from: Vertex2{X: 5, Y: 10}, to: Vertex2{X: 10, Y: 0}},
		{from: Vertex2{X: 10, Y: 0}, to: Vertex2{X: 0, Y: 0}},
	}

	outside := []Vertex2{
		{X: -1, Y: 0},
		{X: 5, Y: 11},
		{X: 11, Y: -2},
	}

	inside := []Vertex2{
		{X: 1, Y: 0.1},
		{X: 5, Y: 3},
		{X: 5, Y: 0.9},
	}

	for _, p := range outside {
		if l.contains(s.perimeters, p) {
			t.Errorf("l.contains(%v, %v) == true, expected false", s.perimeters, p)
		}
	}

	for _, p := range inside {
		if !l.contains(s.perimeters, p) {
			t.Errorf("l.contains(%v, %v) == false, expected true", s.perimeters, p)
		}
	}
}
