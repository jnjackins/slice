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
	fmt.Fprintf(buf, "G1 X%.5f Y%.5f Z%.3f E%.3f\n", s.end1.x, s.end1.y, s.end1.z, 0.0)
	for _, s := range l.perimeters {
		fmt.Fprintf(buf, "G1 X%.5f Y%.5f Z%.5f E%.5f\n", s.end2.x, s.end2.y, s.end2.z, 0.0)
	}
	for _, s := range l.infill {
		fmt.Fprintf(buf, "G1 X%.5f Y%.5f Z%.5f E%.5f\n", s.end2.x, s.end2.y, s.end2.z, 0.0)
	}
	return buf.Bytes()
}
