package slice

// An Infiller Fills a perimeter.
type Infiller interface {
	Fill(*Region)
}
