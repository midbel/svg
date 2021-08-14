package chart

import (
	"bufio"
	"io"

	"github.com/midbel/svg"
)

type Hierarchy struct {
	Label string
	Value float64
	Sub   []Hierarchy
}

func (h Hierarchy) Sum() float64 {
	if h.isLeaf() {
		return h.Value
	}
	var s float64
	for i := range h.Sub {
		s += h.Sub[i].Value
	}
	return s
}

func (h Hierarchy) Depth() int {
	if h.isLeaf() {
		return 1
	}
	var d int
	for i := range h.Sub {
		x := h.Sub[i].Depth()
		if x > d {
			d = x
		}
	}
	return d + 1
}

func (h Hierarchy) Len() int {
	return len(h.Sub)
}

func (h Hierarchy) isLeaf() bool {
	return h.Len() == 0
}

type HeatmapChart struct {
	Chart
}

type TilingMethod uint8

const (
	TilingDefault TilingMethod = iota
	TilingBinary
	TilingVertical
	TilingHorizontal
	TilingAlternate
)

type TreemapChart struct {
	Chart
	Tiling TilingMethod
}

func (c TreemapChart) Render(w io.Writer, series []Hierarchy) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()
	if len(series) == 1 {
		series = series[0].Sub
	}
	cs := c.RenderElement(series)
	cs.Render(ws)
}

func (c TreemapChart) RenderElement(series []Hierarchy) svg.Element {
	c.checkDefault()
	var (
		dim  = svg.NewDim(c.Width, c.Height)
		cs   = svg.NewSVG(dim.Option())
		area = svg.NewGroup(svg.WithID("area"), c.translate())
	)
	switch c.Tiling {
	case TilingDefault:
		c.drawDefault(&area, series)
	case TilingBinary:
		c.drawBinary(&area, series)
	case TilingVertical:
		c.drawVertical(&area, series, c.GetAreaHeight()/getSum(series))
	case TilingHorizontal:
		c.drawHorizontal(&area, series, c.GetAreaWidth()/getSum(series))
	case TilingAlternate:
	default:
	}
	cs.Append(area.AsElement())
	return cs.AsElement()
}

func (c TreemapChart) drawHorizontal(a appender, series []Hierarchy, part float64) {
	var off float64
	for i := range series {
		if series[i].isLeaf() {
			var (
				pos   = svg.NewPos(off, 0)
				dim   = svg.NewDim(series[i].Value*part, c.GetAreaHeight())
				fill  = svg.NewFill("steelblue")
				strok = svg.NewStroke("white", 2)
			)
			r := svg.NewRect(pos.Option(), dim.Option(), fill.Option(), strok.Option())
			r.Title = series[i].Label
			a.Append(r.AsElement())
		} else {
			grp := svg.NewGroup(svg.WithTranslate(off, 0))
			c.drawHorizontal(&grp, series[i].Sub, (series[i].Value*part)/getSum(series[i].Sub))
			a.Append(grp.AsElement())
		}
		off += series[i].Value * part
	}
}

func (c TreemapChart) drawVertical(a appender, series []Hierarchy, part float64) {
	var off float64
	for i := range series {
		if series[i].isLeaf() {
			var (
				pos   = svg.NewPos(0, off)
				dim   = svg.NewDim(c.GetAreaWidth(), series[i].Value*part)
				fill  = svg.NewFill("steelblue")
				strok = svg.NewStroke("white", 2)
			)
			r := svg.NewRect(pos.Option(), dim.Option(), fill.Option(), strok.Option())
			r.Title = series[i].Label
			a.Append(r.AsElement())
		} else {
			grp := svg.NewGroup(svg.WithTranslate(0, off))
			c.drawVertical(&grp, series[i].Sub, (series[i].Value*part)/getSum(series[i].Sub))
			a.Append(grp.AsElement())
		}
		off += series[i].Value * part
	}
}

func (c TreemapChart) drawBinary(a appender, series []Hierarchy) {

}

func (c TreemapChart) drawDefault(a appender, series []Hierarchy) {

}

func getDepth(series []Hierarchy) float64 {
	var d int
	for i := range series {
		x := series[i].Depth()
		if x > d {
			d = x
		}
	}
	return float64(d)
}

func getSum(series []Hierarchy) float64 {
	var sum float64
	for i := range series {
		sum += series[i].Value
	}
	if sum == 0 {
		return 1
	}
	return sum
}
