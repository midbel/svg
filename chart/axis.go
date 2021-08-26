package chart

import (
	"math"
	"time"

	"github.com/midbel/svg"
)

type CategoryAxis struct {
	InnerX  bool
	OuterX  bool
	LabelX  bool
	DomainX bool

	TicksY int
	InnerY bool
	OuterY bool
	LabelY      bool
	DomainY     bool
}

func NewCategoryAxis(ticks int, label, domain bool) CategoryAxis {
	return CategoryAxis{
		TicksY: ticks,
		InnerY: true,
		InnerX: true,
		OuterY: false,
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
	grp.Append(a.drawDomains(c))
	return grp.AsElement()
}

func (a CategoryAxis) drawDomains(c Chart) svg.Element {
	grp := svg.NewGroup(c.translate())
	if a.DomainX {
		var (
			pos1 = svg.NewPos(0, c.GetAreaHeight())
			pos2 = svg.NewPos(c.GetAreaWidth(), c.GetAreaHeight())
			line = svg.NewLine(pos1, pos2, axisstrok.Option(), svg.WithClass("domain"))
		)
		grp.Append(line.AsElement())
	}
	if a.DomainY {
		var (
			pos1 = svg.NewPos(0, 0)
			pos2 = svg.NewPos(0, c.GetAreaHeight()+1)
			line = svg.NewLine(pos1, pos2, axisstrok.Option(), svg.WithClass("domain"))
		)
		grp.Append(line.AsElement())
	}
	return grp.AsElement()
}

func (a CategoryAxis) drawAxisX(c Chart, domains []string) svg.Element {
	var (
		axis = svg.NewGroup(c.getOptionsAxisX()...)
		step = c.GetAreaWidth() / float64(len(domains))
	)
	for i := 0; i < len(domains); i++ {
		var (
			grp = svg.NewGroup(svg.WithClass("tick"))
			off = float64(i) * step
		)
		if a.LabelX {
			var (
				font = svg.NewFont(12)
				pos0 = svg.NewPos(off+(step/3), textick)
				text = svg.NewText(domains[i], pos0.Option(), font.Option())
			)
			grp.Append(text.AsElement())
		}
		if a.InnerX {
			var (
				pos1 = svg.NewPos(off+step, 0)
				pos2 = svg.NewPos(off+step, ticklen)
				line = svg.NewLine(pos1, pos2, axisstrok.Option())
			)
			grp.Append(line.AsElement())
		}
		if a.OuterX {
			var (
				pos1 = svg.NewPos(off+step, 0)
				pos2 = svg.NewPos(off+step, -c.GetAreaHeight())
			)
			grp.Append(getTick(pos1, pos2))
		}
		axis.Append(grp.AsElement())
	}
	return axis.AsElement()
}

func (a CategoryAxis) drawAxisY(c Chart, rg pair) svg.Element {
	var (
		axis  = svg.NewGroup(c.getOptionsAxisY()...)
		coeff = c.GetAreaHeight() / float64(a.TicksY)
		step  = rg.Diff() / float64(a.TicksY)
	)
	for i, j := rg.Min, 0; i < rg.Max+step; i, j = i+step, j+1 {
		var (
			grp  = svg.NewGroup(svg.WithClass("tick"))
			ypos = c.GetAreaHeight() - (float64(j) * coeff)
		)
		if a.LabelY {
			var (
				pos  = svg.NewPos(0, ypos+(ticklen/2))
				anc  = svg.WithAnchor("end")
				font = svg.NewFont(12)
				text = svg.NewText(formatFloat(i), anc, pos.Option(), font.Option())
			)
			text.Shift = svg.NewPos(-ticklen*2, 0)
			grp.Append(text.AsElement())
		}
		if a.InnerY {
			var (
				pos1 = svg.NewPos(-ticklen, ypos)
				pos2 = svg.NewPos(0, ypos)
				line = svg.NewLine(pos1, pos2, axisstrok.Option())
			)
			grp.Append(line.AsElement())
		}
		if a.OuterY {
			var (
				pos1 = svg.NewPos(0, ypos)
				pos2 = svg.NewPos(c.GetAreaWidth(), ypos)
			)
			grp.Append(getTick(pos1, pos2))
		}
		axis.Append(grp.AsElement())
	}
	return axis.AsElement()
}

type LineAxis struct {
	TicksX  int
	OuterX  bool
	InnerX  bool
	LabelX  bool
	DomainX bool

	TicksY  int
	OuterY  bool
	InnerY  bool
	LabelY  bool
	DomainY bool
}

func NewLineAxisWithTicks(c int) LineAxis {
	return NewLineAxis(c, true, true)
}

func NewLineAxis(ticks int, label, domain bool) LineAxis {
	return LineAxis{
		TicksX:  ticks,
		InnerX:  true,
		OuterX:  true,
		TicksY:  ticks,
		InnerY:  true,
		OuterY:  true,
		DomainX: domain,
		DomainY: domain,
		LabelX:  label,
		LabelY:  label,
	}
}

func (a LineAxis) drawAxis(c Chart, rx, ry pair) svg.Element {
	grp := svg.NewGroup()
	grp.Append(a.drawAxisX(c, rx))
	grp.Append(a.drawAxisY(c, ry))
	grp.Append(a.drawDomains(c))
	return grp.AsElement()
}

func (a LineAxis) drawDomains(c Chart) svg.Element {
	grp := svg.NewGroup(c.translate())
	if a.DomainX {
		var (
			pos1 = svg.NewPos(0, c.GetAreaHeight())
			pos2 = svg.NewPos(c.GetAreaWidth(), c.GetAreaHeight())
			line = svg.NewLine(pos1, pos2, axisstrok.Option(), svg.WithClass("domain"))
		)
		grp.Append(line.AsElement())
	}
	if a.DomainY {
		var (
			pos1 = svg.NewPos(0, 0)
			pos2 = svg.NewPos(0, c.GetAreaHeight()+1)
			line = svg.NewLine(pos1, pos2, axisstrok.Option(), svg.WithClass("domain"))
		)
		grp.Append(line.AsElement())
	}
	return grp.AsElement()
}

func (a LineAxis) drawAxisX(c Chart, rg pair) svg.Element {
	var (
		axis  = svg.NewGroup(c.getOptionsAxisX()...)
		coeff = c.GetAreaWidth() / float64(a.TicksX)
		step  = math.Ceil(rg.Diff() / float64(a.TicksX))
	)
	for i, j := rg.Min, 0; i < rg.Max+step; i, j = i+step, j+1 {
		var (
			grp = svg.NewGroup(svg.WithClass("tick"))
			off = float64(j) * coeff
		)
		if a.LabelX {
			var (
				pos0 = svg.NewPos(off-(coeff*0.2), textick+(textick/3))
				font = svg.NewFont(12)
				text = svg.NewText(formatFloat(i), pos0.Option(), font.Option())
			)
			grp.Append(text.AsElement())
		}
		if a.InnerX {
			var (
				pos1 = svg.NewPos(off, 0)
				pos2 = svg.NewPos(off, ticklen)
				line = svg.NewLine(pos1, pos2, axisstrok.Option())
			)
			grp.Append(line.AsElement())
		}
		if a.OuterX {
			var (
				pos1 = svg.NewPos(off, 0)
				pos2 = svg.NewPos(off, -c.GetAreaHeight())
			)
			grp.Append(getTick(pos1, pos2))
		}
		axis.Append(grp.AsElement())
	}
	return axis.AsElement()
}

func (a LineAxis) drawAxisY(c Chart, rg pair) svg.Element {
	var (
		axis  = svg.NewGroup(c.getOptionsAxisY()...)
		coeff = c.GetAreaHeight() / float64(a.TicksY)
		step  = rg.Diff() / float64(a.TicksY)
	)
	for i, j := rg.Min, 0; i < rg.Max+step; i, j = i+step, j+1 {
		var (
			grp  = svg.NewGroup(svg.WithClass("tick"))
			ypos = c.GetAreaHeight() - (float64(j) * coeff)
		)
		if a.LabelY {
			var (
				pos  = svg.NewPos(0, ypos+(ticklen/2))
				anc  = svg.WithAnchor("end")
				font = svg.NewFont(12)
				text = svg.NewText(formatFloat(i), anc, pos.Option(), font.Option())
			)
			text.Shift = svg.NewPos(-ticklen*2, 0)
			grp.Append(text.AsElement())
		}
		if a.InnerY {
			var (
				pos1 = svg.NewPos(-ticklen, ypos)
				pos2 = svg.NewPos(0, ypos)
				line = svg.NewLine(pos1, pos2, axisstrok.Option())
			)
			grp.Append(line.AsElement())
		}
		if a.OuterY {
			var (
				pos1 = svg.NewPos(0, ypos)
				pos2 = svg.NewPos(c.GetAreaWidth(), ypos)
			)
			grp.Append(getTick(pos1, pos2))
		}
		axis.Append(grp.AsElement())
	}
	return axis.AsElement()
}

type TimeAxis struct {
	InnerTicksX int
	OuterTicksY int
	LabelX      bool
	DomainX     bool

	InnerTicksY int
	OuterTicksX int
	LabelY      bool
	DomainY     bool
}

func NewTimeAxisWithTicks(x, y int) TimeAxis {
	a := NewTimeAxisWith(0, true, true)
	a.InnerTicksX = x
	a.InnerTicksY = y
	a.OuterTicksX = x
	a.OuterTicksY = y
	return a
}

func NewTimeAxisWith(ticks int, label, domain bool) TimeAxis {
	return TimeAxis{
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

func (a TimeAxis) drawAxis(c Chart, rx timepair, ry pair) svg.Element {
	grp := svg.NewGroup()
	grp.Append(a.drawAxisX(c, rx))
	grp.Append(a.drawAxisY(c, ry))
	grp.Append(a.drawTicksY(c, ry))
	return grp.AsElement()
}

func (a TimeAxis) drawAxisX(c Chart, rg timepair) svg.Element {
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
	for j := 0; !rg.Min.After(rg.Max); j++ {
		var (
			grp = svg.NewGroup(svg.WithClass("tick"))
			off = float64(j) * coeff
		)
		if a.LabelX {
			var (
				font = svg.NewFont(12)
				pos0 = svg.NewPos(off-(coeff*0.15), textick+(textick/3))
				text = svg.NewText(formatTime(rg.Min), pos0.Option(), font.Option())
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
		rg.Min = rg.Min.Add(time.Second * time.Duration(step))
	}
	return axis.AsElement()
}

func (a TimeAxis) drawAxisY(c Chart, rg pair) svg.Element {
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
				font = svg.NewFont(12)
				text = svg.NewText(formatFloat(i), anc, pos.Option(), font.Option())
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

func (a TimeAxis) drawTicksY(c Chart, rg pair) svg.Element {
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

const (
	ticklen = 7
	textick = 18
)

func getTick(pos1, pos2 svg.Pos) svg.Element {
	tickstrok.Dash.Array = []int{5}
	line := svg.NewLine(pos1, pos2, tickstrok.Option())
	return line.AsElement()
}
