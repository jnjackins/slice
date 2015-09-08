package slice

import "log"

//TODO: case where segment is one of the edges of the triangle
func sliceFacet(f *facet, z float64) segment {
	var ends [3]vertex
	var i int
	v := f.vertices
	// two of these cases will usually be true
	if (v[0].z > z && v[1].z < z) || (v[0].z < z && v[1].z > z) {
		x1, x2 := v[0].x, v[1].x
		y1, y2 := v[0].y, v[1].y
		z1, z2 := v[0].z, v[1].z
		t := (z - z1) / (z2 - z1)
		x := x1 + (x2-x1)*t
		y := y1 + (y2-y1)*t
		ends[i] = vertex{x, y, z}
		i++
	}
	if (v[0].z > z && v[2].z < z) || (v[0].z < z && v[2].z > z) {
		x1, x2 := v[0].x, v[2].x
		y1, y2 := v[0].y, v[2].y
		z1, z2 := v[0].z, v[2].z
		t := (z - z1) / (z2 - z1)
		x := x1 + (x2-x1)*t
		y := y1 + (y2-y1)*t
		ends[i] = vertex{x, y, z}
		i++
	}
	if (v[1].z > z && v[2].z < z) || (v[1].z < z && v[2].z > z) {
		x1, x2 := v[1].x, v[2].x
		y1, y2 := v[1].y, v[2].y
		z1, z2 := v[1].z, v[2].z
		t := (z - z1) / (z2 - z1)
		x := x1 + (x2-x1)*t
		y := y1 + (y2-y1)*t
		ends[i] = vertex{x, y, z}
		i++
	}

	// otherwise, a segment of the facet or the entire facet should coincide with
	// the slice plane
	if i == 0 {
		//TODO
		return segment{}
	} else if i != 2 {
		log.Printf("warning: found %d intersections when finding segment at z=%f", i, z)
		return segment{}
	}

	return segment{ends[0], ends[1]}
}
