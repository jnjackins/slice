package slice

type solid struct {
	perimeters []*segment
	infill     []*segment
	debug      []*segment // extra lines to draw for debugging purposes
}
