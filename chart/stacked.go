package chart

import (
  "bufio"
  "io"
  "math"
  "math/rand"

  "github.com/midbel/svg"
)

type Serie struct {
	Title  string
	Labels []string
	Values []float64
}

func NewSerie(title string) Serie {
	return Serie{
		Title: title,
	}
}

func (s *Serie) Add(label string, value float64) {
	s.Labels = append(s.Labels, label)
	s.Values = append(s.Values, value)
}

func (s Serie) Sum() float64 {
	var sum float64
	for i := range s.Values {
		sum += s.Values[i]
	}
	return sum
}

func (s Serie) Len() int {
	return len(s.Values)
}

type StackedSerie struct {
	Title  string
	Series []Serie
	max    float64
}

func NewStackedSerie(title string) StackedSerie {
	return StackedSerie{
		Title: title,
		max:   math.NaN(),
	}
}

func (sr *StackedSerie) Append(s Serie) {
	if sum := s.Sum(); math.IsNaN(sr.max) || sum > sr.max {
		sr.max = sum
	}
	sr.Series = append(sr.Series, s)
}

func (sr *StackedSerie) Len() int {
	return len(sr.Series)
}

type StackedChart struct {
  Chart

  BarWidth float64
}

func (c StackedChart) Render(w io.Writer, series []StackedSerie) {
  ws := bufio.NewWriter(w)
  defer ws.Flush()
  c.render(ws, series)
}

func (c StackedChart) render(w svg.Writer, series []StackedSerie) {
  c.checkDefault()

  var (
		dim    = svg.NewDim(c.Width, c.Height)
		cs     = svg.NewSVG(dim.Option())
		area   = svg.NewGroup(svg.WithID("area"), c.translate())
		max, _ = getStackedDomains(series)
	)

	offset := c.GetAreaWidth() / float64(len(series))
	for i := range series {
		var (
			off  = offset * float64(i)
			grp  = svg.NewGroup(svg.WithTranslate(off, 0))
			elem = c.drawSerie(series[i], offset, max)
		)
		grp.Append(elem)
		area.Append(grp.AsElement())
	}
	cs.Append(area.AsElement())
	cs.Render(w)
}

func (c StackedChart) drawSerie(s StackedSerie, band, max float64) svg.Element {
	var (
		size   = band / float64(s.Len())
		width  = size * 0.8
		height = c.GetAreaHeight() / max
		grp    = svg.NewGroup(svg.WithTranslate((size*0.2)*2, 0))
	)
	if c.BarWidth > 0 {
		width = c.BarWidth
	}
	for j, s := range s.Series {
		var (
			g  = svg.NewGroup()
			rw = width
			ro = c.Height - c.Top
		)
		for i := range s.Labels {
			var (
				rh = s.Values[i] * height
				rx = (float64(j) * width) + ((width / 2) - (rw / 2))
				ry float64
			)
			ro -= rh
			ry = ro
			r := makeRect(svg.WithPosition(rx, ry), svg.WithDimension(rw, rh))
			r.Title = s.Labels[i]
			g.Append(r.AsElement())
		}
		grp.Append(g.AsElement())
	}
	return grp.AsElement()
}

func getStackedDomains(cs []StackedSerie) (float64, []string) {
	var (
		ys []string
		xs float64
	)
	for i := range cs {
		ys = append(ys, cs[i].Title)
		if i == 0 || cs[i].max > xs {
			xs = cs[i].max
		}
	}
	return xs, ys
}

func makeRect(options ...svg.Option) svg.Rect {
	fill := svg.NewFill(randomColor())
	options = append(options, fill.Option())
	return svg.NewRect(options...)
}

func randomColor() string {
	return svg.Colors[rand.Intn(len(svg.Colors))]
}
