// TODO: inner perimeters - copy and shift perimeters inward (away from their normals),
// and then trim where they intersect. shift+trim 2 segments at a time

package slice

import (
	"container/list"
	"fmt"
	"math"

	"sigint.ca/slice/stl"
	"sigint.ca/slice/vector"
)

func sliceLayer(n int, z float64, s *stl.Solid, cfg Config) *Layer {
	dprintf("slicing layer %d...", n)
	// find the facets which intersect this layer
	facets := make([]stl.Facet, 0)
	for _, f := range s.Facets {
		minz, maxz := math.Inf(+1), math.Inf(-1)
		for _, v := range f.Vertices {
			minz = math.Min(minz, v.Z)
			maxz = math.Max(maxz, v.Z)
		}
		if minz <= z && maxz >= z {
			facets = append(facets, f)
		}
	}

	// first, slice all the facets
	segments := make([]*Segment, 0, len(facets))
	for _, f := range facets {
		s := sliceFacet(f, z)
		if s == nil {
			dprintf("discarding nil Segment")
		} else if s.From.touches(s.To) {
			dprintf("discarding tiny Segment")
		} else {
			segments = append(segments, s)
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

	l := &Layer{
		n:   n,
		z:   z,
		stl: s,
	}

	l.regions = getRegions(getPerimeters(segments))

	return l
}

func sliceFacet(f stl.Facet, z float64) *Segment {
	norm := vector.V2{X: f.Normal.X, Y: f.Normal.Y}
	norm = norm.Normalize()

	var ends [3]Vertex2
	var i int
	v := f.Vertices

	// special case: one or more of the vertices lies
	// exactly on the slice plane
	if v[0].Z == z {
		ends[i] = Vertex2{v[0].X, v[0].Y}
		i++
	}
	if v[1].Z == z {
		ends[i] = Vertex2{v[1].X, v[1].Y}
		i++
	}
	if v[2].Z == z {
		ends[i] = Vertex2{v[2].X, v[2].Y}
		i++
	}
	if i == 1 {
		return &Segment{From: ends[0], To: ends[0], Normal: norm}
	} else if i == 2 {
		return &Segment{From: ends[0], To: ends[1], Normal: norm}
	} else if i == 3 {
		dprintf("facet coincides with slice plane, ignoring")
		// the entire facet coincides with the plane.
		// no need to return any Segment; other facets
		// should be sufficient to draw the perimeter
		return nil
	}

	// two of these cases will normally be true
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

	if i != 2 {
		panic(fmt.Sprintf("facet intersects slice plane %d times at z=%f (impossible)", i, z))
	}

	return &Segment{From: ends[0], To: ends[1], Normal: norm}
}

// order segments into perimeters (brute force)
func getPerimeters(segments []*Segment) [][]*Segment {
	dprintf("finding perimeters...")

	perimeters := make([][]*Segment, 0)
	var current []*Segment

outer:
	for {
		if len(segments) == 0 {
			if current != nil {
				perimeters = append(perimeters, current)
			}
			break
		}
		if current == nil {
			current = make([]*Segment, 1)
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
		dprintf("found %d Segment perimeter", len(current))
		perimeters = append(perimeters, current)
		current = nil
	}
	if len(segments) != 0 {
		wprintf("getPerimeters: segments left over after ordering: %d", len(segments))
	}

	return perimeters
}

func getRegions(perimeters [][]*Segment) []*Region {
	dprintf("grouping regions...")

	regions := make([]*Region, 0)
	Interiors := list.New()

	// first, identify which perimeters our Exteriors and which are interior.
	// store Exteriors perimeters in new solids.
outer:
	for i := 0; i < len(perimeters); i++ {
		for j := 0; j < len(perimeters); j++ {
			if i == j {
				continue
			}
			if contains(perimeters[j], perimeters[i][0].From) {
				// perimeter i is not the outer perimeter of a solid
				Interiors.PushBack(perimeters[i])
				continue outer
			}
		}
		// perimeter i is the outer perimeter of a Region
		r := Region{
			Exterior:  perimeters[i],
			Interiors: make([][]*Segment, 0),
		}
		r.min, r.max = perimeterBounds(r.Exterior)
		regions = append(regions, &r)
	}
	dprintf("found %d regions", len(regions))

	// sort Interiors into their solids
	for _, r := range regions {
		p := Interiors.Front()
		for p != nil {
			v := p.Value.([]*Segment)
			if contains(r.Exterior, v[0].From) {
				r.Interiors = append(r.Interiors, v)
				next := p.Next()
				Interiors.Remove(p)
				p = next
			} else {
				p = p.Next()
			}
		}
	}
	if Interiors.Len() != 0 {
		wprintf("%d leftover Interiors", Interiors.Len())
	}

	return regions
}

func perimeterBounds(p []*Segment) (min, max Vertex2) {
	min = Vertex2{math.Inf(+1), math.Inf(+1)}
	max = Vertex2{math.Inf(-1), math.Inf(-1)}
	for _, s := range p {
		min.X = math.Min(min.X, math.Min(s.From.X, s.To.X))
		min.Y = math.Min(min.Y, math.Min(s.From.Y, s.To.Y))
		max.X = math.Max(max.X, math.Max(s.From.X, s.To.X))
		max.Y = math.Max(max.Y, math.Max(s.From.Y, s.To.Y))
	}
	return
}

// fixOrder returns true if it was able to order the segments (i.e. they are connected)
func fixOrder(first, second *Segment) bool {
	if first.To.touches(second.From) {
		// perfect
		return true
	} else if first.To.touches(second.To) {
		// second is backwards
		second.From, second.To = second.To, second.From
		return true
	} else if first.From.touches(second.From) {
		// first is backwards
		first.From, first.To = first.To, first.From
		return true
	} else if first.From.touches(second.To) {
		// both are backwards
		first.From, first.To = first.To, first.From
		second.From, second.To = second.To, second.From
		return true
	}
	return false
}

// contains returns true if v is inside perimeter, or false if v is outside perimeter
func contains(perimeter []*Segment, v Vertex2) bool {
	// draw a line from outside the perimeter to v. if the line
	// has an odd number of intersections with the perimeter, v
	// is inside the perimeter. if there are an odd number of
	// intersections, v lies outside the perimeter.

	edge := Vertex2{X: -math.MaxFloat64, Y: v.Y}

	ray := &Segment{From: edge, To: v}
	intersections, _ := ray.getIntersections(perimeter)

	// if ray crosses an odd number of segments before reaching v, then v is inside
	// perimeter. otherwise, it is outside perimeter.
	if len(intersections)%2 == 0 {
		return false
	}

	return true
}
