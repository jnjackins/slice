package slice

import (
	"math"
	"testing"
)

func TestGetIntersections(t *testing.T) {
	p := []*segment{
		{from: Vertex2{X: 0, Y: 0}, to: Vertex2{X: 0, Y: 10}},
		{from: Vertex2{X: 0, Y: 10}, to: Vertex2{X: 10, Y: 10}},
		{from: Vertex2{X: 10, Y: 10}, to: Vertex2{X: 10, Y: 0}},
		{from: Vertex2{X: 10, Y: 0}, to: Vertex2{X: 0, Y: 0}},
	}
	testRays := []struct {
		ray *segment
		n   int
	}{
		{ray: &segment{from: Vertex2{X: -math.MaxFloat64, Y: 5}, to: Vertex2{X: 5, Y: 5}}, n: 1},
		{ray: &segment{from: Vertex2{X: math.MaxFloat64, Y: 5}, to: Vertex2{X: -1, Y: 5}}, n: 2},
		{ray: &segment{from: Vertex2{X: -10, Y: 0}, to: Vertex2{X: 5, Y: -5}}, n: 0},
	}

	for _, r := range testRays {
		intersections, _ := r.ray.getIntersections(p)
		if len(intersections) != r.n {
			t.Errorf("%v.getIntersections(%v) == %d, expected %d", r.ray, p, len(intersections), r.n)
		}
	}
}
