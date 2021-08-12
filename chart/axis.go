package chart

import (
	"math"

	"github.com/midbel/svg"
)

type CategoryAxis struct {
	LabelX  bool
	DomainX bool

	InnerTicksY int
	OuterTicksY int
	LabelY      bool
	DomainY     bool
}

func NewCategoryAxisWith(ticks int, label, domain bool) CategoryAxis {
	return CategoryAxis{
		InnerTicksY: ticks,
		OuterTicksY: ticks,
		DomainX:     domain,
		DomainY:     domain,
		LabelX:      label,
		LabelY:      label,
	}
}

func (a CategoryAxis) drawAxis(c Chart, rg pair, domains []string) svg.Element {
	grp := svg.NewGroup()
	grp.Append(a.drawAxisX(c, domains))
	grp.Append(a.drawAxisY(c, rg))
	grp.Append(a.drawTicksY(c, rg))
	return grp.AsElement()
}

func (a CategoryAxis) drawAxisX(c Chart, domains []string) svg.Element {
	var (
		axis = svg.NewGroup(c.getOptionsAxisX()...)
		step = c.GetAreaWidth() / float64(len(domains))
	)
	if a.DomainX {
		var (
			pos1 = svg.NewPos(0, 0)
			pos2 = svg.NewPos(c.GetAreaWidth(), 0)
			line = svg.NewLine(pos1, pos2, axisstrok.Option(), svg.WithClass("domain"))
		)
		axis.Append(line.AsElement())
	}
	for i := 0; i < len(domains); i++ {
		var (
			grp = svg.NewGroup(svg.WithClass("tick"))
			off = float64(i) * step
		)
		if a.LabelX {
			var (
				pos0 = svg.NewPos(off+(step/3), textick)
				text = svg.NewText(domains[i], pos0.Option())
			)
			grp.Append(text.AsElement())
		}
		if a.DomainX {
			var (
				pos1 = svg.NewPos(off+step, 0)
				pos2 = svg.NewPos(off+step, ticklen)
				line = svg.NewLine(pos1, pos2, axisstrok.Option())
			)
			grp.Append(line.AsElement())
		}
		axis.Append(grp.AsElement())
	}
	return axis.AsElement()
}

func (a CategoryAxis) drawAxisY(c Chart, rg pair) svg.Element {
	var (
		axis  = svg.NewGroup(c.getOptionsAxisY()...)
		coeff = c.GetAreaHeight() / float64(a.InnerTicksY)
		step  = rg.Diff() / float64(a.InnerTicksY)
	)
	if a.DomainY {
		var (
			pos1 = svg.NewPos(0, 0)
			pos2 = svg.NewPos(0, c.GetAreaHeight()+1)
			line = svg.NewLine(pos1, pos2, axisstrok.Option(), svg.WithClass("domain"))
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

func (a CategoryAxis) drawTicksY(c Chart, rg pair) svg.Element {
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

type LineAxis struct {
	InnerTicksX int
	OuterTicksY int
	LabelX      bool
	DomainX     bool

	InnerTicksY int
	OuterTicksX int
	LabelY      bool
	DomainY     bool
}

func NewLineAxisWithTicks(x, y int) LineAxis {
	a := NewLineAxisWith(0, true, true)
	a.InnerTicksX = x
	a.InnerTicksY = y
	a.OuterTicksX = x
	a.OuterTicksY = y
	return a
}

func NewLineAxisWith(ticks int, label, domain bool) LineAxis {
	return LineAxis{
		InnerTicksX: ticks,
		InnerTicksY: ticks,
		OuterTicksX: ticks,
		OuterTicksY: ticks,
		DomainX:     domain,
		DomainY:     domain,
		LabelX:      label,
		LabelY:      label,
	}
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
			pos1 = svg.NewPos(0, 0)
			pos2 = svg.NewPos(c.GetAreaWidth(), 0)
			line = svg.NewLine(pos1, pos2, axisstrok.Option(), svg.WithClass("domain"))
		)
		axis.Append(line.AsElement())
	}
	for i, j := rg.Min, 0; i < rg.Max+step; i, j = i+step, j+1 {
		var (
			grp = svg.NewGroup(svg.WithClass("tick"))
			off = float64(j) * coeff
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
			pos1 = svg.NewPos(0, 0)
			pos2 = svg.NewPos(0, c.GetAreaHeight()+1)
			line = svg.NewLine(pos1, pos2, axisstrok.Option(), svg.WithClass("domain"))
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
