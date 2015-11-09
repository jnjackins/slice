package slice

type solid struct {
	exterior  []*segment   // exterior perimeter
	interiors [][]*segment // interior perimeters
	infill    []*segment   // infill lines
	debug     []*segment   // extra lines to draw for debugging purposes
}
