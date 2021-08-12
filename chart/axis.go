package chart

import (
  "math"

  "github.com/midbel/svg"
)

type LineAxis struct {
  InnerTicksX int
  InnerTicksY int
  OuterTicksX int
  OuterTicksY int

  LabelX bool
  LabelY bool
  DomainX bool
  DomainY bool
}

func (a LineAxis) drawAxis(c Chart, rx, ry pair) svg.Element {
  grp := svg.NewGroup()
  grp.Append(a.drawAxisX(c, rx))
  grp.Append(a.drawAxisY(c, ry))
  grp.Append(a.drawTicksY(c, ry))
  return grp.AsElement()
}

func (a LineAxis) drawAxisX(c Chart, rg pair) svg.Element {
	var (
		axis  = svg.NewGroup(c.getOptionsAxisX()...)
		coeff = c.GetAreaWidth() / float64(a.InnerTicksX)
		step  = math.Ceil(rg.Diff() / float64(a.InnerTicksX))
	)
  if a.DomainX {
    var (
      pos1  = svg.NewPos(0, 0)
      pos2  = svg.NewPos(c.GetAreaWidth(), 0)
      line  = svg.NewLine(pos1, pos2, axisstrok.Option(), svg.WithClass("domain"))
    )
    axis.Append(line.AsElement())
  }
	for i, j := rg.Min, 0; i < rg.Max+step; i, j = i+step, j+1 {
		var (
      grp  = svg.NewGroup(svg.WithClass("tick"))
			off  = float64(j) * coeff
		)
    if a.LabelX {
      var (
        pos0 = svg.NewPos(off-(coeff*0.2), textick+(textick/3))
        text = svg.NewText(formatFloat(i), pos0.Option())
      )
      grp.Append(text.AsElement())
    }
    if a.DomainX {
      var (
        pos1 = svg.NewPos(off, 0)
        pos2 = svg.NewPos(off, ticklen)
        line = svg.NewLine(pos1, pos2, axisstrok.Option())
      )
      grp.Append(line.AsElement())
    }
		axis.Append(grp.AsElement())
	}
	return axis.AsElement()
}

func (a LineAxis) drawAxisY(c Chart, rg pair) svg.Element {
	var (
		axis  = svg.NewGroup(c.getOptionsAxisY()...)
		coeff = c.GetAreaHeight() / float64(a.InnerTicksY)
		step  = rg.Diff() / float64(a.InnerTicksY)
	)
  if a.DomainY {
    var (
      pos1  = svg.NewPos(0, 0)
      pos2  = svg.NewPos(0, c.GetAreaHeight()+1)
      line  = svg.NewLine(pos1, pos2, axisstrok.Option(), svg.WithClass("domain"))
    )
    axis.Append(line.AsElement())
  }
	for i, j := rg.Min, 0; i < rg.Max+step; i, j = i+step, j+1 {
		var (
			grp  = svg.NewGroup(svg.WithClass("tick"))
			ypos = c.GetAreaHeight() - (float64(j) * coeff)
		)
    if a.LabelY {
      var (
        pos  = svg.NewPos(0, ypos+(ticklen/2))
        anc  = svg.WithAnchor("end")
        text = svg.NewText(formatFloat(i), anc, pos.Option())
      )
      text.Shift = svg.NewPos(-ticklen*2, 0)
      grp.Append(text.AsElement())
    }
    if a.DomainY {
      var (
        pos1 = svg.NewPos(-ticklen, ypos)
  			pos2 = svg.NewPos(0, ypos)
  			line = svg.NewLine(pos1, pos2, axisstrok.Option())
      )
      grp.Append(line.AsElement())
    }
		axis.Append(grp.AsElement())
	}
	return axis.AsElement()
}

func (a LineAxis) drawTicksY(c Chart, rg pair) svg.Element {
	var (
		max   = rg.AbsMax()
		grp   = svg.NewGroup(svg.WithClass("ticks", "ticks-y"), c.translate())
		step  = c.GetAreaHeight() / max
		coeff = max / float64(a.OuterTicksY)
	)
	for i := a.OuterTicksY; i > 0; i-- {
		var (
			ypos = c.GetAreaHeight() - (float64(i) * coeff * step)
			pos1 = svg.NewPos(0, ypos)
			pos2 = svg.NewPos(c.GetAreaWidth(), ypos)
		)
		grp.Append(getTick(pos1, pos2))
	}
	return grp.AsElement()
}
