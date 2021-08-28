package chart

import (
	"bufio"
	"io"
	"math"

	"github.com/midbel/svg"
)

const (
	defaultStretch = 0.5
	DefaultSize    = 5
)

type CurveStyle uint8

const (
	CurveLinear CurveStyle = iota
	CurveStep
	CurveStepBefore
	CurveStepAfter
	CurveCubic
	CurveQuadratic
)

type LineSerie struct {
	xyserie

	svg.Stroke
	Curve CurveStyle
	Shape ShapeType
}

func NewLineSerie(title string) LineSerie {
	s := xyserie{Title: title}
	return LineSerie{xyserie: s}
}

type ScatterSerie struct {
	xyserie

	svg.Fill
	svg.Stroke
	Size      float64
	Shape     ShapeType
	Highlight bool
}

func NewScatterSerie(title string) ScatterSerie {
	s := xyserie{Title: title}
	return ScatterSerie{xyserie: s, Size: DefaultSize}
}

type ScatterChart struct {
	Chart
	LineAxis
}

func (c ScatterChart) Render(w io.Writer, series []ScatterSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	cs := c.RenderElement(series)
	cs.Render(ws)
}

func (c ScatterChart) RenderElement(series []ScatterSerie) svg.Element {
	c.checkDefault()
	var (
		dim    = svg.NewDim(c.Width, c.Height)
		cs     = svg.NewSVG(dim.Option())
		area   = svg.NewGroup(svg.WithID("area"), c.translate())
		rx, ry = getScatterDomains(series, 1.15)
	)
	cs.Append(c.drawAxis(c.Chart, rx, ry))
	for i := range series {
		grp := c.drawSerie(series[i], rx, ry)
		area.Append(grp)
	}

	cs.Append(area.AsElement())
	return cs.AsElement()
}

func (c ScatterChart) drawSerie(serie ScatterSerie, rx, ry pair) svg.Element {
	var (
		fill = serie.Fill.Option()
		grp  = svg.NewGroup(fill)
		dx   = c.GetAreaWidth() / rx.Diff()
		dy   = c.GetAreaHeight() / ry.Diff()
		pos  svg.Pos
	)
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
		switch xy, radius := pos.Option(), serie.Size; serie.Shape {
		case ShapeCircle:
			elem = serie.Shape.Draw(radius, xy, fill)
		case ShapeTriangle, ShapeStar:
			pos.X -= radius / 2
			pos.Y -= radius / 2
			g := svg.NewGroup(svg.WithTranslate(pos.X, pos.Y))
			i := serie.Shape.Draw(radius, fill)
			g.Append(i)
			elem = g.AsElement()
		case ShapeDiamond:
			pos.X -= radius / 2
			pos.Y -= radius / 2
			rot := svg.WithRotate(45, pos.X, pos.Y)
			elem = serie.Shape.Draw(radius, xy, fill, rot)
		case ShapeSquare, ShapeDefault:
			pos.X -= radius / 2
			pos.Y -= radius / 2
			elem = serie.Shape.Draw(radius, xy, fill)
		default:
		}
		grp.Append(elem)
	}
	if serie.Highlight {
		grp.Append(c.highlightSerie(serie, rx, ry))
	}
	return grp.AsElement()
}

func (c *ScatterChart) highlightSerie(serie ScatterSerie, rx, ry pair) svg.Element {
	var (
		dx = c.GetAreaWidth() / rx.Diff()
		dy = c.GetAreaHeight() / ry.Diff()
		x0 = serie.px.Min * dx
		y0 = c.GetAreaHeight() - (serie.py.Max * dy)
		x1 = serie.px.Max * dx
		y1 = c.GetAreaHeight() - (serie.py.Min * dy)
	)
	if rx.Min < 0 {
		x0 += math.Abs(rx.Min) * dx
		x1 += math.Abs(rx.Min) * dx
	}
	if ry.Min < 0 {
		y1 -= math.Abs(ry.Min) * dy
		y0 -= math.Abs(ry.Min) * dy
	}
	x0 -= serie.Size
	x1 += serie.Size
	y0 -= serie.Size
	y1 += serie.Size * 2
	var (
		pos   = svg.NewPos(x0, y0)
		dim   = svg.NewDim(x1-x0, y1-y0)
		fill  = svg.NewFill("transparent").Option()
		strok = serie.Stroke.Option()
		rect  = svg.NewRect(pos.Option(), dim.Option(), strok, fill)
	)
	return rect.AsElement()
}

type LineChart struct {
	Chart
	LineAxis
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
		switch series[i].Curve {
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
		pat = svg.NewPath(strok.Option(), nonefill.Option())
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
		pat = svg.NewPath(strok.Option(), nonefill.Option())
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
		pat = svg.NewPath(strok.Option(), nonefill.Option())
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
		pat = svg.NewPath(strok.Option(), nonefill.Option())
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
		pat = svg.NewPath(strok.Option(), nonefill.Option())
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

func (p pair) extendBy(by float64) pair {
	p.Min *= by
	p.Max *= by
	return p
}

type xyserie struct {
	Title  string
	values []valuepoint
	px     pair
	py     pair
}

func (xy *xyserie) Add(x, y float64) {
	if len(xy.values) == 0 {
		xy.px.Min = x
		xy.px.Max = x
		xy.py.Min = y
		xy.py.Max = y
	}
	xy.px.Min = getLesser(xy.px.Min, x)
	xy.px.Max = getGreater(xy.px.Max, x)
	xy.py.Min = getLesser(xy.py.Min, y)
	xy.py.Max = getGreater(xy.py.Max, y)
	vp := valuepoint{
		X: x,
		Y: y,
	}
	xy.values = append(xy.values, vp)
}

func (xy *xyserie) Len() int {
	return len(xy.values)
}

func getLineDomains(series []LineSerie, mul float64) (pair, pair) {
	vs := make([]xyserie, len(series))
	for i := range series {
		vs[i] = series[i].xyserie
	}
	return getDomainsXY(vs, mul)
}

func getScatterDomains(series []ScatterSerie, mul float64) (pair, pair) {
	vs := make([]xyserie, len(series))
	for i := range series {
		vs[i] = series[i].xyserie
	}
	return getDomainsXY(vs, mul)
}

func getDomainsXY(series []xyserie, mul float64) (pair, pair) {
	var x, y pair
	for i := range series {
		x.Min = getLesser(series[i].px.Min, x.Min)
		x.Max = getGreater(series[i].px.Max, x.Max)
		y.Min = getLesser(series[i].py.Min, y.Min)
		y.Max = getGreater(series[i].py.Max, y.Max)
	}
	if mul <= 1 {
		return x, y
	}
	return x.extendBy(mul), y.extendBy(mul)
}
