package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/midbel/svg"
)

const (
	defaultWidth  = 640
	defaultHeight = 480
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Point struct {
	X float64
	Y float64
}

func makePoint(x, y float64) Point {
	return Point{
		X: x,
		Y: y,
	}
}

func main() {
	var (
		padding = flag.Float64("p", 50, "padding")
		count   = flag.Int("c", 100, "number of points")
		dwidth  = flag.Float64("x", defaultWidth, "width")
		dheight = flag.Float64("y", defaultHeight, "height")
		limit   = flag.Int("i", 100, "limit")
	)
	flag.Parse()
	var s Serie
	for i := 0; i < *count; i++ {
		s.Add(float64(rand.Intn(*limit)), float64(rand.Intn(*limit)))
	}
	dx, dy := s.Diff()

	var points []Point
	for i := 0; i < *count; i++ {
		x := (s.X[i] / dx) * (*dwidth)
		y := (s.Y[i] / dy) * (*dheight - (*padding * 2))
		points = append(points, makePoint(x, (*dheight-(*padding*2))-y))
	}

	var (
		dim = svg.NewDim(*dwidth, *dheight)
		cs  = svg.NewSVG(dim.Option())
		ws  = bufio.NewWriter(os.Stdout)
	)
	defer ws.Flush()

	horiz := svg.NewGroup(svg.WithTranslate(*padding, *padding))
	band := (*dheight - (*padding * 2)) / 5
	for i := 1; i < 5; i++ {
		h := float64(i) * band
		s := svg.NewStroke("lightgrey", 0)
		s.Width = 0.5
		s.Dash.Array = []int{10, 5}
		line := svg.NewLine(svg.NewPos(0, h), svg.NewPos(*dwidth-*padding, h), s.Option())
		horiz.Append(line.AsElement())
	}

	axis := []svg.Option{
		svg.WithStroke(svg.NewStroke("black", 1)),
	}
	xaxis := svg.NewLine(svg.NewPos(0, *dheight-*padding), svg.NewPos(*dwidth-*padding, *dheight-*padding), axis...)
	horiz.Append(xaxis.AsElement())
	yaxis := svg.NewLine(svg.NewPos(0, 0), svg.NewPos(0, *dheight-*padding), axis...)
	horiz.Append(yaxis.AsElement())
	cs.Append(horiz.AsElement())

	fill := svg.NewFill("transparent")
	fill.Opacity = 0
	options := []svg.Option{
		svg.NewStroke("blue", 2).Option(),
		fill.Option(),
	}
	pat := svg.NewPath(options...)
	for i, p := range points {
		if i == 0 {
			pat.AbsMoveTo(svg.NewPos(p.X, p.Y))
			continue
		}
		pat.AbsLineTo(svg.NewPos(p.X, p.Y))
	}
	area := svg.NewGroup(svg.WithTranslate(*padding, *padding))
	area.Append(pat.AsElement())

	for i, p := range points {
		options := []svg.Option{
			svg.WithRadius(4),
			svg.WithFill(svg.NewFill("white")),
			svg.WithPosition(p.X, p.Y),
			svg.WithStroke(svg.NewStroke("blue", 2)),
		}
		ci := svg.NewCircle(options...)
		ci.Title = fmt.Sprintf("data: (%.0f, %.0f), coord(%.0f, %.0f)", s.X[i], s.Y[i], p.X, p.Y)
		area.Append(ci.AsElement())
	}
	cs.Append(area.AsElement())
	cs.Render(ws)
}

type Serie struct {
	Y []float64
	X []float64
}

func (s *Serie) Diff() (float64, float64) {
	dx := s.X[len(s.X)-1] - s.X[0]
	dy := s.X[len(s.Y)-1] - s.Y[0]
	return dx, dy
}

func (s *Serie) Len() int {
	return len(s.X)
}

func (s *Serie) Add(x, y float64) {
	s.X = append(s.X, x)
	s.Y = append(s.Y, y)
	sort.Float64s(s.X)
	sort.Float64s(s.Y)
}
