package chart

import (
	"bufio"
	"fmt"
	"io"
	"math"

	"github.com/midbel/svg"
)

const (
	defaultStretch = 0.5
	defaultRadius  = 5
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

func (p pair) AbsMax() float64 {
	return math.Max(math.Abs(p.Min), math.Abs(p.Max))
}

func (p pair) AbsMin() float64 {
	return math.Min(math.Abs(p.Min), math.Abs(p.Max))
}

func (p pair) extend() pair {
	x := p
	x.Min *= 1.1
	x.Max *= 1.1
	return x
}

type LineSerie struct {
	Title  string
	Color  string
	values []valuepoint
	min    valuepoint
	max    valuepoint
}

func NewLineSerie(title string) LineSerie {
	return NewLineSerieWithColor(title, "steelblue")
}

func NewLineSerieWithColor(title, color string) LineSerie {
	return LineSerie{
		Title: title,
		Color: color,
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

func (ir *LineSerie) Len() int {
	return len(ir.values)
}

type CurveStyle int8

const (
	CurveLinear CurveStyle = iota
	CurveStep
	CurveStepBefore
	CurveStepAfter
	CurveCubic
	CurveQuadratic
)

type ScatterChart struct {
	Chart
	LineAxis
	Radius float64
}

func (c ScatterChart) Render(w io.Writer, serie LineSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	cs := c.RenderElement(serie)
	cs.Render(ws)
}

func (c ScatterChart) RenderElement(serie LineSerie) svg.Element {
	c.checkDefault()

	var (
		dim    = svg.NewDim(c.Width, c.Height)
		cs     = svg.NewSVG(dim.Option())
		area   = svg.NewGroup(svg.WithID("area"), c.translate())
		rx, ry = getLineDomains([]LineSerie{serie})
		dx     = c.GetAreaWidth() / rx.Diff()
		dy     = c.GetAreaHeight() / ry.Diff()
		fill   = svg.NewFill("steelblue")
		pos    svg.Pos
	)
	cs.Append(c.drawAxisX(c.Chart, rx))
	cs.Append(c.drawAxisY(c.Chart, ry))
	area.Append(c.drawTicksY(c.Chart, ry))
	for i := 0; i < serie.Len(); i++ {
		pos.X = serie.values[i].X * dx
		if rx.Min < 0 {
			pos.X += math.Abs(rx.Min) * dx
		}
		pos.Y = c.GetAreaHeight() - (serie.values[i].Y * dy)
		if ry.Min < 0 {
			pos.Y -= math.Abs(ry.Min) * dy
		}
		ci := svg.NewCircle(pos.Option(), fill.Option(), svg.WithRadius(c.Radius))
		area.Append(ci.AsElement())
	}
	cs.Append(area.AsElement())
	return cs.AsElement()
}

func (c *ScatterChart) checkDefault() {
	c.Chart.checkDefault()
	if c.Radius <= 0 {
		c.Radius = defaultRadius
	}
}

type LineChart struct {
	Chart
	LineAxis
	Curve    CurveStyle
	StretchX float64
	StretchY float64
}

func (c LineChart) Render(w io.Writer, series []LineSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()

	cs := c.RenderElement(series)
	cs.Render(ws)
}

func (c LineChart) RenderElement(series []LineSerie) svg.Element {
	c.checkDefault()

	var (
		dim    = svg.NewDim(c.Width, c.Height)
		cs     = svg.NewSVG(dim.Option())
		area   = svg.NewGroup(svg.WithID("area"), c.translate())
		rx, ry = getLineDomains(series)
	)
	cs.Append(c.drawAxisX(c.Chart, rx))
	cs.Append(c.drawAxisY(c.Chart, ry))
	area.Append(c.drawTicksY(c.Chart, ry))
	for i := range series {
		var elem svg.Element
		switch c.Curve {
		case CurveLinear:
			elem = c.drawLinearSerie(series[i], rx, ry)
		case CurveStep:
			elem = c.drawStepSerie(series[i], rx, ry)
		case CurveStepBefore:
			elem = c.drawStepBeforeSerie(series[i], rx, ry)
		case CurveStepAfter:
			elem = c.drawStepAfterSerie(series[i], rx, ry)
		case CurveCubic:
			elem = c.drawCubicSerie(series[i], rx, ry)
		case CurveQuadratic:
			elem = c.drawQuadraticSerie(series[i], rx, ry)
		default:
		}
		if elem == nil {
			continue
		}
		area.Append(elem)
	}
	cs.Append(area.AsElement())
	return cs.AsElement()
}

func (c LineChart) drawQuadraticSerie(s LineSerie, px, py pair) svg.Element {
	var (
		wx   = c.GetAreaWidth() / px.Diff()
		wy   = c.GetAreaHeight() / py.Diff()
		pat  = getPathLine(s.Color)
		pos  svg.Pos
		old  svg.Pos
		ctrl svg.Pos
	)
	pos.Y = c.GetAreaHeight() - (s.values[0].Y * wy)
	if py.Min < 0 {
		pos.Y -= math.Abs(py.Min) * wy
	}
	pat.AbsMoveTo(pos)
	for i := 1; i < s.Len(); i++ {
		old = pos
		pos.X += (s.values[i].X - s.values[i-1].X) * wx
		pos.Y = c.GetAreaHeight() - (s.values[i].Y * wy)
		if py.Min < 0 {
			pos.Y -= math.Abs(py.Min) * wy
		}
		ctrl.X = old.X
		ctrl.Y = pos.Y
		pat.AbsQuadraticCurve(pos, ctrl)
	}
	return pat.AsElement()
}

func (c LineChart) drawCubicSerie(s LineSerie, px, py pair) svg.Element {
	var (
		wx  = c.GetAreaWidth() / px.Diff()
		wy  = c.GetAreaHeight() / py.Diff()
		pat = getPathLine(s.Color)
		pos svg.Pos
	)
	pos.Y = c.GetAreaHeight() - (s.values[0].Y * wy)
	if py.Min < 0 {
		pos.Y -= math.Abs(py.Min) * wy
	}
	pat.AbsMoveTo(pos)
	for i := 1; i < s.Len(); i++ {
		var (
			ctrl = pos
			old  = pos
		)
		pos.Y = c.GetAreaHeight() - (s.values[i].Y * wy)
		if py.Min < 0 {
			pos.Y -= math.Abs(py.Min) * wy
		}
		pos.X += (s.values[i].X - s.values[i-1].X) * wx
		ctrl.X = old.X - (old.X-pos.X)*c.StretchX
		ctrl.Y = pos.Y
		pat.AbsCubicCurveSimple(pos, ctrl)
	}
	return pat.AsElement()
}

func (c LineChart) drawStepSerie(s LineSerie, px, py pair) svg.Element {
	var (
		wx  = c.GetAreaWidth() / px.Diff()
		wy  = c.GetAreaHeight() / py.Diff()
		pat = getPathLine(s.Color)
		y   = c.GetAreaHeight() - (s.values[0].Y * wy)
		x   float64
	)
	if py.Min < 0 {
		y -= math.Abs(py.Min) * wy
	}
	pat.AbsMoveTo(svg.NewPos(x, y))

	for i := 1; i < s.Len(); i++ {
		delta := (s.values[i].X - s.values[i-1].X) / 2
		x += delta * wx
		pat.AbsLineTo(svg.NewPos(x, y))

		y = c.GetAreaHeight() - (s.values[i].Y * wy)
		if py.Min < 0 {
			y -= math.Abs(py.Min) * wy
		}
		pat.AbsLineTo(svg.NewPos(x, y))
		x += delta * wx
		pat.AbsLineTo(svg.NewPos(x, y))
	}
	return pat.AsElement()
}

func (c LineChart) drawStepBeforeSerie(s LineSerie, px, py pair) svg.Element {
	var (
		wx  = c.GetAreaWidth() / px.Diff()
		wy  = c.GetAreaHeight() / py.Diff()
		pat = getPathLine(s.Color)
		y   = c.GetAreaHeight() - (s.values[0].Y * wy)
		x   float64
	)
	if py.Min < 0 {
		y -= math.Abs(py.Min) * wy
	}
	pat.AbsMoveTo(svg.NewPos(x, y))
	for i := 1; i < s.Len(); i++ {
		y = c.GetAreaHeight() - (s.values[i].Y * wy)
		if py.Min < 0 {
			y -= math.Abs(py.Min) * wy
		}
		pat.AbsLineTo(svg.NewPos(x, y))
		x += (s.values[i].X - s.values[i-1].X) * wx
		pat.AbsLineTo(svg.NewPos(x, y))
	}
	return pat.AsElement()
}

func (c LineChart) drawStepAfterSerie(s LineSerie, px, py pair) svg.Element {
	var (
		wx  = c.GetAreaWidth() / px.Diff()
		wy  = c.GetAreaHeight() / py.Diff()
		pat = getPathLine(s.Color)
		y   = c.GetAreaHeight() - (s.values[0].Y * wy)
		x   float64
	)
	if py.Min < 0 {
		y -= math.Abs(py.Min) * wy
	}
	pat.AbsMoveTo(svg.NewPos(x, y))
	for i := 1; i < s.Len(); i++ {
		x += (s.values[i].X - s.values[i-1].X) * wx
		pat.AbsLineTo(svg.NewPos(x, y))

		y = c.GetAreaHeight() - (s.values[i].Y * wy)
		if py.Min < 0 {
			y -= math.Abs(py.Min) * wy
		}
		pat.AbsLineTo(svg.NewPos(x, y))
	}
	return pat.AsElement()
}

func (c LineChart) drawLinearSerie(s LineSerie, px, py pair) svg.Element {
	var (
		wx  = c.GetAreaWidth() / px.Diff()
		wy  = c.GetAreaHeight() / py.Diff()
		pat = getPathLine(s.Color)
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
		pos.X = (s.values[i].X - s.values[0].X) * wx
		pos.Y = c.GetAreaHeight() - (s.values[i].Y * wy)
		if py.Min < 0 {
			pos.Y -= math.Abs(py.Min) * wy
		}
		pat.AbsLineTo(pos)
	}
	return pat.AsElement()
}

func (c *LineChart) checkDefault() {
	c.Chart.checkDefault()
	if c.Curve == CurveCubic || c.Curve == CurveQuadratic {
		if c.StretchX == 0 {
			c.StretchX = defaultStretch
		}
		if c.StretchY == 0 {
			c.StretchY = defaultStretch
		}
	}
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

type LineAxis struct {
	TicksX int
	TicksY int
}

func (a LineAxis) drawAxisX(c Chart, rg pair) svg.Element {
	var (
		axis  = svg.NewGroup(c.getOptionsAxisX()...)
		pos1  = svg.NewPos(0, 0)
		pos2  = svg.NewPos(c.GetAreaWidth(), 0)
		line  = svg.NewLine(pos1, pos2, axisstrok.Option(), svg.WithClass("domain"))
		coeff = c.GetAreaWidth() / float64(a.TicksX)
		step  = math.Ceil(rg.Diff() / float64(a.TicksX))
	)
	_ = fmt.Sprintf
	for i, j := rg.Min, 0; i < rg.Max+step; i, j = i+step, j+1 {
		var (
			grp  = svg.NewGroup(svg.WithClass("tick"))
			off  = float64(j) * coeff
			pos0 = svg.NewPos(off-(coeff*0.2), textick+(textick/3))
			pos1 = svg.NewPos(off, 0)
			pos2 = svg.NewPos(off, ticklen)
			text = svg.NewText(formatFloat(i), pos0.Option())
			line = svg.NewLine(pos1, pos2, axisstrok.Option())
		)
		grp.Append(text.AsElement())
		grp.Append(line.AsElement())
		axis.Append(grp.AsElement())
	}
	axis.Append(line.AsElement())
	return axis.AsElement()
}

func (a LineAxis) drawAxisY(c Chart, rg pair) svg.Element {
	var (
		axis  = svg.NewGroup(c.getOptionsAxisY()...)
		pos1  = svg.NewPos(0, 0)
		pos2  = svg.NewPos(0, c.GetAreaHeight()+1)
		line  = svg.NewLine(pos1, pos2, axisstrok.Option(), svg.WithClass("domain"))
		coeff = c.GetAreaHeight() / float64(a.TicksY)
		step  = rg.Diff() / float64(a.TicksY)
	)
	axis.Append(line.AsElement())
	for i, j := rg.Min, 0; i < rg.Max+step; i, j = i+step, j+1 {
		var (
			grp  = svg.NewGroup(svg.WithClass("tick"))
			ypos = c.GetAreaHeight() - (float64(j) * coeff)
			pos  = svg.NewPos(0, ypos+(ticklen/2))
			anc  = svg.WithAnchor("end")
			text = svg.NewText(formatFloat(i), anc, pos.Option())
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

func (a LineAxis) drawTicksY(c Chart, rg pair) svg.Element {
	var (
		max   = rg.AbsMax()
		grp   = svg.NewGroup(svg.WithID("ticks"))
		step  = c.GetAreaHeight() / max
		coeff = max / float64(a.TicksY)
	)
	for i := a.TicksY; i > 0; i-- {
		var (
			ypos = c.GetAreaHeight() - (float64(i) * coeff * step)
			pos1 = svg.NewPos(0, ypos)
			pos2 = svg.NewPos(c.GetAreaWidth(), ypos)
		)
		grp.Append(getTick(pos1, pos2))
	}
	return grp.AsElement()
}
