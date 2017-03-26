package slice

import (
	"image"
	"image/draw"
	"os"
	"testing"

	"sigint.ca/slice/stl"
)

var Sink draw.Image

func BenchmarkDraw(b *testing.B) {
	f, err := os.Open("stl/testdata/pikachu.stl")
	if err != nil {
		b.Fatal(err)
	}
	stl, err := stl.Parse(f)
	f.Close()
	if err != nil {
		b.Fatal(err)
	}
	var cfg = Config{
		LayerHeight:   1.0,
		LineWidth:     1.0,
		InfillSpacing: 2.0,
		InfillAngle:   45.0,
	}
	layers, err := Slice(stl, cfg)
	if err != nil {
		b.Fatal(err)
	}

	dst := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	for i := 0; i < b.N; i++ {
		layers[i%len(layers)].Draw(dst)
	}
	Sink = dst
}
