package slice

import (
	"bytes"
	"fmt"
)

// Gcode compiles the layer into gcode.
func (l *Layer) Gcode() []byte {
	buf := new(bytes.Buffer)
	for _, solid := range l.solids {
		if len(solid.perimeters) == 0 {
			return []byte{}
		}
		// first perimeters
		s := solid.perimeters[0]
		fmt.Fprintf(buf, "G1 X%.5f Y%.5f E%.3f\n", s.from.X, s.from.Y, 0.0)
		for _, s := range solid.perimeters {
			fmt.Fprintf(buf, "G1 X%.5f Y%.5f E%.5f\n", s.to.X, s.to.Y, 0.0)
		}
		for _, s := range solid.infill {
			fmt.Fprintf(buf, "G1 X%.5f Y%.5f E%.5f\n", s.to.X, s.to.Y, 0.0)
		}
	}
	return buf.Bytes()
}
