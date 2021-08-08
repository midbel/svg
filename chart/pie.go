package chart

import (
	"bufio"
	"io"
	"math"

	"github.com/midbel/svg"
	"github.com/midbel/svg/colors"
)

const (
	circum  = 360.0
	deg2rad = math.Pi / 180
)

type PieChart struct {
	Chart
	MaxRadius int
	MinRadius int
}

func (c PieChart) Render(w io.Writer, serie Serie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	c.render(ws, serie)
}

func (c PieChart) render(w svg.Writer, serie Serie) {
	c.checkDefault()

	var (
		dim    = svg.NewDim(c.Width, c.Height)
		cs     = svg.NewSVG(dim.Option())
		sum    = serie.Sum()
		part   = circum / sum
		cx, cy = c.GetAreaCenter()
		area   = svg.NewGroup(svg.WithID("area"), svg.WithTranslate(cx, cy))
		angle  float64
	)
	for i, v := range serie.values {
		var next valuelabel
		if i == serie.Len()-1 {
			next = serie.values[0]
		} else {
			next = serie.values[i+1]
		}
		angle += v.Value * part
		var (
			fill = svg.NewFill(colors.Set26[i%len(colors.Set26)])
			pos1 = getPosFromAngle(angle*deg2rad, float64(c.MaxRadius))
			pos2 = getPosFromAngle((angle+(next.Value*part))*deg2rad, float64(c.MaxRadius))
			pos3 = getPosFromAngle((angle+(next.Value*part))*deg2rad, float64(c.MaxRadius-c.MinRadius))
			pos4 = getPosFromAngle(angle*deg2rad, float64(c.MaxRadius-c.MinRadius))
			pat  = svg.NewPath(svg.WithID(v.Label), fill.Option(), whitstrok.Option())
		)
		pat.AbsMoveTo(pos1)
		pat.AbsArcTo(pos2, float64(c.MaxRadius), float64(c.MaxRadius), 0, false, true)
		pat.AbsLineTo(pos3)
		pat.AbsArcTo(pos4, float64(c.MaxRadius-c.MinRadius), float64(c.MaxRadius-c.MinRadius), 0, false, false)
		pat.AbsLineTo(pos1)
		area.Append(pat.AsElement())
	}
	cs.Append(area.AsElement())
	cs.Render(w)
}

func (c *PieChart) checkDefault() {
	c.Chart.checkDefault()
	if c.MaxRadius == 0 {
		c.MaxRadius = int(math.Min(c.Width, c.Height))
	}
	if c.MinRadius == 0 {
		c.MinRadius = c.MaxRadius
	}
}

func getPosFromAngle(angle, radius float64) svg.Pos {
	var (
		x1 = float64(radius) * math.Cos(angle)
		y1 = float64(radius) * math.Sin(angle)
	)
	return svg.NewPos(x1, y1)
}
