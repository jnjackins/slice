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

	// slice each facet to find the perimeters
	segments := make([]*segment, 0, len(facets))
	zs := segment{}
	for _, f := range facets {
		s := sliceFacet(f, z)
		if s != zs {
			segments = append(segments, &s)
		}
	}
	dprintf("layer z=%0.3f: %d facets / %d perimeter segments", z, len(facets), len(segments))

	l := &Layer{
		n:          n,
		z:          z,
		stl:        s,
		facets:     facets,
		perimeters: segments,
	}

	l.genInfill(cfg)

	return l
}

//TODO: case where segment is one of the edges of the triangle
func sliceFacet(f *facet, z float64) segment {
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
		dprintf("warning: no intersections at z=%f", z)
		return segment{}
	} else if i != 2 {
		dprintf("warning: found %d intersections when finding segment at z=%f", i, z)
		return segment{}
	}

	dprintf("sliced segment: {%v - %v}", ends[0], ends[1])
	return segment{from: ends[0], to: ends[1]}
}
