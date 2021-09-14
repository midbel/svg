package chart

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"sort"

	"github.com/midbel/svg"
)

type Hierarchy struct {
	Label string      `json:"name"`
	Value float64     `json:"value"`
	Sub   []Hierarchy `json:"children"`

	svg.Fill
}

func (h Hierarchy) GetValue() float64 {
	if h.isLeaf() {
		return h.Value
	}
	return h.Sum()
}

func (h Hierarchy) Sum() float64 {
	if h.isLeaf() {
		return h.Value
	}
	var s float64
	for i := range h.Sub {
		s += h.Sub[i].Sum()
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
	TilingSquarify
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

	cs := c.RenderElement(series)
	cs.Render(ws)
}

func (c TreemapChart) RenderElement(series []Hierarchy) svg.Element {
	for len(series) == 1 {
		series = series[0].Sub
	}
	return c.renderElement(series)
}

func (c TreemapChart) renderElement(series []Hierarchy) svg.Element {
	c.checkDefault()
	var (
		dim  = svg.NewDim(c.Width, c.Height)
		cs   = svg.NewSVG(dim.Option())
		area = svg.NewGroup(svg.WithID("area"), whitstrok.Option(), c.translate())
	)
	switch c.Tiling {
	case TilingDefault:
		c.drawDefault(&area, series, c.GetAreaWidth(), c.GetAreaHeight())
	case TilingBinary:
		c.drawBinary(&area, series)
	case TilingVertical:
		c.drawVertical(&area, series, c.GetAreaHeight()/getSum(series))
	case TilingHorizontal:
		c.drawHorizontal(&area, series, c.GetAreaWidth()/getSum(series))
	case TilingAlternate:
		c.drawAlternate(&area, series)
	case TilingSquarify:
		c.drawSquarify(&area, series, c.GetAreaWidth(), c.GetAreaHeight())
	default:
	}
	cs.Append(area.AsElement())
	cs.Append(c.drawTitle())
	cs.Append(c.drawLegend())
	return cs.AsElement()
}

func (c TreemapChart) drawAlternate(a Appender, series []Hierarchy) {
	var (
		part = c.GetAreaWidth() / getSum(series)
		off  float64
	)
	for i := range series {
		var (
			grp   = svg.NewGroup(svg.WithTranslate(off, 0))
			width = series[i].GetValue() * part
		)
		series[i].Fill = getFill(i, series[i].Fill, series[i].Fill)
		if m := series[i].Depth() % 2; m == 1 {
			c.alternateHorizontal(&grp, series[i], width, c.GetAreaHeight())
		} else {
			c.alternateVertical(&grp, series[i], width, c.GetAreaHeight())
		}
		off += series[i].GetValue() * part
		a.Append(grp.AsElement())
	}
}

func (c TreemapChart) alternateHorizontal(a Appender, serie Hierarchy, width, height float64) {
	var (
		sum   = serie.Sum()
		wpart = width / sum
		off   float64
	)
	for i := range serie.Sub {
		serie.Sub[i].Fill = getFill(i, serie.Sub[i].Fill, serie.Fill)

		sub := serie.Sub[i].GetValue() * wpart
		if serie.Sub[i].isLeaf() {
			var (
				dim = svg.NewDim(sub, height)
				rec = svg.NewRect(svg.WithPosition(off, 0), dim.Option(), serie.Sub[i].Fill.Option())
			)
			a.Append(rec.AsElement())
		} else {
			grp := svg.NewGroup(svg.WithTranslate(off, 0), svg.WithClass("horizontal"))
			if m := serie.Sub[i].Depth() % 2; m == 1 {
				c.alternateHorizontal(&grp, serie.Sub[i], sub, height)
			} else {
				c.alternateVertical(&grp, serie.Sub[i], sub, height)
			}
			a.Append(grp.AsElement())
		}
		off += sub
	}
}

func (c TreemapChart) alternateVertical(a Appender, serie Hierarchy, width, height float64) {
	var (
		sum   = serie.Sum()
		hpart = height / sum
		off   float64
	)
	for i := range serie.Sub {
		serie.Sub[i].Fill = getFill(i, serie.Sub[i].Fill, serie.Fill)

		sub := serie.Sub[i].GetValue() * hpart
		if serie.Sub[i].isLeaf() {
			var (
				dim = svg.NewDim(width, sub)
				rec = svg.NewRect(svg.WithPosition(0, off), dim.Option(), serie.Sub[i].Fill.Option())
			)
			a.Append(rec.AsElement())
		} else {
			grp := svg.NewGroup(svg.WithTranslate(0, off), svg.WithClass("vertical"))
			if m := serie.Sub[i].Depth() % 2; m == 1 {
				c.alternateHorizontal(&grp, serie.Sub[i], width, sub)
			} else {
				c.alternateVertical(&grp, serie.Sub[i], width, sub)
			}
			a.Append(grp.AsElement())
		}
		off += sub
	}
}

func (c TreemapChart) drawHorizontal(a Appender, series []Hierarchy, part float64) {
	var off float64
	for i := range series {
		series[i].Fill = getFill(i, series[i].Fill, series[i].Fill)
		if series[i].isLeaf() {
			var (
				pos = svg.NewPos(off, 0)
				dim = svg.NewDim(series[i].GetValue()*part, c.GetAreaHeight())
			)
			r := svg.NewRect(pos.Option(), dim.Option(), series[i].Fill.Option())
			r.Title = series[i].Label
			a.Append(r.AsElement())
		} else {
			grp := svg.NewGroup(svg.WithTranslate(off, 0))
			c.drawHorizontal(&grp, series[i].Sub, (series[i].GetValue()*part)/getSum(series[i].Sub))
			a.Append(grp.AsElement())
		}
		off += series[i].GetValue() * part
	}
}

func (c TreemapChart) drawVertical(a Appender, series []Hierarchy, part float64) {
	var off float64
	for i := range series {
		series[i].Fill = getFill(i, series[i].Fill, series[i].Fill)
		if series[i].isLeaf() {
			var (
				pos = svg.NewPos(0, off)
				dim = svg.NewDim(c.GetAreaWidth(), series[i].GetValue()*part)
			)
			r := svg.NewRect(pos.Option(), dim.Option(), series[i].Fill.Option())
			r.Title = series[i].Label
			a.Append(r.AsElement())
		} else {
			grp := svg.NewGroup(svg.WithTranslate(0, off))
			c.drawVertical(&grp, series[i].Sub, (series[i].GetValue()*part)/getSum(series[i].Sub))
			a.Append(grp.AsElement())
		}
		off += series[i].GetValue() * part
	}
}

func (c TreemapChart) drawBinary(a Appender, series []Hierarchy) {

}

func (c TreemapChart) drawDefault(a Appender, series []Hierarchy, width, height float64) {
	sort.Slice(series, func(i, j int) bool {
		return series[i].GetValue() > series[j].GetValue()
	})
	var (
		sum  = getSum(series)
		area = width * height
		offx float64
		offy float64
	)
	for i := range series {
		var (
			s  = series[i].GetValue() / sum
			w  float64
			h  float64
			ox float64
			oy float64
		)
		series[i].Fill = getFill(i, series[i].Fill, series[i].Fill)
		if math.Max(width, height) == width {
			h = height
			w = (s * area) / height
			ox = w
		} else {
			w = width
			h = (s * area) / width
			oy = h
		}
		if series[i].isLeaf() {
			var (
				dim  = svg.NewDim(w, h)
				rect = svg.NewRect(dim.Option(), svg.WithTranslate(offx, offy), series[i].Fill.Option())
			)
			rect.Title = series[i].Label
			a.Append(rect.AsElement())
		} else {
			grp := svg.NewGroup(svg.WithTranslate(offx, offy))
			a.Append(grp.AsElement())
			c.drawDefault(&grp, series[i].Sub, w, h)
		}
		width -= ox
		height -= oy
		offx += ox
		offy += oy
	}
}

func (c TreemapChart) drawSquarify(a Appender, series []Hierarchy, width, height float64) {
	c.squarify(a, series, width, height, 0)
}

const phi = 1.618

func (c TreemapChart) squarify(a Appender, series []Hierarchy, width, height float64, level int) {
	sort.Slice(series, func(i, j int) bool {
		return series[i].GetValue() > series[j].GetValue()
	})
	var (
		area   = width * height
		total  = getSum(series)
		full   = total
		ox, oy float64
		i, n   int
	)
	for i < len(series) {
		var (
			alpha = getAlpha(width, height, total)
			curr  = series[i].GetValue()
			max   = curr
			min   = curr
			sum   = curr
			ratio float64
			last  = getAspectRatio(alpha, min, max, sum)
			j     int
		)
		n++
		series[i].Fill = getFill(n, series[i].Fill, series[i].Fill)
		for j = i + 1; j < len(series); j++ {
			curr = series[j].GetValue()
			sum += curr
			min = math.Min(min, curr)
			max = math.Max(max, curr)
			ratio = getAspectRatio(alpha, min, max, sum)
			if ratio > last {
				sum -= curr
				break
			}
			last = ratio
			series[j].Fill = getFill(j, series[j].Fill, series[i].Fill)
		}
		var w, h float64
		if surface := area * (sum / full); width < height {
			w = width
			h = surface / width
			height -= h
		} else {
			h = height
			w = surface / height
			width -= w
		}
		parent := svg.NewGroup(svg.WithTranslate(ox, oy), svg.WithClass("container"))
		a.Append(parent.AsElement())
		c.layout(&parent, series[i:j], w, h, sum, level+1)
		if w == width {
			oy += h
		} else {
			ox += w
		}
		i = j
		total -= sum
	}
}

func (c TreemapChart) layout(a Appender, series []Hierarchy, width, height, sum float64, level int) {
	var ox, oy float64
	for i := range series {
		var (
			curr = series[i].GetValue()
			used = curr / sum
			w, h float64
		)
		if width > height {
			w = width * used
			h = height
		} else {
			h = height * used
			w = width
		}
		s, ok := getSerie(series[i])
		if !ok {
			grp := svg.NewGroup(svg.WithTranslate(ox, oy), svg.WithID(s.Label), svg.WithClass("row"))
			a.Append(grp.AsElement())
			for j := range s.Sub {
				s.Sub[j].Fill = getFill(j, s.Sub[j].Fill, series[i].Fill)
			}
			c.squarify(&grp, s.Sub, w, h, level)
		} else {
			var (
				p = svg.NewPos(ox, oy)
				d = svg.NewDim(w, h)
				r = svg.NewRect(p.Option(), d.Option(), series[i].Fill.Option())
			)
			r.Title = fmt.Sprintf("%s: %.0f", s.Label, curr)
			a.Append(r.AsElement())
		}
		if width > height {
			ox += w
		} else {
			oy += h
		}
	}
}

func getAlpha(width, height, sum float64) float64 {
	return math.Max(width/height, height/width) / (sum * phi)
}

func getAspectRatio(alpha, min, max, sum float64) float64 {
	beta := (sum * sum) * alpha
	return math.Max(max/beta, beta/min)
}

func getSerie(s Hierarchy) (Hierarchy, bool) {
	if s.isLeaf() {
		return s, true
	}
	return s, false
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
		sum += series[i].GetValue()
	}
	if sum == 0 {
		return 1
	}
	return sum
}
