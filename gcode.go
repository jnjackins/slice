package slice

import (
	"bytes"
	"fmt"
)

func (l *Layer) Gcode() []byte {
	if len(l.perimeters) == 0 {
		return []byte{}
	}
	buf := new(bytes.Buffer)
	// first perimeters
	s := l.perimeters[0]
	fmt.Fprintf(buf, "G1 X%.5f Y%.5f Z%.3f E%.3f\n", s.end1.X, s.end1.Y, s.end1.Z, 0.0)
	for _, s := range l.perimeters {
		fmt.Fprintf(buf, "G1 X%.5f Y%.5f Z%.5f E%.5f\n", s.end2.X, s.end2.Y, s.end2.Z, 0.0)
	}
	for _, s := range l.infill {
		fmt.Fprintf(buf, "G1 X%.5f Y%.5f Z%.5f E%.5f\n", s.end2.X, s.end2.Y, s.end2.Z, 0.0)
	}
	return buf.Bytes()
}
