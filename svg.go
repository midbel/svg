package svg

import (
	"io"
	"strconv"
)

const (
	namespace = "http://www.w3.org/2000/svg"
	prolog    = "<?xml version=\"1.0\" encoding=\"utf-8\"?>"
)

const (
	defaultWidth    = 800
	defaultHeight   = 600
	defaultFontSize = 14
)

var DefaultStroke = NewStroke("black", 1)

type Writer interface {
	io.ByteWriter
	io.StringWriter
}

type Element interface {
	Render(Writer)
	setId(string)
	setClass([]string)
	setStyle(string, []string)
}

type node struct {
	Id     string
	Class  []string
	Styles map[string][]string
}

func (n *node) setId(id string) {
	n.Id = id
}

func (n *node) setClass(class []string) {
	n.Class = append(n.Class, class...)
}

func (n *node) setStyle(prop string, values []string) {
	if n.Styles == nil {
		n.Styles = make(map[string][]string)
	}
	n.Styles[prop] = append(n.Styles[prop], values...)
}

func (n *node) attrs() []string {
	var attrs []string
	if n.Id != "" {
		attrs = append(attrs, appendString("id", n.Id))
	}
	if len(n.Class) > 0 {
		attrs = append(attrs, appendStringArray("class", n.Class, space))
	}
	return attrs
}

type List struct {
	node
	Elems []Element
}

func NewList(es ...Element) List {
	elems := make([]Element, len(es))
	copy(elems, es)
	return List{
		Elems: elems,
	}
}

func (i *List) Append(e Element) {
	i.Elems = append(i.Elems, e)
}

func (i *List) Render(w Writer) {
	for _, e := range i.Elems {
		e.Render(w)
	}
}

type SVG struct {
	node
	List

	Dim
}

func NewSVG(options ...Option) SVG {
	var s SVG
	s.Dim = NewDim(defaultWidth, defaultHeight)
	for _, o := range options {
		o(&s)
	}
	return s
}

func (s *SVG) Render(w Writer) {
	attrs := s.node.attrs()
	attrs = append(attrs, appendString("xmlns", namespace))
	attrs = append(attrs, s.Dim.List()...)

	w.WriteString(prolog)
	writeElement(w, "svg", attrs, func() {
		s.List.Render(w)
	})
}

func (s *SVG) AsElement() Element {
	return s
}

type Group struct {
	node
	List

	Fill string
	Stroke
	Transform
}

func NewGroup(options ...Option) Group {
	var g Group
	for _, o := range options {
		o(&g)
	}
	return g
}

func (g *Group) Render(w Writer) {
	attrs := g.node.attrs()
	attrs = append(attrs, g.Stroke.List()...)
	attrs = append(attrs, g.Transform.List()...)
	if g.Fill != "" {
		attrs = append(attrs, appendString("fill", g.Fill))
	}
	writeElement(w, "g", attrs, func() {
		g.List.Render(w)
	})
}

func (g *Group) AsElement() Element {
	return g
}

type Rect struct {
	node
	List

	Fill string
	RX   float64
	RY   float64
	Dim
	Pos
	Stroke
	Transform
}

func NewRect(options ...Option) Rect {
	r := Rect{
		Fill:   "none",
		Stroke: DefaultStroke,
	}
	for _, o := range options {
		o(&r)
	}
	return r
}

func (r *Rect) Render(w Writer) {
	attrs := r.node.attrs()
	attrs = append(attrs, r.Dim.List()...)
	attrs = append(attrs, r.Pos.List()...)
	attrs = append(attrs, r.Stroke.List()...)
	attrs = append(attrs, r.Transform.List()...)
	if r.RX > 0 {
		attrs = append(attrs, appendFloat("rx", r.RX))
	}
	if r.RY > 0 {
		attrs = append(attrs, appendFloat("ry", r.RY))
	}
	if r.Fill != "" {
		attrs = append(attrs, appendString("fill", r.Fill))
	}
	writeElement(w, "rect", attrs, func() {
		r.List.Render(w)
	})
}

func (r *Rect) AsElement() Element {
	return r
}

type Polygon struct {
	node
	List

	Points []Pos
	Fill   string
	Stroke
	Transform
}

func (p *Polygon) Render(w Writer) {
	attrs := p.node.attrs()
	attrs = append(attrs, p.attrs()...)
	attrs = append(attrs, p.Stroke.List()...)
	attrs = append(attrs, p.Transform.List()...)
	if p.Fill != "" {
		attrs = append(attrs, appendString("fill", p.Fill))
	}
	writeElement(w, "polygon", attrs, func() {
		p.List.Render(w)
	})
}

func (p *Polygon) AsElement() Element {
	return p
}

func (p *Polygon) attrs() []string {
	var list []float64
	for i := range p.Points {
		list = append(list, p.Points[i].array()...)
	}
	a := appendFloatPair("points", list)
	return []string{a}
}

type Ellipse struct {
	node
	List

	Pos
	RX   float64
	RY   float64
	Fill string
	Stroke
	Transform
}

func (e *Ellipse) Render(w Writer) {
	attrs := e.node.attrs()
	attrs = append(attrs, e.attrs()...)
	attrs = append(attrs, e.Pos.Center()...)
	attrs = append(attrs, e.Stroke.List()...)
	attrs = append(attrs, e.Transform.List()...)
	if e.Fill != "" {
		attrs = append(attrs, appendString("fill", e.Fill))
	}
	writeElement(w, "ellipse", attrs, func() {
		e.List.Render(w)
	})
}

func (e *Ellipse) AsElement() Element {
	return e
}

func (e *Ellipse) attrs() []string {
	var attrs []string
	attrs = append(attrs, appendFloat("rx", e.RX))
	attrs = append(attrs, appendFloat("ry", e.RY))
	return attrs
}

type Circle struct {
	node
	List

	Fill   string
	Radius float64
	Pos
	Transform
}

func NewCircle(options ...Option) Circle {
	var c Circle
	for _, o := range options {
		o(&c)
	}
	return c
}

func (c *Circle) Render(w Writer) {
	attrs := c.node.attrs()
	attrs = append(attrs, c.Pos.Center()...)
	if c.Radius != 0 {
		attrs = append(attrs, appendFloat("r", c.Radius))
	}
	if c.Fill != "" {
		attrs = append(attrs, appendString("fill", c.Fill))
	}
	writeElement(w, "circle", attrs, func() {
		c.List.Render(w)
	})
}

func (c *Circle) AsElement() Element {
	return c
}

type Text struct {
	node
	Literal string

	Fill   string
	Anchor string
	Pos
	Font
	Stroke
	Transform
}

func NewText(str string, options ...Option) Text {
	t := Text{Literal: str}
	for _, o := range options {
		o(&t)
	}
	return t
}

func (t *Text) Render(w Writer) {
	attrs := t.node.attrs()
	attrs = append(attrs, t.Pos.List()...)
	attrs = append(attrs, t.Font.List()...)
	attrs = append(attrs, t.Stroke.List()...)
	attrs = append(attrs, t.Transform.List()...)
	if t.Fill != "" {
		attrs = append(attrs, appendString("fill", t.Fill))
	}
	if t.Anchor != "" {
		attrs = append(attrs, appendString("text-anchor", t.Anchor))
	}
	writeElement(w, "text", attrs, func() {
		w.WriteString(t.Literal)
	})
}

func (t *Text) AsElement() Element {
	return t
}

type Line struct {
	node

	Starts Pos
	Ends   Pos
	Fill   string
	Stroke
}

func NewLine(starts, ends Pos, options ...Option) Line {
	i := Line{
		Starts: starts,
		Ends:   ends,
		Stroke: DefaultStroke,
	}
	for _, o := range options {
		o(&i)
	}
	return i
}

func (i *Line) Render(w Writer) {
	attrs := i.node.attrs()
	attrs = append(attrs, i.attrs()...)
	attrs = append(attrs, i.Stroke.List()...)
	if i.Fill != "" {
		attrs = append(attrs, appendString("fill", i.Fill))
	}
	writeOpenElement(w, "line", true, attrs)
}

func (i *Line) AsElement() Element {
	return i
}

func (i *Line) attrs() []string {
	var attrs []string
	attrs = append(attrs, appendFloat("x1", i.Starts.X))
	attrs = append(attrs, appendFloat("y1", i.Starts.Y))
	attrs = append(attrs, appendFloat("x2", i.Ends.X))
	attrs = append(attrs, appendFloat("y2", i.Ends.Y))
	return attrs
}

type PolyLine struct {
	node

	Points []Pos
	Stroke
	Transform
}

func (p *PolyLine) Render(w Writer) {
	attrs := p.node.attrs()
	attrs = append(attrs, p.attrs()...)
	attrs = append(attrs, p.Stroke.List()...)
	attrs = append(attrs, p.Transform.List()...)
	writeOpenElement(w, "polyline", true, attrs)
}

func (p *PolyLine) AsElement() Element {
	return p
}

func (p *PolyLine) attrs() []string {
	var list []float64
	for i := range p.Points {
		list = append(list, p.Points[i].array()...)
	}
	a := appendFloatPair("points", list)
	return []string{a}
}

const (
	cmdMoveToAbs          = "M"
	cmdMoveToRel          = "m"
	cmdLineToAbs          = "L"
	cmdLineToRel          = "l"
	cmdHorizontalAbs      = "H"
	cmdHorizontalRel      = "h"
	cmdVerticalAbs        = "V"
	cmdVerticalRel        = "v"
	cmdArcAbs             = "A"
	cmdArcRel             = "a"
	cmdClosePath          = "Z"
	cmdCubicAbs           = "C"
	cmdCubicRel           = "c"
	cmdCubicSimpleAbs     = "S"
	cmdCubicSimpleRel     = "s"
	cmdQuadraticAbs       = "Q"
	cmdQuadraticRel       = "q"
	cmdQuadraticSimpleAbs = "T"
	cmdQuadraticSimpleRel = "t"
)

type command struct {
	cmd    string
	values [][]float64
}

func makeCommand(cmd string, values ...[]float64) command {
	vs := make([][]float64, 0, len(values))
	vs = append(vs, values...)
	return command{
		cmd:    cmd,
		values: vs,
	}
}

func (c command) String() string {
	buf := []byte(c.cmd)
	for i := range c.values {
		if i > 0 {
			buf = append(buf, comma)
		}
		buf = append(buf, space)
		for j := range c.values[i] {
			if j > 0 {
				buf = append(buf, space)
			}
			buf = strconv.AppendFloat(buf, c.values[i][j], 'f', 2, 64)
		}
	}
	return string(buf)
}

type Path struct {
	node
	commands []command

	Fill string
	Stroke
	Transform
}

func NewPath(options ...Option) Path {
	var p Path
	p.Stroke = DefaultStroke
	for _, o := range options {
		o(&p)
	}
	return p
}

func (p *Path) Render(w Writer) {
	attrs := p.node.attrs()
	attrs = append(attrs, p.attrs()...)
	attrs = append(attrs, p.Stroke.List()...)
	attrs = append(attrs, p.Transform.List()...)
	if p.Fill != "" {
		attrs = append(attrs, appendString("fill", p.Fill))
	}
	writeOpenElement(w, "path", true, attrs)
}

func (p *Path) AsElement() Element {
	return p
}

func (p *Path) attrs() []string {
	var attrs []string
	for _, c := range p.commands {
		attrs = append(attrs, c.String())
	}
	a := appendStringArray("d", attrs, space)
	return []string{a}
}

func (p *Path) AbsMoveTo(pos Pos) {
	c := makeCommand(cmdMoveToAbs, pos.array())
	p.commands = append(p.commands, c)
}

func (p *Path) RelMoveTo(pos Pos) {
	c := makeCommand(cmdMoveToRel, pos.array())
	p.commands = append(p.commands, c)
}

func (p *Path) AbsLineTo(pos Pos) {
	c := makeCommand(cmdLineToAbs, pos.array())
	p.commands = append(p.commands, c)
}

func (p *Path) RelLineTo(pos Pos) {
	c := makeCommand(cmdLineToRel, pos.array())
	p.commands = append(p.commands, c)
}

func (p *Path) AbsHorizontalLine(x float64) {
	c := makeCommand(cmdHorizontalAbs, []float64{x})
	p.commands = append(p.commands, c)
}

func (p *Path) RelHorizontalLine(x float64) {
	c := makeCommand(cmdHorizontalRel, []float64{x})
	p.commands = append(p.commands, c)
}

func (p *Path) AbsVerticalLine(y float64) {
	c := makeCommand(cmdVerticalAbs, []float64{y})
	p.commands = append(p.commands, c)
}

func (p *Path) RelVerticalLine(y float64) {
	c := makeCommand(cmdVerticalRel, []float64{y})
	p.commands = append(p.commands, c)
}

func (p *Path) ClosePath() {
	c := makeCommand(cmdClosePath)
	p.commands = append(p.commands, c)
}

func (p *Path) AbsCubicCurve(pos, start, end Pos) {
	c := makeCommand(cmdCubicAbs, start.array(), end.array(), pos.array())
	p.commands = append(p.commands, c)
}

func (p *Path) RelCubicCurve(pos, start, end Pos) {
	c := makeCommand(cmdCubicRel, start.array(), end.array(), pos.array())
	p.commands = append(p.commands, c)
}

func (p *Path) AbsCubicCurveSimple(pos, ctrl Pos) {
	c := makeCommand(cmdCubicSimpleAbs, ctrl.array(), pos.array())
	p.commands = append(p.commands, c)
}

func (p *Path) RelCubicCurveSimple(pos, ctrl Pos) {
	c := makeCommand(cmdCubicSimpleRel, ctrl.array(), pos.array())
	p.commands = append(p.commands, c)
}

func (p *Path) AbsQuadraticCurve(pos, ctrl Pos) {
	c := makeCommand(cmdQuadraticAbs, ctrl.array(), pos.array())
	p.commands = append(p.commands, c)
}

func (p *Path) RelQuadraticCurve(pos, ctrl Pos) {
	c := makeCommand(cmdQuadraticRel, ctrl.array(), pos.array())
	p.commands = append(p.commands, c)
}

func (p *Path) AbsQuadraticCurveSimple(pos Pos) {
	c := makeCommand(cmdQuadraticSimpleAbs, pos.array())
	p.commands = append(p.commands, c)
}

func (p *Path) RelQuadraticCurveSimple(pos Pos) {
	c := makeCommand(cmdQuadraticSimpleRel, pos.array())
	p.commands = append(p.commands, c)
}

func (p *Path) Arc(pos Pos, rx, ry float64, large, sweep bool) {
}

func writeElement(w Writer, name string, attrs []string, inner func()) {
	closed := inner == nil
	writeOpenElement(w, name, closed, attrs)
	if !closed {
		inner()
		writeCloseElement(w, name)
	}
}

func writeOpenElement(w Writer, name string, closed bool, attrs []string) {
	w.WriteByte(langle)
	w.WriteString(name)
	for i := range attrs {
		w.WriteByte(space)
		w.WriteString(attrs[i])
	}
	if closed {
		w.WriteByte(space)
		w.WriteByte(slash)
	}
	w.WriteByte(rangle)
}

func writeCloseElement(w Writer, name string) {
	w.WriteByte(langle)
	w.WriteByte(slash)
	w.WriteString(name)
	w.WriteByte(rangle)
}
