package stl

import (
	"math"

	"sigint.ca/slice/vector"
)

type Solid struct {
	Facets []Facet

	min, max vector.V3
}

func (s *Solid) Bounds() (min, max vector.V3) {
	return s.min, s.max
}

func (s *Solid) updateBounds() {
	small := math.Inf(-1)
	big := math.Inf(+1)
	s.min, s.max = vector.V3{big, big, big}, vector.V3{small, small, small}

	for _, f := range s.Facets {
		for _, v := range f.Vertices {
			if v.X < s.min.X {
				s.min.X = v.X
			}
			if v.X > s.max.X {
				s.max.X = v.X
			}
			if v.Y < s.min.Y {
				s.min.Y = v.Y
			}
			if v.Y > s.max.Y {
				s.max.Y = v.Y
			}
			if v.Z < s.min.Z {
				s.min.Z = v.Z
			}
			if v.Z > s.max.Z {
				s.max.Z = v.Z
			}
		}
	}
}
