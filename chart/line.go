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

type LineSerie struct {
	xyserie

	Curver
	svg.Stroke
	svg.Fill
	Shape ShapeType
}

func NewLineSerie(title string) LineSerie {
	return LineSerie{
		xyserie: xyserie{Common: makeCommon(title)},
	}
}

func (is *LineSerie) GetStroke() svg.Stroke {
	return is.Stroke
}

func (is *LineSerie) GetFill() svg.Fill {
	return is.Fill
}

type ScatterSerie struct {
	xyserie

	svg.Fill
	svg.Stroke
	Size  float64
	Shape ShapeType
}

func NewScatterSerie(title string) ScatterSerie {
	c := xyserie{Common: makeCommon(title)}
	return ScatterSerie{
		xyserie: c,
		Size:    DefaultSize,
	}
}

type AreaSerie struct {
	Title string
	svg.Stroke
	svg.Fill

	serie1 LineSerie
	serie2 LineSerie
}

func NewAreaSerie(title string, s1, s2 LineSerie) AreaSerie {
	return AreaSerie{
		Title:  title,
		serie1: s1,
		serie2: s2,
	}
}

type AreaChart struct {
	Chart
}

func (c AreaChart) Render(w io.Writer, serie AreaSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	cs := c.RenderElement(serie)
	cs.Render(ws)
}

func (c AreaChart) RenderElement(serie AreaSerie) svg.Element {
	c.checkDefault()

	var (
		cs     = c.getCanvas()
		area   = c.getArea()
		rx, ry = getLineDomains([]LineSerie{serie.serie1, serie.serie2}, 1.1)
	)

	cs.Append(c.Chart.drawAxis(rx.AxisRange(), ry.AxisRange()))

	area.Append(c.drawSerie(serie, rx, ry))
	cs.Append(area.AsElement())
	cs.Append(c.drawTitle())
	cs.Append(c.drawLegend())
	return cs.AsElement()
}

func (c AreaChart) drawSerie(serie AreaSerie, rx, ry Range) svg.Element {
	var (
		dx  = c.GetAreaWidth() / rx.Diff()
		dy  = c.GetAreaHeight() / ry.Diff()
		off float64
		pat = svg.NewPath(serie.Fill.Option(), serie.Stroke.Option())
		pos svg.Pos
	)
	if ry.First() > 0 {
		off = ry.First() * dy
	}
	pos.X = (serie.serie1.values[0].X - rx.First()) * dx
	pos.Y = off + c.GetAreaHeight() - (serie.serie1.values[0].Y * dy)
	if ry.First() < 0 {
		pos.Y -= math.Abs(ry.First()) * dy
	}
	pat.AbsMoveTo(pos)
	for i := 1; i < serie.serie1.Len(); i++ {
		pos.X = (serie.serie1.values[i].X - rx.First()) * dx
		pos.Y = off + c.GetAreaHeight() - (serie.serie1.values[i].Y * dy)
		if ry.First() < 0 {
			pos.Y -= math.Abs(ry.First()) * dy
		}
		pat.AbsLineTo(pos)
	}
	for i := serie.serie2.Len() - 1; i >= 0; i-- {
		pos.X = (serie.serie2.values[i].X - rx.First()) * dx
		pos.Y = off + c.GetAreaHeight() - (serie.serie2.values[i].Y * dy)
		if ry.First() < 0 {
			pos.Y -= math.Abs(ry.First()) * dy
		}
		pat.AbsLineTo(pos)
	}
	pat.ClosePath()
	return pat.AsElement()
}

type ScatterChart struct {
	Chart
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
		cs     = c.getCanvas()
		area   = c.getArea()
		rx, ry = getScatterDomains(series, 1.15)
	)
	cs.Append(c.Chart.drawAxis(rx.AxisRange(), ry.AxisRange()))
	for i := range series {
		grp := c.drawSerie(series[i], rx, ry)
		area.Append(grp)
	}

	cs.Append(area.AsElement())
	cs.Append(c.drawTitle())
	cs.Append(c.drawLegend())
	return cs.AsElement()
}

func (c ScatterChart) drawSerie(serie ScatterSerie, rx, ry Range) svg.Element {
	var (
		fill = serie.Fill.Option()
		grp  = svg.NewGroup(fill)
		dx   = c.GetAreaWidth() / rx.Diff()
		dy   = c.GetAreaHeight() / ry.Diff()
		pos  svg.Pos
	)
	for i := 0; i < serie.Len(); i++ {
		pos.X = (serie.values[i].X - rx.First()) * dx
		pos.Y = c.GetAreaHeight() - (serie.values[i].Y * dy)
		if ry.First() < 0 {
			pos.Y -= math.Abs(ry.First()) * dy
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
	return grp.AsElement()
}

type LineChart struct {
	Chart
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
		cs     = c.getCanvas()
		area   = c.getArea()
		rx, ry = getLineDomains(series, 1)
	)
	if c.XAxis != nil {
		c.XAxis.update(rx.Domain())
	}
	if c.YAxis != nil {
		c.YAxis.update(ry.Domain())
	}
	ry = ry.extendBy(1.1)
	cs.Append(c.drawDefaultAxis())
	for i := range series {
		if series[i].Curver == nil {
			series[i].Curver = LinearCurve()
		}
		if series[i].XAxis != nil {
			series[i].XAxis.update(series[i].px.Domain())
		}
		if series[i].YAxis != nil {
			series[i].YAxis.update(series[i].py.Domain())
		}
		elem := series[i].Curver.Draw(c.Chart, &series[i], rx, ry)
		area.Append(elem)
	}
	cs.Append(area.AsElement())
	cs.Append(c.drawTitle())
	cs.Append(c.drawLegend())
	return cs.AsElement()
}

func (c *LineChart) checkDefault() {
	c.Chart.checkDefault()
}

type pair struct {
	Min float64
	Max float64
}

func (p pair) AxisRange() AxisOption {
	return p.Domain()
}

func (p pair) Domain() AxisOption {
	return WithNumberRange(p.Min, p.Max)
}

func (p pair) First() float64 {
	return p.Min
}

func (p pair) Last() float64 {
	return p.Max
}

func (p pair) Diff() float64 {
	return p.Max - p.Min
}

func (p pair) extendBy(by float64) pair {
	p.Min -= (math.Abs(p.Min) * by) - math.Abs(p.Min)
	p.Max *= by
	return p
}

type xypoint struct {
	X float64
	Y float64
}

type xyserie struct {
	Common

	values []xypoint
	px     pair
	py     pair
}

func (xy *xyserie) At(i int) Point {
	var p Point
	if i >= 0 && i < len(xy.values) {
		p.X = xy.values[i].X
		p.Y = xy.values[i].Y
	}
	return p
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
	vp := xypoint{
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
	var (
		x = series[0].px
		y = series[0].py
	)
	for i := 1; i < len(series); i++ {
		x.Min = getLesser(series[i].px.Min, x.Min)
		x.Max = getGreater(series[i].px.Max, x.Max)
		y.Min = getLesser(series[i].py.Min, y.Min)
		y.Max = getGreater(series[i].py.Max, y.Max)
	}
	if mul <= 1 {
		mul = 1
	}
	return x.extendBy(mul), y.extendBy(mul)
}
