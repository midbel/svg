package chart

import (
	"bufio"
	"fmt"
	"io"
	"math"

	"github.com/midbel/svg"
)

const (
	fullcirc = 360.0
	halfcirc = 180.0
)

type SunburstChart struct {
	Chart
	InnerRadius int
	OuterRadius int
}

func (c SunburstChart) Render(w io.Writer, series []Hierarchy) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	if len(series) == 1 {
		series = series[0].Sub
	}
	cs := c.RenderElement(series)
	cs.Render(ws)
}

func (c SunburstChart) RenderElement(series []Hierarchy) svg.Element {
	c.checkDefault()

	var (
		cx, cy = c.GetAreaCenter()
		cs     = c.getCanvas()
		area   = c.getArea(whitstrok.Option(), svg.WithTranslate(cx, cy))
		height = float64(c.OuterRadius-c.InnerRadius) / getDepth(series)
		part   = fullcirc / getSum(series)
		angle  float64
	)
	for i := range series {
		grp := svg.NewGroup()
		series[i].Fill = getFill(i, series[i].Fill, series[i].Fill)
		c.drawSerie(&grp, series[i], angle, part, float64(height), 0)
		area.Append(grp.AsElement())
		angle += series[i].GetValue() * part
	}
	cs.Append(area.AsElement())
	cs.Append(c.drawTitle())
	cs.Append(c.drawLegend())
	return cs.AsElement()
}

func (c SunburstChart) drawSerie(a Appender, serie Hierarchy, angle, part, height, depth float64) {
	var (
		inner = height
		outer = c.distanceFromCenter() + (height * depth) + inner
		pos1  = getPosFromAngle(angle*deg2rad, outer)
		pos2  = getPosFromAngle((angle+(serie.GetValue()*part))*deg2rad, outer)
		pos3  = getPosFromAngle((angle+(serie.GetValue()*part))*deg2rad, outer-inner)
		pos4  = getPosFromAngle(angle*deg2rad, outer-inner)
		pat   = svg.NewPath(svg.WithID(serie.Label), serie.Fill.Option())
		swap  bool
	)
	if tmp := serie.GetValue() * part; tmp > halfcirc {
		swap = true
	}
	pat.AbsMoveTo(pos1)
	pat.AbsArcTo(pos2, outer, outer, 0, swap, true)
	pat.AbsLineTo(pos3)
	if pos3.X != pos4.X && pos3.Y != pos4.Y {
		pat.AbsArcTo(pos4, outer-inner, outer-inner, 0, swap, false)
	}
	pat.AbsLineTo(pos1)
	pat.Title = fmt.Sprintf("%s - %f", serie.Label, serie.GetValue())
	a.Append(pat.AsElement())

	subpart := (serie.GetValue() * part) / serie.Sum()
	for i := range serie.Sub {
		serie.Sub[i].Fill = getFill(i, serie.Sub[i].Fill, serie.Fill)
		c.drawSerie(a, serie.Sub[i], angle, subpart, height, depth+1)
		angle += serie.Sub[i].GetValue() * subpart
	}
}

func (c *SunburstChart) distanceFromCenter() float64 {
	if c.OuterRadius == c.InnerRadius {
		return 0
	}
	return float64(c.InnerRadius)
}

func (c *SunburstChart) checkDefault() {
	c.Chart.checkDefault()
	if c.OuterRadius == 0 {
		c.OuterRadius = int(math.Min(c.Width, c.Height))
	}
}

type PieChart struct {
	Chart
	OuterRadius int
	InnerRadius int
}

func (c PieChart) Render(w io.Writer, serie BarSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	cs := c.RenderElement(serie)
	cs.Render(ws)
}

func (c PieChart) RenderElement(serie BarSerie) svg.Element {
	c.checkDefault()

	var (
		cx, cy = c.GetAreaCenter()
		cs     = c.getCanvas()
		area   = c.getArea(whitstrok.Option(), svg.WithTranslate(cx, cy))
		sum    = serie.Sum()
		part   = fullcirc / sum
		angle  float64
	)
	if len(serie.Fill) == 0 {
		serie.Fill = DefaultColors
	}
	for i, v := range serie.values {
		var (
			fill = serie.Fill[i%len(serie.Fill)]
			pos1 = getPosFromAngle(angle*deg2rad, float64(c.OuterRadius))
			pos2 = getPosFromAngle((angle+(v.Value*part))*deg2rad, float64(c.OuterRadius))
			pos3 = getPosFromAngle((angle+(v.Value*part))*deg2rad, float64(c.OuterRadius-c.InnerRadius))
			pos4 = getPosFromAngle(angle*deg2rad, float64(c.OuterRadius-c.InnerRadius))
			pat  = svg.NewPath(svg.WithID(v.Label), fill.Option())
			swap bool
		)
		if tmp := v.Value * part; tmp > halfcirc {
			swap = true
		}
		pat.AbsMoveTo(pos1)
		pat.AbsArcTo(pos2, float64(c.OuterRadius), float64(c.OuterRadius), 0, swap, true)
		pat.AbsLineTo(pos3)
		if pos3.X != pos4.X && pos3.Y != pos4.Y {
			pat.AbsArcTo(pos4, float64(c.OuterRadius-c.InnerRadius), float64(c.OuterRadius-c.InnerRadius), 0, swap, false)
		}
		pat.AbsLineTo(pos1)
		pat.Title = v.Label
		area.Append(pat.AsElement())

		angle += v.Value * part
	}
	cs.Append(area.AsElement())
	cs.Append(c.drawTitle())
	cs.Append(c.drawLegend())
	return cs.AsElement()
}

func (c *PieChart) checkDefault() {
	c.Chart.checkDefault()
	if c.OuterRadius == 0 {
		c.OuterRadius = int(math.Min(c.Width, c.Height))
	}
	if c.InnerRadius == 0 {
		c.InnerRadius = c.OuterRadius
	}
}
