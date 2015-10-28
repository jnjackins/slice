package slice

func (s *STL) sliceLayer(n int, z float64, cfg Config) *Layer {
	dprintf("slicing perimeters of layer %d...", n)
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
	dprintf("total for layer z=%0.3f: %d facets / %d perimeter segments", z, len(facets), len(segments))

	if len(segments) == 0 {
		wprintf("no segments, returning nil layer")
		return nil
	}

	perimeters := getPerimeters(segments)
	solids := getSolids(perimeters)

	l := &Layer{
		n:      n,
		z:      z,
		stl:    s,
		solids: solids,
	}

	//l.genInfill(cfg)

	return l
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

	dprintf("sliced segment: %v-%v", ends[0], ends[1])
	return &segment{from: ends[0], to: ends[1]}
}

// order segments into perimeters (brute force)
func getPerimeters(segments []*segment) [][]*segment {
	perimeters := make([][]*segment, 0)
	var current []*segment

	dprintf("ordering and connecting perimeter segments")
outer:
	for {
		if len(segments) == 0 {
			if current != nil {
				perimeters = append(perimeters, current)
			}
			break
		}
		if current == nil {
			dprintf("starting perimeter at %v (%d remaining)", segments[0], len(segments))
			current = make([]*segment, 1)
			current[0] = segments[0]
			segments = segments[1:]
		}
		last := len(current) - 1
		for i := 0; i < len(segments); i++ {
			if fixOrder(current[last], segments[i]) {
				current = append(current, segments[i])
				segments = append(segments[:i], segments[i+1:]...) // delete segments[j]
				dprintf("connected %v to %v (perimeter length: %d)", current[last], current[last+1], len(current))
				continue outer
			}
		}
		perimeters = append(perimeters, current)
		current = nil
	}
	if len(segments) != 0 {
		wprintf("getPerimeters: segments left over after ordering: %d", len(segments))
	}

	return perimeters
}

func getSolids(perimeters [][]*segment) []*solid {
	solids := make([]*solid, 0)
	for _, p := range perimeters {
		solids = append(solids, &solid{perimeters: p})
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
	return false
}

func (l *Layer) contains(perimeter []*segment, v Vertex2) bool {
	edge := Vertex2{X: l.stl.Min.X, Y: v.Y}
	ray := &segment{from: edge, to: v}
	intersections, _ := ray.getIntersections(perimeter)
	if len(intersections)%2 == 0 {
		return false
	}
	return true
}
