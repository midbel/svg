package chart

import (
	"bufio"
	"io"
	"math"

	"github.com/midbel/svg"
	"github.com/midbel/svg/colors"
)

const (
	fullcirc = 360.0
	halfcirc = 180.0
	deg2rad  = math.Pi / halfcirc
)

type PieChart struct {
	Chart
	OutRadius int
	InRadius  int
}

func (c PieChart) Render(w io.Writer, serie Serie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	cs := c.RenderElement(serie)
	cs.Render(ws)
}

func (c PieChart) RenderElement(serie Serie) svg.Element {
	c.checkDefault()

	var (
		dim    = svg.NewDim(c.Width, c.Height)
		cs     = svg.NewSVG(dim.Option())
		sum    = serie.Sum()
		part   = fullcirc / sum
		cx, cy = c.GetAreaCenter()
		area   = svg.NewGroup(svg.WithID("area"), svg.WithTranslate(cx, cy))
		angle  float64
	)
	for i, v := range serie.values {
		var (
			fill = svg.NewFill(colors.RdYlBu11[i%len(colors.RdYlBu11)])
			pos1 = getPosFromAngle(angle*deg2rad, float64(c.OutRadius))
			pos2 = getPosFromAngle((angle+(v.Value*part))*deg2rad, float64(c.OutRadius))
			pos3 = getPosFromAngle((angle+(v.Value*part))*deg2rad, float64(c.OutRadius-c.InRadius))
			pos4 = getPosFromAngle(angle*deg2rad, float64(c.OutRadius-c.InRadius))
			pat  = svg.NewPath(svg.WithID(v.Label), fill.Option(), whitstrok.Option())
			swap bool
		)
		if tmp := v.Value * part; tmp > halfcirc {
			swap = true
		}
		pat.AbsMoveTo(pos1)
		pat.AbsArcTo(pos2, float64(c.OutRadius), float64(c.OutRadius), 0, swap, true)
		pat.AbsLineTo(pos3)
		if pos3.X != pos4.X && pos3.Y != pos4.Y {
			pat.AbsArcTo(pos4, float64(c.OutRadius-c.InRadius), float64(c.OutRadius-c.InRadius), 0, swap, false)
		}
		pat.AbsLineTo(pos1)
		pat.Title = v.Label
		area.Append(pat.AsElement())

		angle += v.Value * part
	}
	cs.Append(area.AsElement())
	return cs.AsElement()
}

func (c *PieChart) checkDefault() {
	c.Chart.checkDefault()
	if c.OutRadius == 0 {
		c.OutRadius = int(math.Min(c.Width, c.Height))
	}
	if c.InRadius == 0 {
		c.InRadius = c.OutRadius
	}
}

func getPosFromAngle(angle, radius float64) svg.Pos {
	var (
		x1 = float64(radius) * math.Cos(angle)
		y1 = float64(radius) * math.Sin(angle)
	)
	return svg.NewPos(x1, y1)
}
