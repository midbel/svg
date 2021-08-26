package chart

import (
	"bufio"
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

type LineSerie struct {
	Title string
	Color string
	Shape ShapeType

	values []valuepoint
	min    valuepoint
	max    valuepoint
}

func NewLineSerie(title string) LineSerie {
	return LineSerie{
		Title: title,
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

type CurveStyle uint8

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
	Shape  ShapeType
	Radius float64
}

func (c ScatterChart) Render(w io.Writer, series []LineSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	cs := c.RenderElement(series)
	cs.Render(ws)
}

func (c ScatterChart) RenderElement(series []LineSerie) svg.Element {
	c.checkDefault()
	var (
		dim    = svg.NewDim(c.Width, c.Height)
		cs     = svg.NewSVG(dim.Option())
		area   = svg.NewGroup(svg.WithID("area"), c.translate())
		rx, ry = getLineDomains(series, 1.15)
	)
	cs.Append(c.drawAxis(c.Chart, rx, ry))
	for i := range series {
		grp := c.drawSerie(series[i], c.GetColor(series[i].Title, i), rx, ry)
		area.Append(grp)
	}

	cs.Append(area.AsElement())
	return cs.AsElement()
}

func (c ScatterChart) drawSerie(serie LineSerie, fill svg.Fill, rx, ry pair) svg.Element {
	var (
		grp   = svg.NewGroup(fill.Option())
		dx    = c.GetAreaWidth() / rx.Diff()
		dy    = c.GetAreaHeight() / ry.Diff()
		strok = svg.NewStroke("none", 0).Option()
		pos   svg.Pos
	)
	if serie.Shape == ShapeDefault {
		serie.Shape = ShapeCircle
	}
	for i := 0; i < serie.Len(); i++ {
		pos.X = serie.values[i].X * dx
		if rx.Min < 0 {
			pos.X += math.Abs(rx.Min) * dx
		}
		pos.Y = c.GetAreaHeight() - (serie.values[i].Y * dy)
		if ry.Min < 0 {
			pos.Y -= math.Abs(ry.Min) * dy
		}
		var elem svg.Element
		switch xy := pos.Option(); serie.Shape {
		case ShapeDefault, ShapeCircle:
			elem = serie.Shape.Draw(c.Radius, xy, strok, fill.Option())
		case ShapeTriangle, ShapeStar:
			pos.X -= c.Radius / 2
			pos.Y -= c.Radius / 2
			g := svg.NewGroup(svg.WithTranslate(pos.X, pos.Y))
			i := serie.Shape.Draw(c.Radius, strok, fill.Option())
			g.Append(i)
			elem = g.AsElement()
		case ShapeDiamond:
			pos.X -= c.Radius / 2
			pos.Y -= c.Radius / 2
			rot := svg.WithRotate(45, pos.X, pos.Y)
			elem = serie.Shape.Draw(c.Radius, xy, fill.Option(), strok, rot)
		case ShapeSquare:
			pos.X -= c.Radius / 2
			pos.Y -= c.Radius / 2
			elem = serie.Shape.Draw(c.Radius, xy, fill.Option(), strok)
		default:
		}
		if elem == nil {
			continue
		}
		grp.Append(elem)
	}
	return grp.AsElement()
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
		rx, ry = getLineDomains(series, 0)
	)
	ry = ry.extendBy(1.2)
	cs.Append(c.drawAxis(c.Chart, rx, ry))
	for i := range series {
		var elem svg.Element
		switch c.Curve {
		case CurveLinear:
			elem = c.drawLinearSerie(series[i], c.GetStroke(series[i].Title, i), rx, ry)
		case CurveStep:
			elem = c.drawStepSerie(series[i], c.GetStroke(series[i].Title, i), rx, ry)
		case CurveStepBefore:
			elem = c.drawStepBeforeSerie(series[i], c.GetStroke(series[i].Title, i), rx, ry)
		case CurveStepAfter:
			elem = c.drawStepAfterSerie(series[i], c.GetStroke(series[i].Title, i), rx, ry)
		case CurveCubic:
			elem = c.drawCubicSerie(series[i], c.GetStroke(series[i].Title, i), rx, ry)
		case CurveQuadratic:
			elem = c.drawQuadraticSerie(series[i], c.GetStroke(series[i].Title, i), rx, ry)
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

func (c LineChart) drawQuadraticSerie(s LineSerie, strok svg.Stroke, px, py pair) svg.Element {
	var (
		wx   = c.GetAreaWidth() / px.Diff()
		wy   = c.GetAreaHeight() / py.Diff()
		pat  = svg.NewPath(strok.Option(), nonefill.Option())
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

func (c LineChart) drawCubicSerie(s LineSerie, strok svg.Stroke, px, py pair) svg.Element {
	var (
		wx  = c.GetAreaWidth() / px.Diff()
		wy  = c.GetAreaHeight() / py.Diff()
		pat  = svg.NewPath(strok.Option(), nonefill.Option())
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

func (c LineChart) drawStepSerie(s LineSerie, strok svg.Stroke, px, py pair) svg.Element {
	var (
		wx  = c.GetAreaWidth() / px.Diff()
		wy  = c.GetAreaHeight() / py.Diff()
		pat  = svg.NewPath(strok.Option(), nonefill.Option())
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

func (c LineChart) drawStepBeforeSerie(s LineSerie, strok svg.Stroke, px, py pair) svg.Element {
	var (
		wx  = c.GetAreaWidth() / px.Diff()
		wy  = c.GetAreaHeight() / py.Diff()
		pat  = svg.NewPath(strok.Option(), nonefill.Option())
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

func (c LineChart) drawStepAfterSerie(s LineSerie, strok svg.Stroke, px, py pair) svg.Element {
	var (
		wx  = c.GetAreaWidth() / px.Diff()
		wy  = c.GetAreaHeight() / py.Diff()
		pat  = svg.NewPath(strok.Option(), nonefill.Option())
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

func (c LineChart) drawLinearSerie(s LineSerie, strok svg.Stroke, px, py pair) svg.Element {
	var (
		wx  = c.GetAreaWidth() / px.Diff()
		wy  = c.GetAreaHeight() / py.Diff()
		pat  = svg.NewPath(strok.Option(), nonefill.Option())
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

func (p pair) extendBy(by float64) pair {
	p.Min *= by
	p.Max *= by
	return p
}

func getLineDomains(series []LineSerie, mul float64) (pair, pair) {
	var x, y pair
	for i := range series {
		x.Min = getLesser(series[i].min.X, x.Min)
		y.Min = getLesser(series[i].min.Y, y.Min)
		x.Max = getGreater(series[i].max.X, x.Max)
		y.Max = getGreater(series[i].max.Y, y.Max)
	}
	if mul <= 1 {
		return x, y
	}
	return x.extendBy(mul), y.extendBy(mul)
}
