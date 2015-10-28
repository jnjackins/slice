package slice

type Layer struct {
	n      int      // layer index
	z      float64  // layer z value
	stl    *STL     // the parent STL
	solids []*solid // one self-contained object, from the layer perspective
}
