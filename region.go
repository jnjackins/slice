package slice

type Region struct {
	min, max  Vertex2
	Exterior  []*Segment   // Exteriors perimeter
	Interiors [][]*Segment // interior perimeters
	Infill    []*Segment   // infill lines
}
