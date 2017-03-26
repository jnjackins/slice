package slice

import (
	"fmt"
	"io"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// EncodeLayer compiles a layer into gcode.
func (e *Encoder) EncodeLayer(l *Layer) {
	for _, region := range l.Regions() {
		//perimeters
		s := region.Exterior[0]
		fmt.Fprintf(e.w, "G1 X%.5f Y%.5f\n", s.From.X, s.From.Y)
		for _, s := range region.Exterior {
			fmt.Fprintf(e.w, "G1 X%.5f Y%.5f E%.5f\n", s.To.X, s.To.Y, 0.0)
		}
		for _, p := range region.Interiors {
			s := p[0]
			fmt.Fprintf(e.w, "G1 X%.5f Y%.5f\n", s.From.X, s.From.Y)
			for _, s := range p {
				fmt.Fprintf(e.w, "G1 X%.5f Y%.5f E%.5f\n", s.To.X, s.To.Y, 0.0)
			}
		}

		//infill
		// TODO: non-printing moves
		for _, s := range region.Infill {
			fmt.Fprintf(e.w, "G1 X%.5f Y%.5f E%.5f\n", s.To.X, s.To.Y, 0.0)
		}
	}
}
