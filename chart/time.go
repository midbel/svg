package chart

import (
	"bufio"
	"io"
	"math"
	"time"

	"github.com/midbel/svg"
)

type timepoint struct {
	X time.Time
	Y float64
}

type TimeSerie struct {
	Title  string
	values []timepoint
	svg.Stroke

	px    timepair
	py    pair
}

func NewTimeSerie(title string) TimeSerie {
	return TimeSerie{
		Title: title,
	}
}

func (ir *TimeSerie) Add(x time.Time, y float64) {
	if len(ir.values) == 0 {
		ir.px.Min = x
		ir.px.Max = x
		ir.py.Min = y
		ir.py.Max = y
	}
	if ir.px.Min.IsZero() || x.Before(ir.px.Min) {
		ir.px.Min = x
	}
	if ir.px.Max.IsZero() || x.After(ir.px.Max) {
		ir.px.Max = x
	}
	ir.py.Min = getLesser(ir.py.Min, y)
	ir.py.Max = getGreater(ir.py.Max, y)
	vp := timepoint{
		X: x,
		Y: y,
	}
	ir.values = append(ir.values, vp)
}

func (ir *TimeSerie) Len() int {
	return len(ir.values)
}

type GanttChart struct {
	Chart
}

type ContribChart struct {
	Chart
}

func (c ContribChart) Render(w io.Writer, series []TimeSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	cs := c.RenderElement(series)
	cs.Render(ws)
}

func (c ContribChart) RenderElement(series []TimeSerie) svg.Element {
	c.checkDefault()

	var (
		dim  = svg.NewDim(c.Width, c.Height)
		cs   = svg.NewSVG(dim.Option())
		area = svg.NewGroup(svg.WithID("area"), c.translate())
	)
	cs.Append(area.AsElement())
	return cs.AsElement()
}

type TimeChart struct {
	Chart
	TimeAxis
}

func (c TimeChart) Render(w io.Writer, series []TimeSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()

	cs := c.RenderElement(series)
	cs.Render(ws)
}

func (c TimeChart) RenderElement(series []TimeSerie) svg.Element {
	c.checkDefault()

	var (
		cs = c.getCanvas()
		area = c.getArea()
		rx, ry = getTimeDomains(series)
	)
	ry = ry.extendBy(1.2)
	cs.Append(c.drawAxis(c.Chart, rx, ry))
	for i := range series {
		elem := c.drawSerie(series[i], rx, ry)
		area.Append(elem)
	}
	cs.Append(area.AsElement())
	return cs.AsElement()
}

func (c TimeChart) drawSerie(s TimeSerie, px timepair, py pair) svg.Element {
	var (
		wx  = c.GetAreaWidth() / px.Diff()
		wy  = c.GetAreaHeight() / py.Diff()
		pat = getPathLine("steelblue")
		pos svg.Pos
	)
	pos.Y = c.GetAreaHeight() - (s.values[0].Y * wy)
	if py.Min < 0 {
		pos.Y -= math.Abs(py.Min) * wy
	}
	pat.AbsMoveTo(pos)
	for i := 1; i < s.Len(); i++ {
		if i > 0 {
		}
		pos.X = (s.values[i].X.Sub(s.values[0].X).Seconds()) * wx
		pos.Y = c.GetAreaHeight() - (s.values[i].Y * wy)
		if py.Min < 0 {
			pos.Y -= math.Abs(py.Min) * wy
		}
		pat.AbsLineTo(pos)
	}
	return pat.AsElement()
}

type timepair struct {
	Min time.Time
	Max time.Time
}

func (t timepair) Diff() float64 {
	diff := t.Max.Sub(t.Min)
	return diff.Seconds()
}

func getTimeDomains(series []TimeSerie) (timepair, pair) {
	var (
		tx timepair
		rx pair
	)
	for i := range series {
		if tx.Min.IsZero() || series[i].px.Min.Before(tx.Min) {
			tx.Min = series[i].px.Min
		}
		if tx.Max.IsZero() || series[i].px.Max.After(tx.Max) {
			tx.Max = series[i].px.Max
		}
		rx.Min = getLesser(series[i].py.Min, rx.Min)
		rx.Max = getGreater(series[i].py.Max, rx.Max)
	}
	return tx, rx
}
