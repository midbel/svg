package chart

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/midbel/svg"
)

type Interval struct {
	Title  string     `json:"title"`
	Starts time.Time  `json:"starts"`
	Ends   time.Time  `json:"ends"`
	Sub    []Interval `json:"children"`

	svg.Fill
}

func (i Interval) Range() (time.Time, time.Time) {
	if i.isLeaf() {
		return i.Starts, i.Ends
	}
	starts, ends := i.Starts, i.Ends
	for j := range i.Sub {
		s, e := i.Sub[j].Range()
		if s.Before(starts) {
			starts = s
		}
		if e.After(ends) {
			ends = e
		}
	}
	return starts, ends
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

func (i Interval) IsZero() bool {
	return i.Starts.IsZero() && i.Ends.IsZero()
}

type TimeSerie struct {
	Title string
	svg.Stroke
	svg.Fill
	Curver

	values []timepoint
	px     timepair
	py     pair
}

func NewTimeSerie(title string) TimeSerie {
	return TimeSerie{
		Title: title,
	}
}

func (ir *TimeSerie) At(i int) Point {
	var p Point
	if i >= 0 && i < len(ir.values) {
		p.X = float64(ir.values[i].X.Unix())
		p.Y = ir.values[i].Y
	}
	return p
}

func (ir *TimeSerie) GetStroke() svg.Stroke {
	return ir.Stroke
}

func (ir *TimeSerie) GetFill() svg.Fill {
	return ir.Fill
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

func (c GanttChart) Render(w io.Writer, series []Interval) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	cs := c.RenderElement(series)
	cs.Render(ws)
}

func (c GanttChart) RenderElement(series []Interval) svg.Element {
	for len(series) == 1 {
		series = series[0].Sub
	}
	var (
		cs     = c.getCanvas()
		area   = c.getArea(whitstrok.Option())
		offset = c.GetAreaHeight() / float64(len(series))
		rx, ds = getIntervalDomains(series)
		bar    = offset / float64(getIntervalDepth(series)) * 0.6
	)
	rx = rx.extendBy(time.Hour)

	cs.Append(c.Chart.drawAxis(rx.AxisRange(), WithLabels(ds...)))

	for i := range series {
		var (
			height = offset / float64(series[i].Depth())
			grp    = svg.NewGroup(svg.WithTranslate(0, float64(i)*offset))
		)
		series[i].Fill = getFill(i, series[i].Fill, series[i].Fill)
		c.drawInterval(&grp, series[i], rx, bar, height, 0)
		area.Append(grp.AsElement())
	}
	cs.Append(area.AsElement())
	cs.Append(c.drawTitle())
	cs.Append(c.drawLegend())
	return cs.AsElement()
}

func (c GanttChart) drawInterval(a Appender, serie Interval, rx timepair, bar, height, level float64) {
	if serie.IsZero() {
		return
	}
	var (
		dx = c.GetAreaWidth() / rx.Diff()
		x0 = serie.Starts.Sub(rx.Min).Seconds() * dx
		x1 = serie.Ends.Sub(rx.Min).Seconds() * dx
		p  = svg.NewPos(x0, (height*level)+(height/2)-(bar/2))
		d  = svg.NewDim(x1-x0, bar)
		r  = svg.NewRect(p.Option(), d.Option(), serie.Fill.Option())
	)
	r.Title = fmt.Sprintf("%s (%s - %s)", serie.Title, serie.Starts.Format(time.RFC3339), serie.Ends.Format(time.RFC3339))
	a.Append(r.AsElement())
	for i := range serie.Sub {
		if serie.Sub[i].Fill.IsZero() {
			serie.Sub[i].Fill = serie.Fill
		}
		c.drawInterval(a, serie.Sub[i], rx, bar, height, level)
		level++
	}
}

type IntervalChart struct {
	Chart
}

func (c IntervalChart) Render(w io.Writer, series []Interval) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	cs := c.RenderElement(series)
	cs.Render(ws)
}

func (c IntervalChart) RenderElement(series []Interval) svg.Element {
	c.checkDefault()

	for len(series) == 1 {
		series = series[0].Sub
	}

	var (
		cs     = c.getCanvas()
		area   = c.getArea(whitstrok.Option())
		offset = c.GetAreaHeight() / float64(len(series))
		rx, ds = getIntervalDomains(series)
		bar    = offset / float64(getIntervalDepth(series)) * 0.6
	)
	rx = rx.extendBy(time.Hour)

	cs.Append(c.Chart.drawAxis(rx.AxisRange(), WithLabels(ds...)))

	for i := range series {
		var (
			height = offset / float64(series[i].Depth())
			grp    = svg.NewGroup(svg.WithTranslate(0, float64(i)*offset))
		)
		series[i].Fill = getFill(i, series[i].Fill, series[i].Fill)
		c.drawInterval(&grp, series[i], rx, bar, height, 0)
		area.Append(grp.AsElement())
	}
	cs.Append(area.AsElement())
	cs.Append(c.drawTitle())
	cs.Append(c.drawLegend())
	return cs.AsElement()
}

func (c IntervalChart) drawInterval(a Appender, serie Interval, rx timepair, bar, height, level float64) {
	if !serie.IsZero() {
		var (
			dx = c.GetAreaWidth() / rx.Diff()
			x0 = serie.Starts.Sub(rx.Min).Seconds() * dx
			x1 = serie.Ends.Sub(rx.Min).Seconds() * dx
			p  = svg.NewPos(x0, (height*level)+(height/2)-(bar/2))
			d  = svg.NewDim(x1-x0, bar)
			r  = svg.NewRect(p.Option(), d.Option(), serie.Fill.Option())
		)
		r.Title = fmt.Sprintf("%s (%s - %s)", serie.Title, serie.Starts.Format(time.RFC3339), serie.Ends.Format(time.RFC3339))
		a.Append(r.AsElement())

		level++
	}
	for i := range serie.Sub {
		if serie.Sub[i].Fill.IsZero() {
			serie.Sub[i].Fill = serie.Fill
		}
		c.drawInterval(a, serie.Sub[i], rx, bar, height, level)
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
	ry = ry.extendBy(1.1)
	cs.Append(c.Chart.drawAxis(rx.AxisRange(), ry.AxisRange()))
	for i := range series {
		if series[i].Curver == nil {
			series[i].Curver = LinearCurve()
		}
		elem := series[i].Draw(c.Chart, &series[i], rx, ry)
		area.Append(elem)
	}
	cs.Append(area.AsElement())
	cs.Append(c.drawTitle())
	cs.Append(c.drawLegend())
	return cs.AsElement()
}

type timepoint struct {
	X time.Time
	Y float64
}

type timepair struct {
	Min time.Time
	Max time.Time
}

func (t timepair) AxisRange() AxisOption {
	return WithTimeRange(t.Min, t.Max)
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

func (t timepair) First() float64 {
	return float64(t.Min.Unix())
}

func (t timepair) Last() float64 {
	return float64(t.Max.Unix())
}

func getIntervalDepth(series []Interval) int {
	var d int
	for i := range series {
		x := series[i].Depth()
		if x > d {
			d = x
		}
	}
	return d
}

func getIntervalDomains(series []Interval) (timepair, []string) {
	var (
		p   timepair
		str []string
	)
	for i := range series {
		str = append(str, series[i].Title)
		starts, ends := series[i].Range()
		if i == 0 || starts.Before(p.Min) {
			p.Min = starts
		}
		if i == 0 || ends.After(p.Max) {
			p.Max = ends
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
