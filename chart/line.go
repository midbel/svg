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

type LineSerie struct {
	Title string
	Color string
	Shape ShapeType

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

type CurveStyle uint8

const (
	CurveLinear CurveStyle = iota
	CurveStep
	CurveStepBefore
	CurveStepAfter
	CurveCubic
	CurveQuadratic
)

type ShapeType uint8

const (
	ShapeDefault ShapeType = iota
	ShapeSquare
	ShapeCircle
	ShapeTriangle
	ShapeStar
	ShapeDiamond
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
	cs.Append(c.drawAxisX(c.Chart, rx))
	cs.Append(c.drawAxisY(c.Chart, ry))
	area.Append(c.drawTicksY(c.Chart, ry))
	for i := range series {
		grp := c.drawSerie(series[i], rx, ry)
		area.Append(grp)
	}

	cs.Append(area.AsElement())
	return cs.AsElement()
}

func (c ScatterChart) drawSerie(serie LineSerie, rx, ry pair) svg.Element {
	var (
		grp   = svg.NewGroup()
		dx    = c.GetAreaWidth() / rx.Diff()
		dy    = c.GetAreaHeight() / ry.Diff()
		fill  = svg.NewFill(serie.Color).Option()
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
			elem = getCircle(c.Radius, xy, strok, fill)
		case ShapeTriangle:
			pos.X -= c.Radius / 2
			pos.Y -= c.Radius / 2
			g := svg.NewGroup(svg.WithTranslate(pos.X, pos.Y))
			i := getTriangle(c.Radius, strok, fill)
			g.Append(i)
			elem = g.AsElement()
		case ShapeStar:
			pos.X -= c.Radius / 2
			pos.Y -= c.Radius / 2
			g := svg.NewGroup(svg.WithTranslate(pos.X, pos.Y))
			i := getStar(c.Radius, strok, fill)
			g.Append(i)
			elem = g.AsElement()
		case ShapeDiamond:
			pos.X -= c.Radius / 2
			pos.Y -= c.Radius / 2
			rot := svg.WithRotate(45, pos.X, pos.Y)
			elem = getDiamond(c.Radius, xy, fill, strok, rot)
		case ShapeSquare:
			pos.X -= c.Radius / 2
			pos.Y -= c.Radius / 2
			elem = getSquare(c.Radius, xy, fill, strok)
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

func getDiamond(rad float64, options ...svg.Option) svg.Element {
	options = append(options, svg.WithDimension(rad, rad))
	i := svg.NewRect(options...)
	return i.AsElement()
}

func getSquare(rad float64, options ...svg.Option) svg.Element {
	options = append(options, svg.WithDimension(rad, rad))
	i := svg.NewRect(options...)
	return i.AsElement()
}

func getTriangle(rad float64, options ...svg.Option) svg.Element {
	points := []svg.Pos{
		svg.NewPos(0, rad),
		svg.NewPos(rad/2, 0),
		svg.NewPos(rad, rad),
	}
	i := svg.NewPolygon(points, options...)
	return i.AsElement()
}

func getCircle(rad float64, options ...svg.Option) svg.Element {
	options = append(options, svg.WithRadius(rad/2))
	i := svg.NewCircle(options...)
	return i.AsElement()
}

func getStar(rad float64, options ...svg.Option) svg.Element {
	rad *= 2
	var (
		onerad   = rad / 5
		tworad   = onerad * 2
		threerad = onerad * 3
		fourrad  = onerad * 4
		halfrad  = rad / 2
	)
	points := []svg.Pos{
		svg.NewPos(onerad, rad),
		svg.NewPos(tworad, halfrad),
		svg.NewPos(0, tworad),
		svg.NewPos(tworad, tworad),
		svg.NewPos(halfrad, 0),
		svg.NewPos(threerad, tworad),
		svg.NewPos(rad, tworad),
		svg.NewPos(threerad, halfrad),
		svg.NewPos(fourrad, rad),
		svg.NewPos(halfrad, threerad),
		svg.NewPos(onerad, rad),
	}
	i := svg.NewPolygon(points, options...)
	return i.AsElement()
}
