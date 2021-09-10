package chart

import (
	"bufio"
	"fmt"
	"io"

	"github.com/midbel/svg"
)

type PolarSerie struct {
	xyserie
	svg.Fill
}

func NewPolarSerie(title string) PolarSerie {
	s := xyserie{Title: title}
	return PolarSerie{xyserie: s}
}

type PolarChart struct {
	Chart
	Radius float64
	Zone   float64
}

func (c PolarChart) Render(w io.Writer, series []PolarSerie) {
	ws := bufio.NewWriter(w)
	defer ws.Flush()

	cs := c.RenderElement(series)
	cs.Render(ws)
}

func (c PolarChart) RenderElement(series []PolarSerie) svg.Element {
	c.checkDefault()

	var (
		cs     = c.getCanvas()
		area   = c.getArea()
		space  = (c.Radius / c.Zone) / 2
		cx, cy = c.GetAreaCenter()
	)
	for i := 0; i < int(c.Zone); i++ {
		var (
			t = svg.WithTranslate(cx, cy)
			r = svg.WithRadius(space * float64(i+1))
			z = svg.NewCircle(r, t, nonefill.Option(), axisstrok.Option())
		)
		area.Append(z.AsElement())
	}

	var (
		hx, hy = c.GetAreaWidth() / 2, c.GetAreaHeight() / 2
		long   = getLongestSerie(series)
	)
	_ = fmt.Sprintf
	for i := 0; i < long; i++ {
		var (
			a = float64(i) * (90 / float64(long))
			g = svg.NewGroup(svg.WithRotate(a, cx, cy))
			x = svg.NewLine(svg.NewPos(0, hy), svg.NewPos(c.GetAreaWidth(), hy), axisstrok.Option())
			y = svg.NewLine(svg.NewPos(hx, 0), svg.NewPos(hx, c.GetAreaWidth()), axisstrok.Option())
		)
		g.Append(x.AsElement())
		g.Append(y.AsElement())
		area.Append(g.AsElement())
	}

	cs.Append(area.AsElement())
	return cs.AsElement()
}

func (c *PolarChart) checkDefault() {
	c.Chart.checkDefault()
	if c.Zone <= 0 {
		c.Zone = 5
	}
	if c.Radius <= 0 {
		c.Radius = getGreater(c.GetAreaWidth(), c.GetAreaHeight())
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
