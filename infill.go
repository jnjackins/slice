package slice

const (
	rectolinear = iota
)

type infill struct {
	kind    int
	spacing float64
}

func (l *Layer) genInfill() {
}
