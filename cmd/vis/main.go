package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"

	"sigint.ca/slice"
)

const maskOpacity = 0xFF

var (
	imgs  []*image.RGBA
	layer int
	mask  *image.Uniform
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("vis: ")
	flag.Usage = func() {
		log.Print("Usage: vis file")
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	stl, err := slice.Parse(f)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	var cfg = slice.Config{
		LayerHeight:   0.4,
		InfillSpacing: 1.0,
		InfillAngle:   45.0,
	}

	err = stl.Slice(nil, cfg)
	if err != nil {
		log.Fatal(err)
	}

	// pre-draw all the layers

	imgs = make([]*image.RGBA, len(stl.Layers))
	mask = image.NewUniform(color.Alpha{maskOpacity})
	first := stl.Layers[0].Image()
	r := first.Bounds()

	// draw the first layer onto a plain white background
	imgs[0] = image.NewRGBA(r)
	draw.Draw(imgs[0], r, image.White, r.Min, draw.Src)
	draw.Draw(imgs[0], r, first, r.Min, draw.Over)
	for i := 1; i < len(stl.Layers); i++ {
		// draw a semi-transparent version of the previous layer
		tmp := image.NewRGBA(r)
		draw.Draw(tmp, r, imgs[i-1], r.Min, draw.Src)
		draw.Draw(tmp, r, mask, r.Min, draw.Over)

		// draw the transparent layer on white, and then draw the new layer on that.
		imgs[i] = image.NewRGBA(r)
		draw.Draw(imgs[i], r, image.White, r.Min, draw.Src)
		draw.Draw(imgs[i], r, tmp, r.Min, draw.Over)
		draw.Draw(imgs[i], r, stl.Layers[i].Image(), r.Min, draw.Over)
	}

	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(nil)
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		r := image.Rect(int(stl.Min.X), int(stl.Min.Y), int(stl.Max.X)+1, int(stl.Max.Y)+1)
		winSize := r.Size()
		var b screen.Buffer
		defer func() {
			if b != nil {
				b.Release()
			}
		}()

		var sz size.Event
		var lastClick mouse.Event

		redraw := func() {
			draw.Draw(b.RGBA(), b.RGBA().Bounds(), imgs[layer], imgs[layer].Bounds().Min, draw.Src)

			// draw layer number in top-left corner
			drawLayerNumber(b.RGBA(), layer)

		}

		for e := range w.Events() {
			switch e := e.(type) {
			default:

			case mouse.Event:
				if e.Button == mouse.ButtonLeft {
					if e.Y > lastClick.Y && layer < len(imgs)-1 {
						layer++
						redraw()
					} else if e.Y < lastClick.Y && layer > 0 {
						layer--
						redraw()
					}
					lastClick = e
				}

			case key.Event:
				if e.Code == key.CodeEscape {
					return
				}

			case paint.Event:
				w.Upload(image.Point{}, b, b.Bounds(), w)
				w.Publish()

			case screen.UploadedEvent:
				// No-op.

			case size.Event:
				sz = e
				if b != nil {
					b.Release()
				}
				winSize = image.Point{sz.WidthPx, sz.HeightPx}
				b, err = s.NewBuffer(winSize)
				if err != nil {
					log.Fatal(err)
				}
				redraw()

			case error:
				log.Printf("error: %v", e)
			}
		}
	})
}

func drawLayerNumber(dst draw.Image, n int) {
	d := font.Drawer{
		Dst:  dst,
		Src:  image.Black,
		Face: basicfont.Face7x13,
		Dot:  fixed.P(2, 13),
	}
	d.DrawString(fmt.Sprintf("Layer %03d", n))
}
