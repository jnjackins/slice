package slice

type Concentric struct {
	Spacing float64
}

func (in *Concentric) Fill(r *Region) {
	// each concentric circle is generated generated based on the previous.
	// start with the Region's exterior.
	lastRound := r.Exterior

	for round := 0; round < 1; round++ {
		dprintf("starting concentric infill round %d", round)
		roundStart := len(r.Infill)

		// shift everything inwards
		for _, s := range lastRound {
			shift := s.Normal.Mul(in.Spacing)

			infillSegment := *s
			infillSegment.ShiftBy(shift)

			r.Infill = append(r.Infill, &infillSegment)
		}

		roundEnd := len(r.Infill)
		if roundEnd == roundStart {
			// no new segments, we're done.
			break
		}
		dprintf("added %d segments in round %d", roundEnd-roundStart, round)

		lastRound = r.Infill[roundStart:roundEnd]

		// join and trim newly shifted segments.
		// there will be some overlapping, which we will eliminate in the next phase.
		in.connect(lastRound)

		// eliminate overlapping segments
		in.trim(lastRound)

		// regroup into regions
		// TODO
	}
}

func (in *Concentric) connect(segments []*Segment) {
	for i := 0; i < len(segments); i++ {
		a := segments[i]
		j := (i + 1) % len(segments)
		b := segments[j]
		aLine := a.getLine()
		bLine := b.getLine()
		intersection, err := aLine.intersect(bLine)
		if err == errNoIntersections {
			wprintf("can't connect: no intersections!")
		} else {
			a.To = intersection
			b.From = intersection
		}
	}
}

func (in *Concentric) trim(segments []*Segment) {
	// for i := 0; i < len(segments); i++ {
	// 	current := segments[i]
	// 	for {
	// 		ss, vv := current.getIntersections(segments)
	// 		if len(ss) == 0 {
	// 			break
	// 		}
	// 		if vv[0].distFrom(current.From) < vv[0].distFrom(current.To) {
	// 			ss[0].To = vv[0]
	// 			current.From = vv[0]
	// 		} else {
	// 			current.To = vv[0]
	// 			ss[0].From = vv[0]
	// 		}
	// 	}
	// }
}
