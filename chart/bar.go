package chart

import (
	"bufio"
	"fmt"
	"io"

	"github.com/midbel/svg"
)

type BarSerie struct {
	Title string
	Fill  []svg.Fill

	values []valuelabel
	pair
}

func NewBarSerie(title string) BarSerie {
	return BarSerie{
		Title: title,
	}
}

func (s *BarSerie) Add(label string, val float64) {
	s.pair.Min = getLesser(s.pair.Min, val)
	s.pair.Max = getGreater(s.pair.Max, val)
	vl := valuelabel{
		Label: label,
		Value: val,
	}
	s.values = append(s.values, vl)
}

func (s *BarSerie) Sum() float64 {
	var sum float64
	for i := range s.values {
		sum += s.values[i].Value
	}
	return sum
}

func (s *BarSerie) Len() int {
	return len(s.values)
}

type StackedBarSerie struct {
	Title  string
	series []BarSerie
	pair
}

func NewStackedBarSerie(title string) StackedBarSerie {
	return StackedBarSerie{
		Title: title,
	}
}

func (sr *StackedBarSerie) Append(s BarSerie) {
	sum := s.Sum()
	sr.pair.Min = getLesser(sr.pair.Min, sum)
	sr.pair.Max = getGreater(sr.pair.Max, sum)
	sr.series = append(sr.series, s)
}

func (sr *StackedBarSerie) Len() int {
	return len(sr.series)
}

type BarChart struct {
	Chart
	CategoryAxis
}

func (c BarChart) Render(w io.Writer, series []BarSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	cs := c.RenderElement(series)
	cs.Render(ws)
}

func (c BarChart) RenderElement(series []BarSerie) svg.Element {
	c.checkDefault()
	var (
		cs     = c.getCanvas()
		area   = c.getArea()
		rg, ds = getBarDomains(series)
		offset = c.GetAreaWidth() / float64(len(series))
	)
	rg.Max *= 1.05
	cs.Append(c.drawAxis(c.Chart, rg, ds))
	for i := range series {
		var (
			width = offset * 0.8
			off   = offset*float64(i) + (offset / 2) - (width / 2)
			grp   = svg.NewGroup(svg.WithTranslate(off, 0))
		)
		c.drawSerie(&grp, series[i], rg, width)
		area.Append(grp.AsElement())
	}
	cs.Append(area.AsElement())
	return cs.AsElement()
}

func (c BarChart) drawSerie(a appender, serie BarSerie, rg pair, width float64) {
	var (
		step = width / float64(serie.Len())
		bar  = step * 0.7
		dy   = c.GetAreaHeight() / rg.Diff()
	)
	for i := range serie.values {
		var (
			x = (step * float64(i)) + (step / 2) - (bar / 2)
			h = serie.values[i].Value * dy
			p = svg.NewPos(x, c.GetAreaHeight()-h)
			d = svg.NewDim(bar, h)
			f = serie.Fill[i%len(serie.Fill)]
			r = svg.NewRect(p.Option(), d.Option(), f.Option())
		)
		r.Title = serie.values[i].Label
		a.Append(r.AsElement())
	}
}

type StackedBarChart struct {
	Chart
	CategoryAxis
}

func (c StackedBarChart) Render(w io.Writer, series []StackedBarSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	cs := c.RenderElement(series)
	cs.Render(ws)
}

func (c StackedBarChart) RenderElement(series []StackedBarSerie) svg.Element {
	c.checkDefault()

	var (
		cs     = c.getCanvas()
		area   = c.getArea()
		rg, ds = getStackedBarDomains(series)
		offset = c.GetAreaWidth() / float64(len(series))
	)
	rg.Max *= 1.05
	cs.Append(c.drawAxis(c.Chart, rg, ds))
	for i := range series {
		var (
			width = offset * 0.8
			off   = (offset * float64(i)) + (offset / 2) - (width / 2)
			grp   = svg.NewGroup(svg.WithTranslate(off, 0))
		)
		area.Append(grp.AsElement())
		c.drawSerie(&grp, series[i], rg, width)
	}
	cs.Append(area.AsElement())
	return cs.AsElement()
}

func (c StackedBarChart) drawSerie(a appender, serie StackedBarSerie, rg pair, width float64) {
	var (
		step = width / float64(serie.Len())
		bar  = step * 0.7
		dy   = c.GetAreaHeight() / rg.Diff()
	)
	for j, s := range serie.series {
		var (
			off = c.GetAreaHeight()
			grp = svg.NewGroup(svg.WithTranslate(float64(j)*step, 0))
		)
		for i := range s.values {
			var (
				h = s.values[i].Value * dy
				p = svg.NewPos(0, off-h)
				d = svg.NewDim(bar, h)
				f = s.Fill[i%len(s.Fill)]
				r = svg.NewRect(p.Option(), d.Option(), f.Option())
			)
			r.Title = fmt.Sprintf("%s/%s", s.Title, s.values[i].Label)
			grp.Append(r.AsElement())
			off -= h
		}
		a.Append(grp.AsElement())
	}
}

type valuelabel struct {
	Label string
	Value float64
}

func getBarDomains(series []BarSerie) (pair, []string) {
	var (
		ds []string
		rg pair
	)
	for i := range series {
		ds = append(ds, series[i].Title)
		rg.Min = getLesser(rg.Min, series[i].pair.Min)
		rg.Max = getGreater(rg.Max, series[i].pair.Max)
	}
	return rg, ds
}

func getStackedBarDomains(series []StackedBarSerie) (pair, []string) {
	var (
		ds []string
		rg pair
	)
	for i := range series {
		ds = append(ds, series[i].Title)
		rg.Min = getLesser(rg.Min, series[i].pair.Min)
		rg.Max = getGreater(rg.Max, series[i].pair.Max)
	}
	return rg, ds
}
