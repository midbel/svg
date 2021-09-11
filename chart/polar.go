package chart

import (
	"bufio"
	"io"

	"github.com/midbel/svg"
)

type PolarSerie struct {
	Title  string
	values []float64
	pair

	Radius float64
	svg.Fill
	svg.Stroke
}

func NewPolarSerie(title string) PolarSerie {
	return PolarSerie{Title: title}
}

func (s *PolarSerie) Add(v float64) {
	if len(s.values) == 0 {
		s.pair.Min = v
		s.pair.Max = v
	}
	s.pair.Min = getLesser(s.pair.Min, v)
	s.pair.Max = getGreater(s.pair.Max, v)
	s.values = append(s.values, v)
}

func (s *PolarSerie) Len() int {
	return len(s.values)
}

type PolarChart struct {
	Chart
	Ticks int
	Radius float64
	Zone   float64
}

func (c PolarChart) Render(w io.Writer, serie PolarSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()

	cs := c.RenderElement(serie)
	cs.Render(ws)
}

func (c PolarChart) RenderElement(serie PolarSerie) svg.Element {
	c.checkDefault()
	var (
		cs     = c.getCanvas()
		area   = c.getArea()
	)
	c.drawInnerCircles(&area)

	serie.Fill = getFill(0, serie.Fill, serie.Fill)
	var (
		cx, cy = c.GetAreaCenter()
		dx  = c.Radius / serie.pair.Diff()
		lg  = serie.Len()
		pat = svg.NewPath(serie.Fill.Stroke().Option())
	)
	if serie.Radius <= 0 {
		serie.Radius = 5
	}

	for i := 0; i < lg; i++ {
		var (
			g   = svg.NewGroup(svg.WithTranslate(cx, cy))
			a   = float64(i) * (halfcirc / float64(lg))
			p   = (serie.values[i] - serie.pair.Min) * dx
			r   = svg.WithRadius(serie.Radius)
			pos = getPosFromAngle(a, p)
			pt  = svg.NewCircle(r, serie.Fill.Option(), pos.Option())
		)
		if i == 0 {
			pat.AbsMoveTo(pos)
		} else {
			pat.AbsLineTo(pos)
		}
		g.Append(pt.AsElement())
		area.Append(g.AsElement())
	}
	pat.ClosePath()

	g := svg.NewGroup(svg.WithTranslate(cx, cy))
	g.Append(pat.AsElement())
	area.Append(g.AsElement())

	cs.Append(area.AsElement())
	return cs.AsElement()
}

func (c *PolarChart) drawInnerCircles(a appender) {
	var (
		space  = c.Radius / c.Zone
		cx, cy = c.GetAreaCenter()
	)
	for i := 0; i < int(c.Zone); i++ {
		var (
			t = svg.WithTranslate(cx, cy)
			r = svg.WithRadius(space * float64(i+1))
			z = svg.NewCircle(r, t, nonefill.Option(), axisstrok.Option())
		)
		a.Append(z.AsElement())
	}
	for i := 0; i < c.Ticks; i++ {
		var (
			g = svg.NewGroup(svg.WithTranslate(cx, cy))
			ang   = float64(i) * (halfcirc / float64(c.Ticks))
			pos1 = getPosFromAngle(ang*deg2rad, c.Radius)
			pos2 = getPosFromAngle((ang+halfcirc)*deg2rad, c.Radius)
			line = svg.NewLine(pos1, pos2, axisstrok.Option())
		)
		g.Append(line.AsElement())
		a.Append(g.AsElement())
	}
}

func (c *PolarChart) checkDefault() {
	c.Chart.checkDefault()
	if c.Zone <= 0 {
		c.Zone = 5
	}
	if c.Radius <= 0 {
		c.Radius = getGreater(c.GetAreaWidth(), c.GetAreaHeight()) / 2
	}
	if c.Ticks <= 0 {
		c.Ticks = 5
	}
}

func getLongestSerie(series []PolarSerie) int {
	var j int
	for i := range series {
		if n := series[i].Len(); i == 0 || n > j {
			j = n
		}
	}
	return j
}
