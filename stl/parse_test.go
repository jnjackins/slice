package stl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseExact(t *testing.T) {
	tests := []struct {
		path             string
		facets, min, max string
	}{
		{
			path: "testdata/cube20_ascii.stl",
			facets: `[n=(-0.0, 0.0, -1.0) v=[(20.0, 0.0, 0.0) (0.0, -20.0, 0.0) (0.0, 0.0, 0.0)]
 n=(0.0, 0.0, -1.0) v=[(0.0, -20.0, 0.0) (20.0, 0.0, 0.0) (20.0, -20.0, 0.0)]
 n=(0.0, -1.0, -0.0) v=[(20.0, -20.0, 20.0) (0.0, -20.0, 0.0) (20.0, -20.0, 0.0)]
 n=(0.0, -1.0, 0.0) v=[(0.0, -20.0, 0.0) (20.0, -20.0, 20.0) (0.0, -20.0, 20.0)]
 n=(1.0, 0.0, 0.0) v=[(20.0, 0.0, 0.0) (20.0, -20.0, 20.0) (20.0, -20.0, 0.0)]
 n=(1.0, -0.0, 0.0) v=[(20.0, -20.0, 20.0) (20.0, 0.0, 0.0) (20.0, 0.0, 20.0)]
 n=(0.0, 0.0, 1.0) v=[(20.0, -20.0, 20.0) (0.0, 0.0, 20.0) (0.0, -20.0, 20.0)]
 n=(-0.0, 0.0, 1.0) v=[(0.0, 0.0, 20.0) (20.0, -20.0, 20.0) (20.0, 0.0, 20.0)]
 n=(-1.0, -0.0, 0.0) v=[(0.0, 0.0, 20.0) (0.0, -20.0, 0.0) (0.0, -20.0, 20.0)]
 n=(-1.0, 0.0, 0.0) v=[(0.0, -20.0, 0.0) (0.0, 0.0, 20.0) (0.0, 0.0, 0.0)]
 n=(0.0, 1.0, 0.0) v=[(0.0, 0.0, 20.0) (20.0, 0.0, 0.0) (0.0, 0.0, 0.0)]
 n=(0.0, 1.0, -0.0) v=[(20.0, 0.0, 0.0) (0.0, 0.0, 20.0) (20.0, 0.0, 20.0)]
]`,
			min: "(0.0, -20.0, 0.0)",
			max: "(20.0, 0.0, 20.0)",
		},

		{
			path: "testdata/cube40_binary.stl",
			facets: `[n=(0.0, -0.0, -1.0) v=[(20.0, 20.0, 0.0) (20.0, -20.0, 0.0) (-20.0, -20.0, 0.0)]
 n=(-0.0, 0.0, -1.0) v=[(20.0, 20.0, 0.0) (-20.0, -20.0, 0.0) (-20.0, 20.0, 0.0)]
 n=(0.0, 0.0, 1.0) v=[(20.0, 20.0, 40.0) (-20.0, 20.0, 40.0) (-20.0, -20.0, 40.0)]
 n=(0.0, 0.0, 1.0) v=[(20.0, 20.0, 40.0) (-20.0, -20.0, 40.0) (20.0, -20.0, 40.0)]
 n=(1.0, -0.0, -0.0) v=[(20.0, 20.0, 0.0) (20.0, 20.0, 40.0) (20.0, -20.0, 40.0)]
 n=(1.0, 0.0, 0.0) v=[(20.0, 20.0, 0.0) (20.0, -20.0, 40.0) (20.0, -20.0, 0.0)]
 n=(-0.0, -1.0, -0.0) v=[(20.0, -20.0, 0.0) (20.0, -20.0, 40.0) (-20.0, -20.0, 40.0)]
 n=(-0.0, -1.0, 0.0) v=[(20.0, -20.0, 0.0) (-20.0, -20.0, 40.0) (-20.0, -20.0, 0.0)]
 n=(-1.0, 0.0, -0.0) v=[(-20.0, -20.0, 0.0) (-20.0, -20.0, 40.0) (-20.0, 20.0, 40.0)]
 n=(-1.0, 0.0, -0.0) v=[(-20.0, -20.0, 0.0) (-20.0, 20.0, 40.0) (-20.0, 20.0, 0.0)]
 n=(0.0, 1.0, 0.0) v=[(20.0, 20.0, 40.0) (20.0, 20.0, 0.0) (-20.0, 20.0, 0.0)]
 n=(0.0, 1.0, 0.0) v=[(20.0, 20.0, 40.0) (-20.0, 20.0, 0.0) (-20.0, 20.0, 40.0)]
]`,
			min: "(-20.0, -20.0, 0.0)",
			max: "(20.0, 20.0, 40.0)",
		},
	}

	for _, test := range tests {
		t.Logf("parsing %s", test.path)

		f, err := os.Open(test.path)
		if err != nil {
			t.Fatal(err)
		}

		s, err := Parse(f)
		if err != nil {
			t.Fatal(err)
		}

		facets := fmt.Sprint(s.Facets)
		if facets != test.facets {
			t.Errorf("%s: bad facets:\ngot=%v\nwant=%v", test.path, facets, test.facets)
		}

		min := fmt.Sprint(s.min)
		if min != test.min {
			t.Errorf("%s: bad min:\ngot=%v, want=%v", test.path, min, test.min)
		}

		max := fmt.Sprint(s.max)
		if max != test.max {
			t.Errorf("%s: bad max: got=%v, want=%v", test.path, max, test.max)
		}
		f.Close()
	}
}

func TestParseCount(t *testing.T) {
	tests := []struct {
		path     string
		nfacets  int
		min, max string
	}{
		{
			path:    "testdata/pikachu.stl",
			nfacets: 412,
			min:     "(-10.2, -11.2, -0.0)",
			max:     "(20.2, 36.0, 59.0)",
		},

		{
			path:    "testdata/3DBenchy.stl",
			nfacets: 225706,
			min:     "(-29.2, -15.5, 0.0)",
			max:     "(30.8, 15.5, 48.0)",
		},
	}

	for _, test := range tests {
		t.Logf("parsing %s", test.path)

		f, err := os.Open(test.path)
		if err != nil {
			t.Fatal(err)
		}

		s, err := Parse(f)
		if err != nil {
			t.Fatal(err)
		}

		nfacets := len(s.Facets)
		if nfacets != test.nfacets {
			t.Errorf("%s: bad nfacets: got=%d, want=%d", test.path, nfacets, test.nfacets)
		}

		min := fmt.Sprint(s.min)
		if min != test.min {
			t.Errorf("%s: bad min: got=%v, want=%v", test.path, min, test.min)
		}

		max := fmt.Sprint(s.max)
		if max != test.max {
			t.Errorf("%s: bad max: got=%v, want=%v", test.path, max, test.max)
		}
		f.Close()
	}
}

var Sink *Solid

func BenchmarkParseAsciiCube(b *testing.B) {
	buf, err := ioutil.ReadFile("testdata/cube20_ascii.stl")
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(buf)
		s, err := Parse(r)
		if err != nil {
			b.Fatal(err)
		}
		Sink = s
	}
}

func BenchmarkParsePikachu(b *testing.B) {
	buf, err := ioutil.ReadFile("testdata/Pikachu.stl")
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(buf)
		s, err := Parse(r)
		if err != nil {
			b.Fatal(err)
		}
		Sink = s
	}
}

func BenchmarkParse3DBenchy(b *testing.B) {
	buf, err := ioutil.ReadFile("testdata/3DBenchy.stl")
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		r := bytes.NewBuffer(buf)
		s, err := Parse(r)
		if err != nil {
			b.Fatal(err)
		}
		Sink = s
	}
}
