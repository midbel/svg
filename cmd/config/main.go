package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/midbel/toml"
)

type Datum struct {
	Key   string
	Value interface{}
}

type Common struct {
	Id   string
	Data []Datum
}

type Shape struct {
	Common
	Paths []Path
}

type Path struct {
	Common
}

type Group struct {
	Common
	Shapes []Shape
	Paths  []Path
}

type Config struct {
	Width  float64
	Height float64
	Shapes []Shape `toml:"shape"`
	Paths  []Path  `toml:"path"`
	Groups []Group `toml:"group"`
}

func (c Config) Render(w io.Writer) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()

	canvas := svg.NewSVG(svg.WithDim(c.Width, c.Height))
	canvas.Render(ws)
}

func main() {
	flag.Parse()
	var cfg Config
	if err := toml.DecodeFile(flag.Arg(0), &cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	fmt.Printf("%+v\n", cfg)

	cfg.Render(os.Stdout)
}

func parse(str string) interface{} {
	var (
		fn func(string) interface{}
		pf string
	)
	switch {
	case strings.HasPrefix(str, circle):
		fn = parseCircle
		pf = circle
	case strings.HasPrefix(str, rect):
		fn = parseRect
		pf = rect
	case strings.HasPrefix(str, image):
		fn = parseImage
		pf = image
	default:
		return nil
	}
	str, ok := checkBracket(strings.TrimPrefix(str, pf))
	if !ok {
		return nil
	}
	return fn(str)
}

func parseCircle(str string) interface{} {
	attrs, ok := checkAttrs(str, 1, 3)
	if !ok {
		return nil
	}
	var (
		circ   Circle
		err    error
		fields = []*float64{&circ.R, &circ.X, &circ.Y}
	)
	for i := range attrs {
		*fields[i], err = strconv.ParseFloat(strings.TrimSpace(attrs[i]), 64)
		if err != nil {
			return nil
		}
	}
	return circ
}

func parseRect(str string) interface{} {
	attrs, ok := checkAttrs(str, 2, 4)
	if !ok {
		return nil
	}
	var (
		rect   Rect
		err    error
		fields = []*float64{&rect.W, &rect.H, &rect.X, &rect.Y}
	)
	for i := range attrs {
		*fields[i], err = strconv.ParseFloat(strings.TrimSpace(attrs[i]), 64)
		if err != nil {
			return nil
		}
	}
	return rect
}

func parseImage(str string) interface{} {
	attrs, ok := checkAttrs(str, 1, 3)
	if !ok {
		return nil
	}
	var (
		err    error
		img    = Image{File: attrs[0]}
		fields = []*float64{&img.W, &img.H}
	)
	for j := range attrs[1:] {
		*fields[j], err = strconv.ParseFloat(strings.TrimSpace(attrs[j+1]), 64)
		if err != nil {
			return nil
		}
	}
	return img
}

func checkAttrs(str string, min, max int) ([]string, bool) {
	var (
		attrs = strings.Split(str, ",")
		size  = len(attrs)
	)
	return attrs, size >= min && size <= max
}

func checkBracket(str string) (string, bool) {
	var (
		left  = strings.IndexRune(str, '(')
		right = strings.IndexRune(str, ')')
	)
	if left == 0 && right == len(str)-1 {
		return str[left+1 : right], true
	}
	return "", false
}

const (
	circle = "circle"
	rect   = "rect"
	image  = "image"
)

type Circle struct {
	R float64
	X float64
	Y float64
}

type Rect struct {
	W float64
	H float64
	X float64
	Y float64
}

type Image struct {
	File string
	W    float64
	H    float64
}
