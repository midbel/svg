package chart

import (
	"bufio"
	"io"
	"math"
	"sort"

	"github.com/midbel/svg"
	"github.com/midbel/svg/colors"
)

type Hierarchy struct {
	Label string      `json:"name"`
	Value float64     `json:"value"`
	Sub   []Hierarchy `json:"children"`
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
	return cs.AsElement()
}

func (c TreemapChart) drawAlternate(a appender, series []Hierarchy) {
	var (
		part = c.GetAreaWidth() / getSum(series)
		off  float64
	)
	for i := range series {
		var (
			color = colors.RdYlBu11[i%len(colors.RdYlBu11)]
			fill  = svg.NewFill(color)
			grp   = svg.NewGroup(svg.WithTranslate(off, 0), fill.Option())
			width = series[i].GetValue() * part
		)
		if m := series[i].Depth() % 2; m == 1 {
			c.alternateHorizontal(&grp, series[i], width, c.GetAreaHeight())
		} else {
			c.alternateVertical(&grp, series[i], width, c.GetAreaHeight())
		}
		off += series[i].GetValue() * part
		a.Append(grp.AsElement())
	}
}

func (c TreemapChart) alternateHorizontal(a appender, serie Hierarchy, width, height float64) {
	var (
		sum   = serie.Sum()
		wpart = width / sum
		off   float64
	)
	for i := range serie.Sub {
		sub := serie.Sub[i].GetValue() * wpart
		if serie.Sub[i].isLeaf() {
			var (
				dim = svg.NewDim(sub, height)
				rec = svg.NewRect(svg.WithPosition(off, 0), dim.Option())
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

func (c TreemapChart) alternateVertical(a appender, serie Hierarchy, width, height float64) {
	var (
		sum   = serie.Sum()
		hpart = height / sum
		off   float64
	)
	for i := range serie.Sub {
		sub := serie.Sub[i].GetValue() * hpart
		if serie.Sub[i].isLeaf() {
			var (
				dim = svg.NewDim(width, sub)
				rec = svg.NewRect(svg.WithPosition(0, off), dim.Option())
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

func (c TreemapChart) drawHorizontal(a appender, series []Hierarchy, part float64) {
	var off float64
	for i := range series {
		if series[i].isLeaf() {
			var (
				pos  = svg.NewPos(off, 0)
				dim  = svg.NewDim(series[i].GetValue()*part, c.GetAreaHeight())
				fill = svg.NewFill("steelblue")
			)
			r := svg.NewRect(pos.Option(), dim.Option(), fill.Option())
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

func (c TreemapChart) drawVertical(a appender, series []Hierarchy, part float64) {
	var off float64
	for i := range series {
		if series[i].isLeaf() {
			var (
				pos  = svg.NewPos(0, off)
				dim  = svg.NewDim(c.GetAreaWidth(), series[i].GetValue()*part)
				fill = svg.NewFill("steelblue")
			)
			r := svg.NewRect(pos.Option(), dim.Option(), fill.Option())
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

func (c TreemapChart) drawBinary(a appender, series []Hierarchy) {

}

func (c TreemapChart) drawDefault(a appender, series []Hierarchy, width, height float64) {
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
				fill = svg.NewFill("steelblue")
				rect = svg.NewRect(dim.Option(), fill.Option(), svg.WithTranslate(offx, offy))
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

func (c TreemapChart) drawSquarify(a appender, series []Hierarchy, width, height float64) {
	sort.Slice(series, func(i, j int) bool {
		return series[i].GetValue() > series[j].GetValue()
	})
	var (
		i    int
		offx float64
		offy float64
		sum  = getSum(series)
		area = width * height
	)
	for i < len(series) {
		var (
			val   float64
			ratio float64
			prev  = math.NaN()
			min   = prev
			max   = prev
			short = width
			j     = i
		)
		if width > height {
			short = height
		}
		for j < len(series) {
			var (
				ssq  = math.Pow(short, 2)
				vsq  = math.Pow(val, 2)
				r1   = (ssq * max) / vsq
				r2   = vsq / (ssq * min)
				curr float64
			)
			ratio = math.Max(r1, r2)
			if !math.IsNaN(prev) && ratio >= prev {
				break
			}
			prev = ratio

			curr = series[j].GetValue()
			val += curr
			if math.IsNaN(min) || curr < min {
				min = curr
			}
			if math.IsNaN(max) || curr > max {
				max = curr
			}
			j++
		}
		var w, h float64
		if used := val / sum; short == width {
			w = short
			h = (area * used) / short
			height -= h
		} else {
			w = (area * used) / short
			h = short
			width -= w
		}
		parent := svg.NewGroup(svg.WithTranslate(offx, offy))
		a.Append(parent.AsElement())
		var (
			tw, th float64
			ox, oy float64
		)
		if short == w {
			tw = w
			th = h / float64(j-i)
			oy = th
		} else {
			th = h
			tw = w / float64(j-i)
			ox = tw
		}
		for k := 0; i < j; k++ {
			if series[i].isLeaf() {
				var (
					p = svg.NewPos(float64(k)*ox, float64(k)*oy)
					f = svg.NewFill("steelblue")
					d = svg.NewDim(tw, th)
					r = svg.NewRect(p.Option(), f.Option(), d.Option())
				)
				r.Title = series[i].Label
				parent.Append(r.AsElement())
			} else {
				grp := svg.NewGroup(svg.WithTranslate(float64(k)*ox, float64(k)*oy))
				c.drawSquarify(&grp, series[i].Sub, tw, th)
				parent.Append(grp.AsElement())
			}
			i++
		}
		if short == w {
			offy += h
		} else {
			offx += w
		}
		i = j
	}
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
