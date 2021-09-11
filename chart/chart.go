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

type appender interface {
	Append(svg.Element)
}

type Chart struct {
	Title  string
	Width  float64
	Height float64
	Padding

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

func (c *Chart) getCanvas() svg.SVG {
	var (
		dim = svg.NewDim(c.Width, c.Height)
		cs  = svg.NewSVG(dim.Option())
		bg  = svg.NewGroup(svg.WithClass("bg-chart"))
	)

	if !c.Background.IsZero() && !c.Padding.IsZero() {
		var (
			d = svg.NewDim(c.Width, c.Height)
			r = svg.NewRect(d.Option(), c.Border.Option(), c.Background.Option())
		)
		bg.Append(r.AsElement())
	}
	if !c.Area.IsZero() {
		var (
			p = svg.NewPos(c.Left, c.Top)
			d = svg.NewDim(c.GetAreaWidth(), c.GetAreaHeight())
			r = svg.NewRect(p.Option(), d.Option(), c.Area.Option())
		)
		bg.Append(r.AsElement())
	}
	cs.Append(bg.AsElement())
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

func formatTime(t time.Time) string {
	return t.Format("15:04:05")
}

func formatFloat(val float64) string {
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
