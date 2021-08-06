package chart

import (
	"github.com/midbel/svg"
)

const (
	DefaultWidth  = 1360
	DefaultHeight = 768
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

func (p Padding) Horizontal() float64 {
  return p.Left+p.Right
}

func (p Padding) Vertical() float64 {
  return p.Top+p.Bottom
}

func (p Padding) translate() svg.Option {
	return svg.WithTranslate(p.Left, p.Top)
}
