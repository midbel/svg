package chart

import (
	"bufio"
	"fmt"
	"io"

	"github.com/midbel/svg"
	"github.com/midbel/svg/colors"
)

type valuelabel struct {
	Label string
	Value float64
}

type Serie struct {
	Title  string
	values []valuelabel
	colors []string
}

func NewSerie(title string) Serie {
	return NewSerieWithColors(title, colors.Reverse(colors.PuBu6))
}

func NewSerieWithColors(title string, colors []string) Serie {
	return Serie{
		Title:  title,
		colors: colors,
	}
}

func (s *Serie) Add(label string, value float64) {
	vl := valuelabel{
		Label: label,
		Value: value,
	}
	s.values = append(s.values, vl)
}

func (s Serie) Sum() float64 {
	var sum float64
	for i := range s.values {
		sum += s.values[i].Value
	}
	return sum
}

func (s Serie) Len() int {
	return len(s.values)
}

func (s Serie) peekFill(i int) svg.Option {
	color := s.colors[i%len(s.colors)]
	return svg.NewFill(color).Option()
}

type StackedSerie struct {
	Title  string
	Series []Serie
	max    float64
	min    float64
}

func NewStackedSerie(title string) StackedSerie {
	return StackedSerie{
		Title: title,
	}
}

func (sr *StackedSerie) Append(s Serie) {
	sum := s.Sum()
	sr.min = getLesser(sr.min, sum)
	sr.max = getGreater(sr.max, sum)
	sr.Series = append(sr.Series, s)
}

func (sr *StackedSerie) Len() int {
	return len(sr.Series)
}

type BarChart struct {
	Chart
	CategoryAxis

	BarWidth float64
	Ticks    int
}

func (c BarChart) Render(w io.Writer, serie []Serie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	cs := c.RenderElement(serie)
	cs.Render(ws)
}

func (c BarChart) RenderElement(serie []Serie) svg.Element {
	c.checkDefault()
	var (
		dim  = svg.NewDim(c.Width, c.Height)
		cs   = svg.NewSVG(dim.Option())
		area = svg.NewGroup(svg.WithID("area"), c.translate())
	)
	cs.Append(area.AsElement())
	return cs.AsElement()
}

type StackedChart struct {
	Chart
	CategoryAxis

	BarWidth float64
	Ticks    int
}

func (c StackedChart) Render(w io.Writer, series []StackedSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	cs := c.RenderElement(series)
	cs.Render(ws)
}

func (c StackedChart) RenderElement(series []StackedSerie) svg.Element {
	c.checkDefault()

	var (
		dim    = svg.NewDim(c.Width, c.Height)
		cs     = svg.NewSVG(dim.Option())
		area   = svg.NewGroup(svg.WithID("area"), c.translate())
		rg, ds = getStackedDomains(series)
		offset = c.GetAreaWidth() / float64(len(series))
	)

	cs.Append(c.drawAxis(c.Chart, rg, ds))
	for i := range series {
		var (
			off  = offset * float64(i)
			grp  = svg.NewGroup(svg.WithTranslate(off, 0))
			elem = c.drawSerie(series[i], offset, rg.Max)
		)
		grp.Append(elem)
		area.Append(grp.AsElement())
	}
	cs.Append(area.AsElement())
	return cs.AsElement()
}

func (c StackedChart) drawSerie(s StackedSerie, band, max float64) svg.Element {
	var (
		size   = band / float64(s.Len())
		width  = size * 0.8
		height = c.GetAreaHeight() / max
		trst   svg.Option
	)
	if c.BarWidth > 0 {
		width = c.BarWidth
		off := band - (width * float64(s.Len()))
		trst = svg.WithTranslate(off/2, 0)
	} else {
		trst = svg.WithTranslate((size*0.2)*2, 0)
	}
	grp := svg.NewGroup(trst, svg.WithClass("bar"))
	for j := range s.Series {
		var (
			g  = svg.NewGroup()
			rw = width
			ro = c.GetAreaHeight()
		)
		for i, v := range s.Series[j].values {
			if v.Value == 0 {
				continue
			}
			var (
				rh   = v.Value * height
				rx   = (float64(j) * width) + ((width / 2) - (rw / 2))
				ry   float64
				fill = s.Series[j].peekFill(i)
			)
			ro -= rh
			ry = ro

			r := getRect(svg.WithPosition(rx, ry), svg.WithDimension(rw, rh), fill)
			r.Title = fmt.Sprintf("%s/%s", s.Title, v.Label)
			g.Append(r.AsElement())
		}
		grp.Append(g.AsElement())
	}
	return grp.AsElement()
}

func getStackedDomains(cs []StackedSerie) (pair, []string) {
	var (
		ys []string
		xs pair
	)
	for i := range cs {
		ys = append(ys, cs[i].Title)
		xs.Max = getGreater(xs.Max, cs[i].max)
		xs.Min = getLesser(xs.Min, cs[i].min)
	}
	return xs, ys
}
