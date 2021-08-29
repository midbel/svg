package chart

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/midbel/svg"
)

type Interval struct {
	Title  string     `json:"title"`
	Starts time.Time  `json:"starts"`
	Ends   time.Time  `json:"ends"`
	Sub    []Interval `json:"children"`
}

func (i Interval) Depth() int {
	if i.isLeaf() {
		return 1
	}
	var d int
	for j := range i.Sub {
		x := i.Sub[j].Depth()
		if x > d {
			d = x
		}
	}
	return d + 1
}

func (i Interval) isLeaf() bool {
	return len(i.Sub) == 0
}

type GanttSerie struct {
	Title  string
	values []Interval

	timepair
	svg.Fill
}

func NewGanttSerie(title string) GanttSerie {
	return GanttSerie{
		Title: title,
	}
}

func (g *GanttSerie) Append(i Interval) {
	if len(g.values) == 0 {
		g.timepair.Min = i.Starts
		g.timepair.Max = i.Ends
	}
	if g.timepair.Min.IsZero() || i.Starts.Before(g.timepair.Min) {
		g.timepair.Min = i.Starts
	}
	if g.timepair.Max.IsZero() || i.Ends.After(g.timepair.Max) {
		g.timepair.Max = i.Ends
	}
	g.values = append(g.values, i)
}

func (g *GanttSerie) Depth() int {
	var d int
	for i := range g.values {
		x := g.values[i].Depth()
		if x > d {
			d = x
		}
	}
	return d
}

func (g *GanttSerie) Len() int {
	return len(g.values)
}

type TimeSerie struct {
	Title string
	svg.Stroke

	values []timepoint
	px     timepair
	py     pair
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
	GanttAxis
}

func (c GanttChart) Render(w io.Writer, series []GanttSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	cs := c.RenderElement(series)
	cs.Render(ws)
}

func (c GanttChart) RenderElement(series []GanttSerie) svg.Element {
	var (
		cs     = c.getCanvas()
		area   = c.getArea(whitstrok.Option())
		height = c.GetAreaHeight() / float64(len(series))
		rx, ds = getGanttDomains(series)
		bar    = height / float64(getMaxGanttDepth(series)) * 0.6
	)
	rx = rx.extendBy(time.Hour * 4)
	cs.Append(c.GanttAxis.drawAxis(c.Chart, rx, ds))
	for i := range series {
		var (
			depth = series[i].Depth()
			zone  = height / float64(depth)
			grp   = svg.NewGroup(svg.WithTranslate(0, float64(i)*height))
		)
		c.drawSerie(&grp, series[i], rx, height, bar, zone, 0)
		area.Append(grp.AsElement())
	}
	cs.Append(area.AsElement())
	return cs.AsElement()
}

func (c GanttChart) drawSerie(a appender, serie GanttSerie, rx timepair, height, bar, part, level float64) {
	dx := c.GetAreaWidth() / rx.Diff()
	for i, v := range serie.values {
		var (
			x0 = v.Starts.Sub(rx.Min).Seconds() * dx
			x1 = v.Ends.Sub(rx.Min).Seconds() * dx
			y0 = (height / 2) - (bar / 2)
		)
		if !v.isLeaf() || level > 0 {
			y0 = part * level
			y0 += (part / 2) - (bar / 2)
		}
		var (
			p = svg.NewPos(x0, y0)
			d = svg.NewDim(x1-x0, bar)
			r = svg.NewRect(p.Option(), d.Option(), serie.Fill.Option())
		)
		r.Title = fmt.Sprintf("%s (%s - %s)", serie.values[i].Title, v.Starts.Format(time.RFC3339), v.Ends.Format(time.RFC3339))
		a.Append(r.AsElement())

		sx := GanttSerie{values: v.Sub}
		sx.Fill = serie.Fill
		c.drawSerie(a, sx, rx, height, bar, part, level+1)
	}
}

type CalendarChart struct {
	Chart
}

func (c CalendarChart) Render(w io.Writer, series []TimeSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	cs := c.RenderElement(series)
	cs.Render(ws)
}

func (c CalendarChart) RenderElement(series []TimeSerie) svg.Element {
	return nil
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
		cs     = c.getCanvas()
		area   = c.getArea()
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
		pat = svg.NewPath(s.Stroke.Option(), nonefill.Option())
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

type timepoint struct {
	X time.Time
	Y float64
}

type timepair struct {
	Min time.Time
	Max time.Time
}

func (t timepair) extendBy(by time.Duration) timepair {
	t.Min = t.Min.Add(-by)
	t.Max = t.Max.Add(by)
	return t
}

func (t timepair) Diff() float64 {
	diff := t.Max.Sub(t.Min)
	return diff.Seconds()
}

func getMaxGanttDepth(series []GanttSerie) int {
	var d int
	for i := range series {
		x := series[i].Depth()
		if x > d {
			d = x
		}
	}
	return d
}

func getGanttDomains(series []GanttSerie) (timepair, []string) {
	var (
		p   timepair
		str []string
	)
	for i := range series {
		str = append(str, series[i].Title)
		if i == 0 || p.Min.After(series[i].timepair.Min) {
			p.Min = series[i].timepair.Min
		}
		if i == 0 || p.Max.Before(series[i].timepair.Max) {
			p.Max = series[i].timepair.Max
		}
	}
	return p, str
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
