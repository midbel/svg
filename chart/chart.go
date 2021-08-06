package chart

import (
	"hash/adler32"

	"github.com/midbel/svg"
)

const (
	DefaultWidth  = 800
	DefaultHeight = 600
)

var DefaultColours = []string{"yellow", "orange", "red", "purple", "blue", "green"}

type Chart struct {
	Width  float64
	Height float64
	Padding
	Colours []string
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
	if len(c.Colours) == 0 {
		c.Colours = DefaultColours
	}
}

func (c *Chart) peekFillFromNumber(str int64) svg.Option {
	set := c.Colours
	if len(set) == 0 {
		set = svg.Colours
	}
	col := set[int(str)%len(set)]
	return svg.NewFill(col).Option()
}

func (c *Chart) peekFillFromString(str string) svg.Option {
	set := c.Colours
	if len(set) == 0 {
		set = svg.Colours
	}
	var (
		sum = adler32.Checksum([]byte(str))
		col = set[int(sum)%len(set)]
	)
	return svg.NewFill(col).Option()
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
