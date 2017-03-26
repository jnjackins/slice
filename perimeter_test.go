package slice

import (
	"testing"

	"sigint.ca/slice/stl"
	"sigint.ca/slice/vector"
)

func TestSliceFacet(t *testing.T) {
	f := stl.Facet{
		Vertices: [3]vector.V3{
			{X: 0, Y: 10, Z: 0},
			{X: 10, Y: 20, Z: 0},
			{X: 5, Y: 15, Z: 10},
		},
	}

	tests := []struct {
		z    float64
		want Segment
	}{
		{
			z: 0,
			want: Segment{
				From: Vertex2{X: 0, Y: 10},
				To:   Vertex2{X: 10, Y: 20},
			},
		},
		{
			z: 1,
			want: Segment{
				From: Vertex2{X: 0.5, Y: 10.5},
				To:   Vertex2{X: 9.5, Y: 19.5},
			},
		},
		{
			z: 5,
			want: Segment{
				From: Vertex2{X: 2.5, Y: 12.5},
				To:   Vertex2{X: 7.5, Y: 17.5},
			},
		},
		{
			z: 9,
			want: Segment{
				From: Vertex2{X: 4.5, Y: 14.5},
				To:   Vertex2{X: 5.5, Y: 15.5},
			},
		},
		{
			z: 10,
			want: Segment{
				From: Vertex2{X: 5, Y: 15},
				To:   Vertex2{X: 5, Y: 15},
			},
		},
	}

	for _, test := range tests {
		s := sliceFacet(f, test.z)
		if !s.From.touches(test.want.From) || !s.To.touches(test.want.To) {
			t.Errorf("bad Segment for z=%f: got %v, want %v", test.z, s, &test.want)
		}
	}
}

func TestContains(t *testing.T) {
	perimeter := []*Segment{
		{From: Vertex2{X: 0, Y: 0}, To: Vertex2{X: 0, Y: 10}},
		{From: Vertex2{X: 0, Y: 10}, To: Vertex2{X: 10, Y: 10}},
		{From: Vertex2{X: 10, Y: 10}, To: Vertex2{X: 10, Y: 0}},
		{From: Vertex2{X: 10, Y: 0}, To: Vertex2{X: 0, Y: 0}},
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
