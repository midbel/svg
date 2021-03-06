package chart

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/midbel/svg"
)

var (
	tickstrok = svg.NewStroke("darkgray", 0.5)
	axisstrok = svg.NewStroke("darkgray", 0.5)
)

type Axis interface {
	Draw(Appender, float64, float64, ...svg.Option)
	Domain() Range
	update(options ...AxisOption)
	Left() bool
	Bottom() bool
	Horizontal() bool
	Vertical() bool
}

type AxisOption func(Axis)

type FormatterFunc func(interface{}, Position) string

func WithFormatter(fn FormatterFunc) AxisOption {
	return func(a Axis) {
		if fn == nil {
			return
		}
		switch a := a.(type) {
		case *labelAxis:
			a.Formatter = fn
		case *numberAxis:
			a.Formatter = fn
		case *timeAxis:
			a.Formatter = fn
		default:
		}
	}
}

func WithTimeRange(starts, ends time.Time) AxisOption {
	return func(a Axis) {
		switch a := a.(type) {
		case *timeAxis:
			a.starts, a.ends = starts, ends
		case *numberAxis:
			a.starts, a.ends = float64(starts.Unix()), float64(ends.Unix())
		case *labelAxis:
			a.labels = append(a.labels, starts.Format(time.RFC3339))
			a.labels = append(a.labels, ends.Format(time.RFC3339))
		default:
		}
	}
}

func WithNumberRange(starts, ends float64) AxisOption {
	return func(a Axis) {
		switch a := a.(type) {
		case *timeAxis:
			a.starts, a.ends = time.Unix(int64(starts), 0), time.Unix(int64(ends), 0)
		case *numberAxis:
			if len(a.domains) > 0 {
				return
			}
			a.starts, a.ends = starts, ends
		case *labelAxis:
			a.labels = append(a.labels, strconv.FormatFloat(starts, 'f', -1, 64))
			a.labels = append(a.labels, strconv.FormatFloat(ends, 'f', -1, 64))
		default:
		}
	}
}

func WithLabels(labels ...string) AxisOption {
	return func(a Axis) {
		if a, ok := a.(*labelAxis); ok {
			a.labels = append(a.labels[:0], labels...)
		}
	}
}

func WithTicks(n int, inner, outer bool) AxisOption {
	return func(a Axis) {
		switch a := a.(type) {
		case *timeAxis:
			a.WithTicks = n
			a.WithInner = inner
			a.WithOuter = outer
		case *numberAxis:
			a.WithTicks = n
			a.WithInner = inner
			a.WithOuter = outer
		case *labelAxis:
			a.WithTicks = n
			a.WithInner = inner
			a.WithOuter = outer
		default:
		}
	}
}

func WithPosition(p Position) AxisOption {
	return func(a Axis) {
		switch a := a.(type) {
		case *timeAxis:
			a.Position = p
		case *numberAxis:
			a.Position = p
		case *labelAxis:
			a.Position = p
		default:
		}
	}
}

const (
	ticklen = 7
	textick = 18
	numtick = 5
)

type AxisConfig struct {
	WithTitle  bool
	WithInner  bool
	WithOuter  bool
	WithLabel  bool
	WithDomain bool
	WithTicks  int
	Formatter  FormatterFunc
	Position

	svg.Font
	svg.Fill
	svg.Stroke
}

func (a AxisConfig) Left() bool {
	return a.Position == 0 || a.canLeft()
}

// func (a AxisConfig )Horizontal() bool {
// 	return a.Position.Horizontal()
// }

func (a AxisConfig) Bottom() bool {
	return a.Position == 0 || a.canBottom()
}

// func (a AxisConfig) Vertical() bool {
// 	return a.Position.Vertical()
// }

func (a AxisConfig) drawDomain(ap Appender, size float64) {
	if !a.WithDomain {
		return
	}
	var (
		pos1 = svg.NewPos(0, 0)
		pos2 svg.Pos
		line svg.Line
	)
	if a.Horizontal() {
		pos2.X, pos2.Y = size, 0
	} else {
		pos2.X, pos2.Y = 0, size
	}
	line = svg.NewLine(pos1, pos2, axisstrok.Option(), svg.WithClass("domain"))
	ap.Append(line.AsElement())
}

func (a AxisConfig) getInnerTick(off float64) svg.Element {
	if !a.WithInner {
		return nil
	}
	var (
		pos1 = a.adjustPosition(svg.NewPos(-ticklen, off))
		pos2 = svg.NewPos(0, off)
	)
	if a.Horizontal() {
		pos1.X, pos1.Y = off, 0
		pos2.X, pos2.Y = pos1.X, ticklen
	}
	line := svg.NewLine(pos1, pos2, axisstrok.Option())
	return line.AsElement()
}

func (a AxisConfig) getOuterTick(off, rsize float64) svg.Element {
	if !a.WithOuter {
		return nil
	}
	var (
		pos1 = svg.NewPos(0, off)
		pos2 = a.adjustPosition(svg.NewPos(rsize, off))
	)
	if a.Horizontal() {
		pos1.X, pos1.Y = pos1.Y, pos1.X
		pos2.X, pos2.Y = pos1.X, -rsize
		pos2 = a.adjustPosition(pos2)
	}
	tickstrok.Dash.Array = []int{5}
	line := svg.NewLine(pos1, pos2, tickstrok.Option())
	return line.AsElement()
}

func (a AxisConfig) getTickLabel(v interface{}, off float64) svg.Element {
	if !a.WithLabel {
		return nil
	}
	str := a.Formatter(v, a.Position)
	if str == "" {
		return nil
	}
	var (
		font    = svg.NewFont(12)
		shift   float64
		pos     svg.Pos
		options []svg.Option
	)
	switch a.Position {
	case Top, Bottom:
		pos.X, pos.Y = off, textick
		pos = a.adjustPosition(pos)
		options = append(options, svg.WithAnchor("middle"))
	case Left:
		pos.X, pos.Y = 0, off+(ticklen/2)
		shift = -ticklen * 2
		options = append(options, svg.WithAnchor("end"))
	case Right:
		pos.X, pos.Y = 0, off+(ticklen/2)
		shift = ticklen * 2
		options = append(options, svg.WithAnchor("start"))
	default:
	}
	options = append(options, pos.Option())
	options = append(options, font.Option())
	text := svg.NewText(str, options...)
	if shift != 0 {
		text.Shift = svg.NewPos(shift, 0)
	}
	return text.AsElement()
}

func (a AxisConfig) skip() bool {
	return a.WithTicks == 0 || (!a.WithInner && !a.WithOuter && !a.WithLabel)
}

func (a AxisConfig) adjustPosition(pos svg.Pos) svg.Pos {
	switch a.Position {
	case Top:
		pos.Y = -pos.Y
	case Bottom:
	case Left:
	case Right:
		pos.X = -pos.X
	default:
	}
	return pos
}

type labelAxis struct {
	AxisConfig
	labels []string
}

func CreateLabelAxis(options ...AxisOption) Axis {
	var a labelAxis
	a.Formatter = formatString
	a.WithDomain = true
	a.WithInner = true
	a.WithLabel = true
	a.WithTicks = numtick
	a.update(options...)
	return &a
}

func (a *labelAxis) Domain() Range {
	return nil
}

func (a *labelAxis) Draw(ap Appender, size, rsize float64, options ...svg.Option) {
	a.drawDomain(ap, size)
	if len(a.labels) == 0 && a.skip() {
		return
	}
	a.drawTicks(ap, size, rsize)
}

func (a *labelAxis) drawTicks(ap Appender, size, rsize float64) {
	step := size / float64(len(a.labels))
	for i := 0; i < len(a.labels); i++ {
		var (
			grp = svg.NewGroup(svg.WithClass("tick"))
			off = float64(i) * step
		)
		grp.Append(a.getTickLabel(a.labels[i], off+(step/2)))
		grp.Append(a.getInnerTick(off))
		grp.Append(a.getOuterTick(off, rsize))

		ap.Append(grp.AsElement())
	}
}

func (a *labelAxis) update(options ...AxisOption) {
	for _, o := range options {
		o(a)
	}
}

type timeAxis struct {
	AxisConfig

	starts  time.Time
	ends    time.Time
	domains []time.Time
}

func CreateTimeAxis(options ...AxisOption) Axis {
	var a timeAxis
	a.Formatter = formatTime
	a.WithDomain = true
	a.WithInner = true
	a.WithOuter = true
	a.WithLabel = true
	a.WithTicks = numtick
	a.update(options...)
	return &a
}

func (a *timeAxis) Domain() Range {
	return timepair{
		Min: a.starts,
		Max: a.ends,
	}
}

func (a *timeAxis) Draw(ap Appender, size, rsize float64, options ...svg.Option) {
	a.drawDomain(ap, size)
	if a.skip() {
		return
	}
	a.drawTicks(ap, size, rsize)
}

func (a *timeAxis) drawTicks(ap Appender, size, rsize float64) {
	var (
		coeff = size / float64(a.WithTicks)
		half  = coeff / 2
	)
	for j, w := range a.getDomains() {
		var (
			grp = svg.NewGroup(svg.WithClass("tick"))
			off = (float64(j) * coeff) + half
		)
		grp.Append(a.getInnerTick(off))
		grp.Append(a.getTickLabel(w, off))
		grp.Append(a.getOuterTick(off, rsize))
		ap.Append(grp.AsElement())
	}
}

func (a *timeAxis) getDomains() []time.Time {
	if len(a.domains) > 0 {
		return a.domains
	}
	var (
		starts  = a.starts
		ends    = a.ends
		diff    = ends.Sub(starts).Seconds()
		step    = math.Ceil(diff / float64(a.WithTicks))
		delta   = time.Duration(step) * time.Second
		domains []time.Time
	)
	for starts.Before(ends) {
		domains = append(domains, starts.Add(delta/2))
		starts = starts.Add(delta)
	}
	return domains
}

func (a *timeAxis) update(options ...AxisOption) {
	for _, o := range options {
		o(a)
	}
}

type numberAxis struct {
	AxisConfig

	starts  float64
	ends    float64
	domains []float64
}

func CreateNumberAxis(options ...AxisOption) Axis {
	var a numberAxis
	a.Formatter = formatFloat
	a.WithDomain = true
	a.WithInner = true
	a.WithOuter = true
	a.WithLabel = true
	a.WithTicks = numtick
	a.update(options...)

	if a.starts != 0 && a.ends != 0 && a.WithTicks > 0 {
		diff := a.ends - a.starts
		step := diff / float64(a.WithTicks)
		value := a.starts + step
		for i := 0; i < a.WithTicks; i++ {
			a.domains = append(a.domains, value)
			value += step
		}
	}

	return &a
}

func (a *numberAxis) Domain() Range {
	return pair{
		Min: a.starts,
		Max: a.ends,
	}
}

func (a *numberAxis) Draw(ap Appender, size, rsize float64, options ...svg.Option) {
	a.drawDomain(ap, size)
	if a.skip() {
		return
	}
	a.drawTicks(ap, size, rsize)
}

func (a *numberAxis) drawTicks(ap Appender, size, rsize float64) {
	var (
		coeff   = size / float64(a.WithTicks)
		half    = coeff / 2
		domains = a.getDomains()
		num     = len(domains) - 1
	)
	for j, v := range domains {
		var (
			grp = svg.NewGroup(svg.WithClass("tick"))
			off = size - (float64(num-j) * coeff) - half
		)
		if a.Position.Vertical() {
			off = size - (float64(j) * coeff) - half
		}
		grp.Append(a.getInnerTick(off))
		grp.Append(a.getTickLabel(v, off))
		grp.Append(a.getOuterTick(off, rsize))
		ap.Append(grp.AsElement())
	}
}

func (a *numberAxis) getDomains() []float64 {
	if len(a.domains) > 0 {
		return a.domains
	}
	var (
		domains []float64
		step    = math.Abs(a.ends-a.starts) / float64(a.WithTicks)
		starts  = a.starts
		ends    = a.ends - (step / 2)
	)
	for starts < ends {
		domains = append(domains, starts+(step/2))
		starts += step
	}
	return domains
}

func (a *numberAxis) update(options ...AxisOption) {
	for _, o := range options {
		o(a)
	}
}

func formatTime(v interface{}, _ Position) string {
	t, ok := v.(time.Time)
	if !ok {
		return fmt.Sprintf("%v", v)
	}
	// return t.Format("15:04:05")
	return t.Format("2006-01-02 15:04")
}

func formatFloat(v interface{}, _ Position) string {
	f, ok := v.(float64)
	if !ok {
		return fmt.Sprintf("%v", v)
	}
	if almostZero(f) {
		return "0.00"
	}
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func formatString(v interface{}, _ Position) string {
	s, ok := v.(string)
	if !ok {
		return fmt.Sprintf("%v", v)
	}
	return s
}

const threshold = 1e-9

func almostZero(val float64) bool {
	return math.Abs(val-0) <= threshold
}
