package slice

import (
	"container/list"
	"math"
)

func (s *STL) sliceLayer(n int, z float64, cfg Config) *Layer {
	dprintf("slicing layer %d...", n)
	// find the facets which interect this layer
	facets := make([]*facet, 0)
	for _, f := range s.facets {
		if f.lowZ <= z && f.highZ >= z {
			facets = append(facets, f)
		}
	}

	// first, slice all the facets
	segments := make([]*segment, 0, len(facets))
	zs := segment{}
	for _, f := range facets {
		s := sliceFacet(f, z)
		if *s != zs {
			segments = append(segments, s)
		} else {
			dprintf("discarding empty segment")
		}
	}
	dprintf("sliced %d segments", len(segments))

	if len(segments) == 0 {
		wprintf("no segments, returning empty layer")
		return &Layer{
			n:   n,
			z:   z,
			stl: s,
		}
	}

	perimeters := getPerimeters(segments)
	solids := getSolids(perimeters)

	return &Layer{
		n:      n,
		z:      z,
		stl:    s,
		solids: solids,
	}
}

//TODO: case where segment is one of the edges of the triangle
func sliceFacet(f *facet, z float64) *segment {
	var ends [3]Vertex2
	var i int
	v := f.vertices
	// two of these cases will usually be true
	if (v[0].Z > z && v[1].Z < z) || (v[0].Z < z && v[1].Z > z) {
		x1, x2 := v[0].X, v[1].X
		y1, y2 := v[0].Y, v[1].Y
		z1, z2 := v[0].Z, v[1].Z
		t := (z - z1) / (z2 - z1)
		x := x1 + (x2-x1)*t
		y := y1 + (y2-y1)*t
		ends[i] = Vertex2{x, y}
		i++
	}
	if (v[0].Z > z && v[2].Z < z) || (v[0].Z < z && v[2].Z > z) {
		x1, x2 := v[0].X, v[2].X
		y1, y2 := v[0].Y, v[2].Y
		z1, z2 := v[0].Z, v[2].Z
		t := (z - z1) / (z2 - z1)
		x := x1 + (x2-x1)*t
		y := y1 + (y2-y1)*t
		ends[i] = Vertex2{x, y}
		i++
	}
	if (v[1].Z > z && v[2].Z < z) || (v[1].Z < z && v[2].Z > z) {
		x1, x2 := v[1].X, v[2].X
		y1, y2 := v[1].Y, v[2].Y
		z1, z2 := v[1].Z, v[2].Z
		t := (z - z1) / (z2 - z1)
		x := x1 + (x2-x1)*t
		y := y1 + (y2-y1)*t
		ends[i] = Vertex2{x, y}
		i++
	}

	// otherwise, a segment of the facet or the entire facet should coincide with
	// the slice plane
	if i == 0 {
		//TODO
		wprintf("no intersections at z=%f", z)
		return &segment{}
	} else if i != 2 {
		wprintf("found %d intersections when finding segment at z=%f", i, z)
		return &segment{}
	}

	return &segment{from: ends[0], to: ends[1]}
}

// order segments into perimeters (brute force)
func getPerimeters(segments []*segment) [][]*segment {
	dprintf("finding perimeters...")

	perimeters := make([][]*segment, 0)
	var current []*segment

outer:
	for {
		if len(segments) == 0 {
			if current != nil {
				perimeters = append(perimeters, current)
			}
			break
		}
		if current == nil {
			current = make([]*segment, 1)
			current[0] = segments[0]
			segments = segments[1:]
		}
		last := len(current) - 1
		for i := 0; i < len(segments); i++ {
			if fixOrder(current[last], segments[i]) {
				current = append(current, segments[i])
				segments = append(segments[:i], segments[i+1:]...) // delete segments[i]
				continue outer
			}
		}
		dprintf("found %d segment perimeter", len(current))
		perimeters = append(perimeters, current)
		current = nil
	}
	if len(segments) != 0 {
		wprintf("getPerimeters: segments left over after ordering: %d", len(segments))
	}

	return perimeters
}

func getSolids(perimeters [][]*segment) []*solid {
	dprintf("grouping solids...")

	solids := make([]*solid, 0)
	interiors := list.New()

	// first, identify which perimeters our exterior and which are interior.
	// store exterior perimeters in new solids.
outer:
	for i := 0; i < len(perimeters); i++ {
		for j := 0; j < len(perimeters); j++ {
			if i == j {
				continue
			}
			if contains(perimeters[j], perimeters[i][0].from) {
				// perimeter i is not the outer perimeter of a solid
				interiors.PushBack(perimeters[i])
				continue outer
			}
		}
		// perimeter i is the outer perimeter of a solid
		s := solid{
			exterior:  perimeters[i],
			interiors: make([][]*segment, 0),
		}
		solids = append(solids, &s)
	}
	dprintf("found %d solids", len(solids))

	// sort interiors into their solids
	for _, s := range solids {
		p := interiors.Front()
		for p != nil {
			v := p.Value.([]*segment)
			if contains(s.exterior, v[0].from) {
				s.interiors = append(s.interiors, v)
				next := p.Next()
				interiors.Remove(p)
				p = next
			} else {
				p = p.Next()
			}
		}
	}
	if interiors.Len() != 0 {
		wprintf("%d leftover interiors", interiors.Len())
		for i := 0; i < interiors.Len(); i++ {
			p := (interiors.Remove(interiors.Front()).([]*segment))
			s := solid{
				exterior:  p,
				interiors: make([][]*segment, 0),
			}
			solids = append(solids, &s)
		}
	}

	return solids
}

// fixOrder returns true if it was able to order the segments (i.e. they are connected)
func fixOrder(first, second *segment) bool {
	if first.to.touches(second.from) {
		// perfect
		return true
	} else if first.to.touches(second.to) {
		// second is backwards
		second.from, second.to = second.to, second.from
		return true
	} else if first.from.touches(second.from) {
		// first is backwards
		first.from, first.to = first.to, first.from
		return true
	} else if first.from.touches(second.to) {
		// both are backwards
		first.from, first.to = first.to, first.from
		second.from, second.to = second.to, second.from
		return true
	}
	if first.from.near(second.from) || first.from.near(second.to) ||
		first.to.near(second.from) || first.to.near(second.to) {
	}
	return false
}

// contains returns true if v is inside perimeter, or false if v is outside perimeter
func contains(perimeter []*segment, v Vertex2) bool {
	// edge is a vertex outside of perimeter, with the same y value as v
	edge := Vertex2{X: -math.MaxFloat64, Y: v.Y}

	ray := &segment{from: edge, to: v}
	intersections, _ := ray.getIntersections(perimeter)

	// if ray crosses an odd number of segments before reaching v, then v is inside
	// perimeter. otherwise, it is outside perimeter.
	if len(intersections)%2 == 0 {
		return false
	}

	return true
}
