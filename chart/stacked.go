package chart

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"

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

func (s Serie) Labels() []string {
	var str []string
	for i := range s.values {
		str = append(str, s.values[i].Label)
	}
	return str
}

func (s Serie) Values() []float64 {
	var str []float64
	for i := range s.values {
		str = append(str, s.values[i].Value)
	}
	return str
}

func (s Serie) peekFill(i int) svg.Option {
	color := s.colors[i%len(s.colors)]
	return svg.NewFill(color).Option()
}

type StackedSerie struct {
	Title  string
	Series []Serie
	max    float64
}

func NewStackedSerie(title string) StackedSerie {
	return StackedSerie{
		Title: title,
		max:   math.NaN(),
	}
}

func (sr *StackedSerie) Append(s Serie) {
	if sum := s.Sum(); math.IsNaN(sr.max) || sum > sr.max {
		sr.max = sum
	}
	sr.Series = append(sr.Series, s)
}

func (sr *StackedSerie) Len() int {
	return len(sr.Series)
}

type StackedChart struct {
	Chart

	BarWidth float64
	Ticks    int
}

func (c StackedChart) Render(w io.Writer, series []StackedSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	c.render(ws, series)
}

func (c StackedChart) render(w svg.Writer, series []StackedSerie) {
	c.checkDefault()

	var (
		dim     = svg.NewDim(c.Width, c.Height)
		cs      = svg.NewSVG(dim.Option())
		area    = svg.NewGroup(svg.WithID("area"), c.translate())
		max, ds = getStackedDomains(series)
	)
	max *= 1.01
	if c.Ticks > 0 {
		area.Append(c.drawTicks(max))
	}

	offset := c.GetAreaWidth() / float64(len(series))
	for i := range series {
		var (
			off  = offset * float64(i)
			grp  = svg.NewGroup(svg.WithTranslate(off, 0))
			elem = c.drawSerie(series[i], offset, max)
		)
		grp.Append(elem)
		area.Append(grp.AsElement())
	}
	cs.Append(area.AsElement())
	cs.Append(c.drawAxisX(ds))
	if c.Ticks > 0 {
		cs.Append(c.drawAxisY(max))
	}
	cs.Render(w)
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
	grp := svg.NewGroup(trst)
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

func (c StackedChart) drawAxisX(domains []string) svg.Element {
	options := []svg.Option{
		svg.WithID("x-axis"),
		svg.WithClass("axis"),
		svg.WithTranslate(c.Padding.Left, c.Height-c.Padding.Bottom),
	}
	var (
		axis = svg.NewGroup(options...)
		pos1 = svg.NewPos(0, 0)
		pos2 = svg.NewPos(c.GetAreaWidth(), 0)
		line = svg.NewLine(pos1, pos2, axisstrok.Option(), svg.WithClass("domain"))
		step = c.GetAreaWidth() / float64(len(domains))
	)
	axis.Append(line.AsElement())
	for i := 0; i < len(domains); i++ {
		var (
			grp  = svg.NewGroup(svg.WithClass("tick"))
			off  = float64(i) * step
			pos0 = svg.NewPos(off+(step/3), textick)
			pos1 = svg.NewPos(off+step, 0)
			pos2 = svg.NewPos(off+step, ticklen)
			text = svg.NewText(domains[i], pos0.Option())
			line = svg.NewLine(pos1, pos2, axisstrok.Option())
		)
		grp.Append(text.AsElement())
		grp.Append(line.AsElement())
		axis.Append(grp.AsElement())
	}
	return axis.AsElement()
}

func (c StackedChart) drawAxisY(max float64) svg.Element {
	options := []svg.Option{
		svg.WithID("y-axis"),
		svg.WithClass("axis"),
		c.translate(),
	}
	var (
		axis  = svg.NewGroup(options...)
		pos1  = svg.NewPos(0, 0)
		pos2  = svg.NewPos(0, c.GetAreaHeight()+1)
		line  = svg.NewLine(pos1, pos2, axisstrok.Option(), svg.WithClass("domain"))
		step  = c.GetAreaHeight() / max
		coeff = max / float64(c.Ticks)
	)
	axis.Append(line.AsElement())
	for i := c.Ticks; i >= 0; i-- {
		var (
			grp   = svg.NewGroup(svg.WithClass("tick"))
			val   = coeff * float64(i)
			pos   = svg.NewPos(0, c.GetAreaHeight()-(step*val)+(ticklen/2))
			anc   = svg.WithAnchor("end")
			label = strconv.FormatFloat(val, 'f', 2, 64)
			text  = svg.NewText(label, anc, pos.Option())
			pos1  = svg.NewPos(-ticklen, c.GetAreaHeight()-(step*val))
			pos2  = svg.NewPos(0, c.GetAreaHeight()-(step*val))
			line  = svg.NewLine(pos1, pos2, axisstrok.Option())
		)
		text.Shift = svg.NewPos(-ticklen*2, 0)
		grp.Append(text.AsElement())
		grp.Append(line.AsElement())
		axis.Append(grp.AsElement())
	}
	return axis.AsElement()
}

func (c StackedChart) drawTicks(max float64) svg.Element {
	var (
		grp   = svg.NewGroup(svg.WithID("ticks"))
		step  = c.GetAreaHeight() / max
		coeff = max / float64(c.Ticks)
	)
	for i := c.Ticks; i > 0; i-- {
		var (
			ypos = c.GetAreaHeight() - (float64(i) * coeff * step)
			pos1 = svg.NewPos(0, ypos)
			pos2 = svg.NewPos(c.GetAreaWidth(), ypos)
		)
		grp.Append(getTick(pos1, pos2))
	}
	return grp.AsElement()
}

func getStackedDomains(cs []StackedSerie) (float64, []string) {
	var (
		ys []string
		xs float64
	)
	for i := range cs {
		ys = append(ys, cs[i].Title)
		if i == 0 || cs[i].max > xs {
			xs = cs[i].max
		}
	}
	return xs, ys
}
