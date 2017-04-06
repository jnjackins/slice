package main

import (
	"flag"
	"image"
	"image/draw"
	"log"
	"os"
	"time"

	"sigint.ca/slice"
	"sigint.ca/slice/stl"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

var (
	debug = flag.Bool("d", false, "debug mode")
)

var bgcol = image.Black

func main() {
	log.SetFlags(0)
	log.SetPrefix("preview: ")
	flag.Usage = func() {
		log.Print("Usage: preview [options] file")
		flag.PrintDefaults()
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

	log.Print("parsing STL...")
	t := time.Now()
	stl, err := stl.Parse(f)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	log.Printf("parsing took %v", time.Now().Sub(t))

	layers, err := sliceSTL(stl)
	if err != nil {
		log.Fatal(err)
	}
	var (
		layer     = 0
		buf       screen.Buffer
		winSize   image.Point
		lastPaint time.Time
		dirty     = true
	)
	driver.Main(func(s screen.Screen) {
		log.Print("launching window...")
		w, err := s.NewWindow(nil)
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		// for some reason it matters that we publish before
		// the first size event
		w.Publish()
		time.Sleep(100 * time.Millisecond)

		up := func() {
			if layer < len(layers)-1 {
				layer++
			}
			dirty = true
		}
		down := func() {
			if layer > 0 {
				layer--
			}
			dirty = true
		}

		rate := time.Duration(float64(time.Second / 60))
		go func() {
			for range time.NewTicker(rate).C {
				w.Send(paint.Event{})
			}
		}()
		for {
			e := w.NextEvent()
			switch e := e.(type) {
			default:
			case mouse.Event:
				switch e.Button {
				case mouse.ButtonWheelUp:
					up()
				case mouse.ButtonWheelDown:
					down()
				}
			case key.Event:
				switch e.Code {
				case key.CodeEscape:
					log.Print("quitting")
					return
				case key.CodeUpArrow:
					if e.Direction == key.DirPress || e.Direction == key.DirNone {
						up()
					}
				case key.CodeDownArrow:
					if e.Direction == key.DirPress || e.Direction == key.DirNone {
						down()
					}
				}

			case paint.Event:
				if buf == nil {
					log.Print("fatal: unexpected nil buffer")
					os.Exit(1)
				}
				if e.External || (dirty && time.Since(lastPaint) > rate) {
					l := layers[layer]

					draw.Draw(buf.RGBA(), buf.Bounds(), bgcol, image.ZP, draw.Src)
					l.Draw(buf.RGBA())

					w.Upload(image.ZP, buf, buf.Bounds())
					w.Publish()

					dirty = false
					lastPaint = time.Now()
				}

			case size.Event:
				r := image.Point{e.WidthPx, e.HeightPx}
				if r != winSize {
					winSize = r
					dirty = true

					if buf != nil {
						buf.Release()
					}
					buf, err = s.NewBuffer(winSize)
					if err != nil {
						log.Printf("alloc screen.Buffer: %v", err)
						os.Exit(1)
					}
				}

			case error:
				log.Printf("error: %v", e)
			}
		}
	})
}

func sliceSTL(stl *stl.Solid) ([]*slice.Layer, error) {
	log.Print("slicing...")
	t := time.Now()

	var cfg = slice.Config{
		DebugMode:   true,
		LayerHeight: 0.4,
		LineWidth:   1.0,
		Infill:      &slice.Concentric{Spacing: 0.2},
	}

	layers, err := slice.Slice(stl, cfg)
	if err != nil {
		return nil, err
	}

	log.Printf("slicing took %v", time.Now().Sub(t))
	return layers, nil
}
