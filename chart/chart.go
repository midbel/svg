package chart

import (
	"math"
	"strconv"

	"github.com/midbel/svg"
)

const (
	DefaultWidth  = 800
	DefaultHeight = 600
)

type Chart struct {
	Width  float64
	Height float64
	Padding
}

func (c *Chart) GetAreaWidth() float64 {
	c.checkDefault()
	return c.Width - c.Horizontal()
}

func (c *Chart) GetAreaHeight() float64 {
	c.checkDefault()
	return c.Height - c.Vertical()
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

var (
	tickstrok = svg.NewStroke("lightgrey", 1)
	axisstrok = svg.NewStroke("black", 1)
	whitstrok = svg.NewStroke("white", 1)
	linestrok = svg.NewStroke("steelblue", 1)
)

const (
	ticklen = 7
	textick = 18
)

func getRect(options ...svg.Option) svg.Rect {
	options = append(options, whitstrok.Option())
	return svg.NewRect(options...)
}

func getTick(pos1, pos2 svg.Pos) svg.Element {
	tickstrok.Dash.Array = []int{5}
	line := svg.NewLine(pos1, pos2, tickstrok.Option())
	return line.AsElement()
}

func getLesser(v1, v2 float64) float64 {
	if math.IsNaN(v1) || v2 < v1 {
		return v2
	}
	return v1
}

func getGreater(v1, v2 float64) float64 {
	if math.IsNaN(v1) || v2 > v1 {
		return v2
	}
	return v1
}

func getPathLine(stk string) svg.Path {
	var (
		fill  = svg.NewFill("transparent")
		strok = svg.NewStroke(stk, 1)
	)
	fill.Opacity = 0
	return svg.NewPath(fill.Option(), strok.Option())
}

func formatFloat(val float64) string {
	return strconv.FormatFloat(val, 'f', 2, 64)
}
