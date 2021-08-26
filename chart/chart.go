package chart

import (
	"math"
	"strconv"
	"time"

	"github.com/midbel/svg"
)

const (
	DefaultWidth  = 800
	DefaultHeight = 600
)

type appender interface {
	Append(svg.Element)
}

type Chart struct {
	Width  float64
	Height float64
	Padding

	GetColor  func(string, int) svg.Fill
	GetStroke func(string, int) svg.Stroke
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

func (c *Chart) checkDefault() {
	if c.Width == 0 {
		c.Width = DefaultWidth
	}
	if c.Height == 0 {
		c.Height = DefaultHeight
	}

	if c.GetColor == nil {
		c.GetColor = defaultFill
	}
	if c.GetStroke == nil {
		c.GetStroke = defaultStroke
	}
}

func defaultFill(_ string, _ int) svg.Fill {
	return svg.NewFill("steelblue")
}

func defaultStroke(_ string, _ int) svg.Stroke {
	return svg.NewStroke("black", 1)
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

var (
	tickstrok = svg.NewStroke("lightgrey", 1)
	axisstrok = svg.NewStroke("black", 1)
	whitstrok = svg.NewStroke("white", 1)
	linestrok = svg.NewStroke("steelblue", 1)
	nonefill  = svg.NewFill("none")
)

func getRect(options ...svg.Option) svg.Rect {
	options = append(options)
	return svg.NewRect(options...)
}

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
