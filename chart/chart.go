package chart

import (
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
