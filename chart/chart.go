package chart

import (
	"math"
	"strconv"
	"time"

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

type Appender interface {
	Append(svg.Element)
}

type Axis interface {
	Draw(Appender, float64, ...svg.Option)
	update(options ...AxisOption)
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

type Orientation uint8

const (
	Horizontal Orientation = 1 << iota
	Vertical
)

type Legend struct {
	Show bool
	Orientation
	Position
}

type Chart struct {
	Title  string
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

func (c *Chart) drawAxis(options ...AxisOption) svg.Element {
	if c.Axis.Top == nil && c.Axis.Bottom == nil && c.Axis.Left == nil && c.Axis.Right == nil {
		return nil
	}
	ap := svg.NewGroup(svg.WithClass("axis"), svg.WithTranslate(c.Padding.Left, c.Padding.Right))
	if c.Axis.Top != nil {
		grp := svg.NewGroup()
		ap.Append(grp.AsElement())

		c.Axis.Top.update(options...)
		c.Axis.Top.Draw(&grp, c.GetAreaWidth())
	}
	if c.Axis.Left != nil {
		grp := svg.NewGroup()
		ap.Append(grp.AsElement())

		c.Axis.Left.update(options...)
		c.Axis.Left.Draw(&grp, c.GetAreaHeight())
	}
	if c.Axis.Bottom != nil {
		grp := svg.NewGroup(svg.WithTranslate(0, c.GetAreaHeight()))
		ap.Append(grp.AsElement())

		options = append(options, withOrientation(Horizontal), withPosition(Top))
		c.Axis.Bottom.update(options...)
		c.Axis.Bottom.Draw(&grp, c.GetAreaWidth())
	}
	if c.Axis.Right != nil {
		grp := svg.NewGroup(svg.WithTranslate(0, c.GetAreaWidth()))
		ap.Append(grp.AsElement())

		c.Axis.Right.update(options...)
		c.Axis.Right.Draw(&grp, c.GetAreaHeight())
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
	tickstrok = svg.NewStroke("darkgray", 1)
	axisstrok = svg.NewStroke("darkgray", 1)
	whitstrok = svg.NewStroke("white", 1)
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

func formatDay(t time.Time) string {
	return t.Format("2006-01-02")
}

func formatTime(t time.Time, _ int) string {
	return t.Format("15:04:05")
}

func formatFloat(val float64, _ int) string {
	if almostZero(val) {
		return "0.00"
	}
	return strconv.FormatFloat(val, 'f', 2, 64)
}

const threshold = 1e-9

func almostZero(val float64) bool {
	return math.Abs(val-0) <= threshold
}

const deg2rad = math.Pi / halfcirc

func getPosFromAngle(angle, radius float64) svg.Pos {
	var (
		x1 = float64(radius) * math.Cos(angle)
		y1 = float64(radius) * math.Sin(angle)
	)
	return svg.NewPos(x1, y1)
}
