package chart

import (
	"bufio"
	"io"
	"math"

	"github.com/midbel/svg"
	"github.com/midbel/svg/colors"
)

type valuepoint struct {
	X float64
	Y float64
}

type pair struct {
	Min float64
	Max float64
}

func (p pair) Diff() float64 {
	return p.Max - p.Min
}

type LineSerie struct {
	Title  string
	values []valuepoint
	min    valuepoint
	max    valuepoint
	colors []string
}

func NewLineSerie(title string) LineSerie {
	return NewLineSerieWithColors(title, colors.Set26)
}

func NewLineSerieWithColors(title string, colors []string) LineSerie {
	return LineSerie{
		Title: title,
		min: valuepoint{
			X: math.NaN(),
			Y: math.NaN(),
		},
		max: valuepoint{
			X: math.NaN(),
			Y: math.NaN(),
		},
		colors: colors,
	}
}

func (ir *LineSerie) Add(x, y float64) {
	ir.min.X = getLesser(ir.min.X, x)
	ir.min.Y = getLesser(ir.min.Y, y)
	ir.max.X = getGreater(ir.max.X, x)
	ir.max.Y = getGreater(ir.max.Y, y)
	vp := valuepoint{
		X: x,
		Y: y,
	}
	ir.values = append(ir.values, vp)
}

type LineChart struct {
	Chart
	TicksY int
	TicksX int
}

func (c LineChart) Render(w io.Writer, series []LineSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	c.render(ws, series)
}

func (c LineChart) render(w svg.Writer, series []LineSerie) {
	c.checkDefault()

	var (
		dim    = svg.NewDim(c.Width, c.Height)
		cs     = svg.NewSVG(dim.Option())
		area   = svg.NewGroup(svg.WithID("area"), c.translate())
		rx, ry = getLineDomains(series)
	)
	if c.TicksY > 0 {
		area.Append(c.drawTicks(ry))
	}
	for i := range series {
		grp := svg.NewGroup()
		grp.Append(c.drawSerie(series[i], rx, ry))
		area.Append(grp.AsElement())
	}
	cs.Append(area.AsElement())
	if c.TicksX > 0 {
		cs.Append(c.drawAxisX(rx))
	}
	if c.TicksY > 0 {
		cs.Append(c.drawAxisY(ry))
	}
	cs.Render(w)
}

func (c LineChart) drawSerie(s LineSerie, px, py pair) svg.Element {
	var (
		grp = svg.NewGroup(svg.WithClass("line"))
		dx  = px.Diff()
		dy  = py.Diff()
		pat = getPathLine("steelblue")
	)
	for i := range s.values {
		var (
			x = (s.values[i].X / dx) * c.GetAreaWidth()
			y = c.GetAreaHeight() - (s.values[i].Y/dy)*c.GetAreaHeight()
			p = svg.NewPos(x, y)
		)
		if i == 0 {
			pat.AbsMoveTo(p)
			continue
		}
		pat.AbsLineTo(p)
	}
	grp.Append(pat.AsElement())
	return grp.AsElement()
}

func (c LineChart) drawAxisX(rg pair) svg.Element {
	var (
		axis  = svg.NewGroup(c.getOptionsAxisX()...)
		pos1  = svg.NewPos(0, 0)
		pos2  = svg.NewPos(c.GetAreaWidth(), 0)
		line  = svg.NewLine(pos1, pos2, axisstrok.Option(), svg.WithClass("domain"))
		coeff = c.GetAreaWidth() / float64(c.TicksX)
		step  = rg.Max / float64(c.TicksX)
	)
	for i := 0; i < c.TicksX+1; i++ {
		var (
			grp  = svg.NewGroup(svg.WithClass("tick"))
			off  = float64(i) * coeff
			pos0 = svg.NewPos(off-(step*1.5), textick+(textick/3))
			pos1 = svg.NewPos(off, 0)
			pos2 = svg.NewPos(off, ticklen)
			text = svg.NewText(formatFloat(float64(i)*step), pos0.Option())
			line = svg.NewLine(pos1, pos2, axisstrok.Option())
		)
		grp.Append(text.AsElement())
		grp.Append(line.AsElement())
		axis.Append(grp.AsElement())
	}
	axis.Append(line.AsElement())
	return axis.AsElement()
}

func (c LineChart) drawAxisY(rg pair) svg.Element {
	var (
		axis  = svg.NewGroup(c.getOptionsAxisY()...)
		pos1  = svg.NewPos(0, 0)
		pos2  = svg.NewPos(0, c.GetAreaHeight()+1)
		line  = svg.NewLine(pos1, pos2, axisstrok.Option(), svg.WithClass("domain"))
		coeff = c.GetAreaHeight() / float64(c.TicksY)
		step  = rg.Max / float64(c.TicksY)
	)
	axis.Append(line.AsElement())
	for i := c.TicksY; i >= 0; i-- {
		var (
			grp  = svg.NewGroup(svg.WithClass("tick"))
			ypos = c.GetAreaHeight() - (float64(i) * coeff)
			val  = step * float64(i)
			pos  = svg.NewPos(0, ypos+(ticklen/2))
			anc  = svg.WithAnchor("end")
			text = svg.NewText(formatFloat(val), anc, pos.Option())
			pos1 = svg.NewPos(-ticklen, ypos)
			pos2 = svg.NewPos(0, ypos)
			line = svg.NewLine(pos1, pos2, axisstrok.Option())
		)
		text.Shift = svg.NewPos(-ticklen*2, 0)
		grp.Append(text.AsElement())
		grp.Append(line.AsElement())
		axis.Append(grp.AsElement())
	}
	return axis.AsElement()
}

func (c LineChart) drawTicks(rg pair) svg.Element {
	var (
		grp   = svg.NewGroup(svg.WithID("ticks"))
		step  = c.GetAreaHeight() / rg.Max
		coeff = rg.Max / float64(c.TicksY)
	)
	for i := c.TicksY; i > 0; i-- {
		var (
			ypos = c.GetAreaHeight() - (float64(i) * coeff * step)
			pos1 = svg.NewPos(0, ypos)
			pos2 = svg.NewPos(c.GetAreaWidth(), ypos)
		)
		grp.Append(getTick(pos1, pos2))
	}
	return grp.AsElement()
}

func getLineDomains(series []LineSerie) (pair, pair) {
	var x, y pair
	for i := range series {
		x.Min = getLesser(series[i].min.X, x.Min)
		y.Min = getLesser(series[i].min.Y, y.Min)
		x.Max = getGreater(series[i].max.X, x.Max)
		y.Max = getGreater(series[i].max.Y, y.Max)
	}
	return x, y
}
