package chart

import (
	"bufio"
	"io"
	"math"

	"github.com/midbel/svg"
	"github.com/midbel/svg/colors"
)

const (
	circle  = 360.0
	deg2rad = math.Pi / 180
)

type PieChart struct {
	Chart
	Radius int
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
		cx, cy = c.GetAreaCenter()
		area   = svg.NewGroup(svg.WithID("area"), svg.WithTranslate(cx, cy))
		sum    = serie.Sum()
		part   = circle / sum
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
			pos1 = getPosFromAngle(angle*deg2rad, float64(c.Radius))
			pos2 = getPosFromAngle((angle+(next.Value*part))*deg2rad, float64(c.Radius))
			pat  = svg.NewPath(svg.WithID(v.Label), fill.Option(), whitstrok.Option())
		)
		pat.AbsMoveTo(svg.NewPos(0, 0))
		pat.AbsLineTo(pos1)
		pat.AbsArcTo(pos2, float64(c.Radius), float64(c.Radius), 0, false, true)
		pat.AbsLineTo(svg.NewPos(0, 0))
		area.Append(pat.AsElement())
	}
	cs.Append(area.AsElement())
	cs.Render(w)
}

func getPosFromAngle(angle, radius float64) svg.Pos {
	var (
		x1 = float64(radius) * math.Cos(angle)
		y1 = float64(radius) * math.Sin(angle)
	)
	return svg.NewPos(x1, y1)
}
