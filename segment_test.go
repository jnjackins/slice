package slice

import (
	"math"
	"testing"
)

func TestGetIntersections(t *testing.T) {
	p := []*Segment{
		{From: Vertex2{X: 0, Y: 0}, To: Vertex2{X: 0, Y: 10}},
		{From: Vertex2{X: 0, Y: 10}, To: Vertex2{X: 10, Y: 10}},
		{From: Vertex2{X: 10, Y: 10}, To: Vertex2{X: 10, Y: 0}},
		{From: Vertex2{X: 10, Y: 0}, To: Vertex2{X: 0, Y: 0}},
	}
	testRays := []struct {
		ray *Segment
		n   int
	}{
		{ray: &Segment{From: Vertex2{X: -math.MaxFloat64, Y: 5}, To: Vertex2{X: 5, Y: 5}}, n: 1},
		{ray: &Segment{From: Vertex2{X: math.MaxFloat64, Y: 5}, To: Vertex2{X: -1, Y: 5}}, n: 2},
		{ray: &Segment{From: Vertex2{X: -10, Y: 0}, To: Vertex2{X: 5, Y: -5}}, n: 0},
	}

	for _, r := range testRays {
		intersections, ip := r.ray.getIntersections(p)
		if len(intersections) != r.n {
			t.Errorf("%v.getIntersections(%v) == %d, expected %d (intersection point is %v", r.ray, p, len(intersections), r.n, ip)
		}
	}
}

func TestBound(t *testing.T) {
	min, max := Vertex2{X: -1, Y: -1}, Vertex2{X: 1, Y: 1}
	s := &Segment{From: Vertex2{X: -2, Y: 0}, To: Vertex2{X: 2, Y: 0}}
	l := s.getLine()
	bounded, err := l.bound(min, max)
	if err != nil {
		t.Fatalf("%v out of bounds (%v-%v)", l, min, max)
	}
	if !bounded.From.touches(Vertex2{X: -1, Y: 0}) {
		t.Errorf("expected bounded.from ~= (-1,0), got %v", bounded.From)
	}
	if !bounded.To.touches(Vertex2{X: 1, Y: 0}) {
		t.Errorf("expected bounded.to ~= (-1,0), got %v", bounded.To)
	}
}
