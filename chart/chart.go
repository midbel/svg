package chart

import (
	"math"

	"github.com/midbel/svg"
	"github.com/midbel/svg/colors"
)

const (
	DefaultWidth  = 800
	DefaultHeight = 600
)

var DefaultColors []svg.Fill

func init() {
	DefaultColors = make([]svg.Fill, len(colors.Paired10))
	for i := range colors.Paired10 {
		DefaultColors[i] = svg.NewFill(colors.Paired10[i])
	}
}

// type Serie interface {
// 	X() Axis
// 	Y() Axis
// }

type Common struct {
	Title string

	XAxis Axis
	YAxis Axis
}

func makeCommon(title string) Common {
	return Common{
		Title: title,
	}
}

func (c Common) hasAxis() bool {
	return c.XAxis != nil || c.YAxis != nil
}

type Appender interface {
	Append(svg.Element)
}

type Position uint8

const (
	Top Position = 1 << iota
	TopLeft
	Left
	BottomLeft
	Bottom
	BottomRight
	Right
	TopRight
)

func (p Position) Horizontal() bool {
	return p == Top || p == Bottom
}

func (p Position) Vertical() bool {
	return p == Left || p == Right
}

func (p Position) canLeft() bool {
	return p == BottomLeft || p == Left || p == TopLeft || p == Top
}

func (p Position) canTop() bool {
	return p.canLeft()
}

func (p Position) canRight() bool {
	return p == BottomRight || p == Right || p == TopRight || p == Bottom
}

func (p Position) canBottom() bool {
	return p.canRight()
}

func (p Position) adjust(pos svg.Pos) svg.Pos {
	switch p {
	case Top:
		pos.Y = -pos.Y
	case Bottom:
	case Left:
	case Right:
		pos.X = -pos.X
	default:
	}
	return pos
}

type Legend struct {
	Show bool
	Position
}

type Chart struct {
	Common

	Width  float64
	Height float64
	Legend
	Padding
	Axis struct {
		Top    Axis
		Left   Axis
		Bottom Axis
		Right  Axis
	}

	Border     svg.Stroke
	Background svg.Fill
	Area       svg.Fill
}

func (c *Chart) GetAreaWidth() float64 {
	c.checkDefault()
	return c.Width - c.Horizontal()
}

func (c *Chart) GetAreaHeight() float64 {
	c.checkDefault()
	return c.Height - c.Vertical()
}

func (c *Chart) GetAreaCenter() (float64, float64) {
	return c.GetAreaWidth() / 2, c.GetAreaHeight() / 2
}

func (c *Chart) getArea(options ...svg.Option) svg.Group {
	os := []svg.Option{
		svg.WithClass("area"),
		c.translate(),
	}
	return svg.NewGroup(append(os, options...)...)
}

func (c *Chart) drawTitle() svg.Element {
	if c.Title == "" {
		return nil
	}
	y := 16.0
	if c.Padding.Top > 0 {
		y = c.Padding.Top / 2
	}
	var (
		pos  = svg.NewPos(c.Width/2, y)
		font = svg.NewFont(16)
		anc  = svg.WithAnchor("middle")
		base = svg.WithDominantBaseline("middle")
		text = svg.NewText(c.Title, pos.Option(), font.Option(), anc, base)
	)
	return text.AsElement()
}

func (c *Chart) drawLegend() svg.Element {
	return nil
}

func (c *Chart) drawDefaultAxis() svg.Element {
	var (
		ap = svg.NewGroup(svg.WithClass("axis"), svg.WithTranslate(c.Padding.Left, c.Padding.Top))
		ok = c.hasAxis()
	)
	if c.YAxis != nil {
		var (
			opt    = svg.WithTranslate(0, 0)
			width  = c.GetAreaWidth()
			height = c.GetAreaHeight()
		)
		if !c.YAxis.Left() {
			opt = svg.WithTranslate(c.GetAreaWidth(), 0)
			width, height = height, width
		}
		grp := svg.NewGroup(opt)
		ap.Append(grp.AsElement())
		c.YAxis.Draw(&grp, c.GetAreaHeight(), c.GetAreaWidth())
	}
	if c.XAxis != nil {
		var (
			opt    = svg.WithTranslate(0, c.GetAreaHeight())
			width  = c.GetAreaWidth()
			height = c.GetAreaHeight()
		)
		if !c.XAxis.Bottom() {
			opt = svg.WithTranslate(0, 0)
		}
		grp := svg.NewGroup(opt)
		ap.Append(grp.AsElement())
		c.XAxis.Draw(&grp, width, height)
	}
	if !ok {
		return nil
	}
	return ap.AsElement()
}

func (c *Chart) drawAxisBis(axis []Axis) svg.Element {
	if len(axis) == 0 {
		return nil
	}
	ap := svg.NewGroup(svg.WithClass("axis"), svg.WithTranslate(c.Padding.Left, c.Padding.Top))
	for i := range axis {
		c.appendAxis(&ap, axis[i])
	}
	return ap.AsElement()
}

func (c *Chart) appendAxis(ap Appender, a Axis) {
	if a.Vertical() {
		var (
			opt    = svg.WithTranslate(0, 0)
			width  = c.GetAreaWidth()
			height = c.GetAreaHeight()
		)
		if !a.Left() {
			opt = svg.WithTranslate(c.GetAreaWidth(), 0)
			width, height = height, width
		}
		grp := svg.NewGroup(opt)
		ap.Append(grp.AsElement())
		a.Draw(&grp, c.GetAreaHeight(), c.GetAreaWidth())
	}
	if a.Horizontal() {
		var (
			opt    = svg.WithTranslate(0, c.GetAreaHeight())
			width  = c.GetAreaWidth()
			height = c.GetAreaHeight()
		)
		if !a.Bottom() {
			opt = svg.WithTranslate(0, 0)
		}
		grp := svg.NewGroup(opt)
		ap.Append(grp.AsElement())
		a.Draw(&grp, width, height)
	}
}

func (c *Chart) drawAxis(rx, ry AxisOption, options ...AxisOption) svg.Element {
	if c.Axis.Top == nil && c.Axis.Bottom == nil && c.Axis.Left == nil && c.Axis.Right == nil {
		return nil
	}
	ap := svg.NewGroup(svg.WithClass("axis"), svg.WithTranslate(c.Padding.Left, c.Padding.Top))
	if c.Axis.Top != nil {
		grp := svg.NewGroup()
		ap.Append(grp.AsElement())

		options = append(options, rx, WithPosition(Top))
		c.Axis.Top.update(options...)
		c.Axis.Top.Draw(&grp, c.GetAreaWidth(), c.GetAreaHeight())
	}
	if c.Axis.Bottom != nil {
		grp := svg.NewGroup(svg.WithTranslate(0, c.GetAreaHeight()))
		ap.Append(grp.AsElement())

		options = append(options, rx, WithPosition(Bottom))
		c.Axis.Bottom.update(options...)
		c.Axis.Bottom.Draw(&grp, c.GetAreaWidth(), c.GetAreaHeight())
	}
	if c.Axis.Left != nil {
		grp := svg.NewGroup()
		ap.Append(grp.AsElement())

		options = append(options, ry, WithPosition(Left))
		c.Axis.Left.update(options...)
		c.Axis.Left.Draw(&grp, c.GetAreaHeight(), c.GetAreaWidth())
	}
	if c.Axis.Right != nil {
		grp := svg.NewGroup(svg.WithTranslate(c.GetAreaWidth(), 0))
		ap.Append(grp.AsElement())

		options = append(options, ry, WithPosition(Right))
		c.Axis.Right.update(options...)
		c.Axis.Right.Draw(&grp, c.GetAreaHeight(), c.GetAreaWidth())
	}
	return ap.AsElement()
}

func (c *Chart) getCanvas() svg.SVG {
	var (
		dim = svg.NewDim(c.Width, c.Height)
		cs  = svg.NewSVG(dim.Option())
		bg  = svg.NewGroup(svg.WithClass("bg-chart"))
		ok  bool
	)

	if !c.Background.IsZero() && !c.Padding.IsZero() {
		var (
			d = svg.NewDim(c.Width, c.Height)
			r = svg.NewRect(d.Option(), c.Border.Option(), c.Background.Option())
		)
		bg.Append(r.AsElement())
		ok = true
	}
	if !c.Area.IsZero() {
		var (
			p = svg.NewPos(c.Left, c.Top)
			d = svg.NewDim(c.GetAreaWidth(), c.GetAreaHeight())
			r = svg.NewRect(p.Option(), d.Option(), c.Area.Option())
		)
		bg.Append(r.AsElement())
		ok = true
	}
	if ok {
		cs.Append(bg.AsElement())
	}
	return cs
}

func (c *Chart) checkDefault() {
	if c.Width == 0 {
		c.Width = DefaultWidth
	}
	if c.Height == 0 {
		c.Height = DefaultHeight
	}
}

func (c *Chart) getOptionsAxisX() []svg.Option {
	if c.Horizontal() == 0 {
		return nil
	}
	return []svg.Option{
		svg.WithID("x-axis"),
		svg.WithClass("axis"),
		svg.WithTranslate(c.Padding.Left, c.Height-c.Padding.Bottom),
	}
}

func (c *Chart) getOptionsAxisY() []svg.Option {
	if c.Vertical() == 0 {
		return nil
	}
	return []svg.Option{
		svg.WithID("y-axis"),
		svg.WithClass("axis"),
		c.translate(),
	}
}

type Padding struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}

func CreatePadding(horiz, vert float64) Padding {
	return Padding{
		Left:   horiz,
		Right:  horiz,
		Top:    vert,
		Bottom: vert,
	}
}

func (p Padding) Horizontal() float64 {
	return p.Left + p.Right
}

func (p Padding) Vertical() float64 {
	return p.Top + p.Bottom
}

func (p Padding) translate() svg.Option {
	return svg.WithTranslate(p.Left, p.Top)
}

func (p Padding) IsZero() bool {
	return p.Top == 0 && p.Bottom == 0 && p.Right == 0 && p.Left == 0
}

func getFill(i int, fill, other svg.Fill) svg.Fill {
	if !fill.IsZero() {
		return fill
	}
	if !other.IsZero() {
		return other
	}
	return DefaultColors[i%len(DefaultColors)]
}

var (
	whitstrok = svg.NewStroke("white", 0.5)
	nonefill  = svg.NewFill("none")
)

func getLesser(v1, v2 float64) float64 {
	return math.Min(v1, v2)
}

func getGreater(v1, v2 float64) float64 {
	return math.Max(v1, v2)
}

func getPathLine(color string) svg.Path {
	var (
		fill  = svg.NewFill("none")
		strok = svg.NewStroke(color, 2)
	)
	fill.Opacity = 0
	return svg.NewPath(fill.Option(), strok.Option(), svg.WithClass("line"))
}

const deg2rad = math.Pi / halfcirc

func getPosFromAngle(angle, radius float64) svg.Pos {
	var (
		x1 = float64(radius) * math.Cos(angle)
		y1 = float64(radius) * math.Sin(angle)
	)
	return svg.NewPos(x1, y1)
}
