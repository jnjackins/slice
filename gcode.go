package slice

import (
	"bytes"
	"fmt"
)

// Gcode compiles the layer into gcode.
func (l *Layer) Gcode() []byte {
	buf := new(bytes.Buffer)
	for _, solid := range l.solids {
		//perimeters
		s := solid.exterior[0]
		fmt.Fprintf(buf, "G1 X%.5f Y%.5f\n", s.from.X, s.from.Y)
		for _, s := range solid.exterior {
			fmt.Fprintf(buf, "G1 X%.5f Y%.5f E%.5f\n", s.to.X, s.to.Y, 0.0)
		}
		for _, p := range solid.interiors {
			s := p[0]
			fmt.Fprintf(buf, "G1 X%.5f Y%.5f\n", s.from.X, s.from.Y)
			for _, s := range p {
				fmt.Fprintf(buf, "G1 X%.5f Y%.5f E%.5f\n", s.to.X, s.to.Y, 0.0)
			}
		}

		//infill
		// TODO: non-printing moves
		for _, s := range solid.infill {
			fmt.Fprintf(buf, "G1 X%.5f Y%.5f E%.5f\n", s.to.X, s.to.Y, 0.0)
		}
	}
	return buf.Bytes()
}
