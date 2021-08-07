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
	min    timepoint
	max    timepoint
}

func NewTimeSerie(title string) TimeSerie {
	return TimeSerie{
		Title: title,
		min: timepoint{
			Y: math.NaN(),
		},
		max: timepoint{
			Y: math.NaN(),
		},
	}
}

func (ir *TimeSerie) Add(x time.Time, y float64) {
	vp := timepoint{
		X: x,
		Y: y,
	}
	ir.values = append(ir.values, vp)
}

type ContribChart struct {
	Chart
}

func (c ContribChart) Render(w io.Writer, series []TimeSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	c.render(ws, series)
}

func (c ContribChart) render(w svg.Writer, series []TimeSerie) {
	c.checkDefault()

	var (
		dim  = svg.NewDim(c.Width, c.Height)
		cs   = svg.NewSVG(dim.Option())
		area = svg.NewGroup(svg.WithID("area"), c.translate())
	)
	cs.Append(area.AsElement())
	cs.Render(w)
}

type TimeChart struct {
	Chart
	TicksY int
	TicksX int
}

func (c TimeChart) Render(w io.Writer, series []TimeSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	c.render(ws, series)
}

func (c TimeChart) render(w svg.Writer, series []TimeSerie) {
	c.checkDefault()

	var (
		dim  = svg.NewDim(c.Width, c.Height)
		cs   = svg.NewSVG(dim.Option())
		area = svg.NewGroup(svg.WithID("area"), c.translate())
	)
	cs.Append(area.AsElement())
	cs.Render(w)
}

func (c TimeChart) drawSerie(s TimeSerie) svg.Element {
	return nil
}

func (c TimeChart) drawAxisX() svg.Element {
	return nil
}

func (c TimeChart) drawAxisY() svg.Element {
	return nil
}

func (c TimeChart) drawTicks() svg.Element {
	return nil
}

func getTimeDomains(series []TimeSerie) {

}
