package slice

import "sigint.ca/slice/stl"

type Layer struct {
	n           int        // layer index
	z           float64    // layer z value
	stl         *stl.Solid // the parent STL
	regions     []*Region  // one self-contained object, from the layer perspective
	scaleFactor float64    // for drawing
}

func (l *Layer) Regions() []*Region {
	return l.regions
}
