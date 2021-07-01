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
	Title string
	Desc  string

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
	List []Element
}

func NewList(es ...Element) List {
	list := make([]Element, len(es))
	copy(list, es)
	return List{
		List: list,
	}
}

func (i *List) Append(e Element) {
	i.List = append(i.List, e)
}

func (i *List) Render(w Writer) {
	for _, e := range i.List {
		e.Render(w)
	}
}

type Defs struct {
	List
}

func (d *Defs) Render(w Writer) {
	writeElement(w, "defs", nil, func() {
		d.List.Render(w)
	})
}

func (d *Defs) AsElement() Element {
	return d
}

type Use struct {
	node
	Ref string

	Pos
	Dim
	Stroke
	Fill
}

func (u *Use) Render(w Writer) {
	if u.Ref == "" {
		return
	}
	attrs := u.node.attrs()
	attrs = append(attrs, u.Pos.List()...)
	attrs = append(attrs, u.Dim.List()...)
	attrs = append(attrs, u.Stroke.List()...)
	attrs = append(attrs, u.Fill.List()...)
	attrs = append(attrs, appendString("href", u.Ref))
	writeOpenElement(w, "use", true, attrs)
}

func (u *Use) AsElement() Element {
	return u
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
		writeTitle(w, s.Title)
		writeDesc(w, s.Desc)
		s.List.Render(w)
	})
}

func (s *SVG) AsElement() Element {
	return s
}

type Group struct {
	node
	List

	Fill
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
	attrs = append(attrs, g.Fill.List()...)
	writeElement(w, "g", attrs, func() {
		writeTitle(w, g.Title)
		writeDesc(w, g.Desc)
		g.List.Render(w)
	})
}

func (g *Group) AsElement() Element {
	return g
}

type Rect struct {
	node
	List

	RX float64
	RY float64
	Fill
	Dim
	Pos
	Stroke
	Transform
}

func NewRect(options ...Option) Rect {
	var r Rect
	r.Stroke = DefaultStroke
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
	attrs = append(attrs, r.Fill.List()...)
	writeElement(w, "rect", attrs, func() {
		writeTitle(w, r.Title)
		writeDesc(w, r.Desc)
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
	Fill
	Stroke
	Transform
}

func (p *Polygon) Render(w Writer) {
	attrs := p.node.attrs()
	attrs = append(attrs, p.attrs()...)
	attrs = append(attrs, p.Stroke.List()...)
	attrs = append(attrs, p.Transform.List()...)
	attrs = append(attrs, p.Fill.List()...)
	writeElement(w, "polygon", attrs, func() {
		writeTitle(w, p.Title)
		writeDesc(w, p.Desc)
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
	RX float64
	RY float64
	Fill
	Stroke
	Transform
}

func (e *Ellipse) Render(w Writer) {
	attrs := e.node.attrs()
	attrs = append(attrs, e.attrs()...)
	attrs = append(attrs, e.Pos.Center()...)
	attrs = append(attrs, e.Stroke.List()...)
	attrs = append(attrs, e.Transform.List()...)
	attrs = append(attrs, e.Fill.List()...)
	writeElement(w, "ellipse", attrs, func() {
		writeTitle(w, e.Title)
		writeDesc(w, e.Desc)
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

	Radius float64
	Pos
	Fill
	Stroke
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
	attrs = append(attrs, c.Fill.List()...)
	attrs = append(attrs, c.Stroke.List()...)
	writeElement(w, "circle", attrs, func() {
		writeTitle(w, c.Title)
		writeDesc(w, c.Desc)
		c.List.Render(w)
	})
}

func (c *Circle) AsElement() Element {
	return c
}

type Text struct {
	node
	Literal string

	Anchor string
	Fill
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
	attrs = append(attrs, t.Fill.List()...)
	if t.Anchor != "" {
		attrs = append(attrs, appendString("text-anchor", t.Anchor))
	}
	writeElement(w, "text", attrs, func() {
		writeTitle(w, t.Title)
		writeDesc(w, t.Desc)
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
	Fill
	Stroke
	Transform
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
	attrs = append(attrs, i.Fill.List()...)
	attrs = append(attrs, i.Transform.List()...)
	writeElement(w, "line", attrs, func() {
		writeTitle(w, i.Title)
		writeDesc(w, i.Desc)
	})
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
	Fill
	Transform
}

func (p *PolyLine) Render(w Writer) {
	attrs := p.node.attrs()
	attrs = append(attrs, p.attrs()...)
	attrs = append(attrs, p.Stroke.List()...)
	attrs = append(attrs, p.Fill.List()...)
	attrs = append(attrs, p.Transform.List()...)
	writeElement(w, "polyline", attrs, func() {
		writeTitle(w, p.Title)
		writeDesc(w, p.Desc)
	})
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

	Fill
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
	attrs = append(attrs, p.Fill.List()...)
	attrs = append(attrs, p.Transform.List()...)
	attrs = append(attrs, p.Fill.List()...)
	writeElement(w, "path", attrs, func() {
		writeTitle(w, p.Title)
		writeDesc(w, p.Desc)
	})
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

func writeTitle(w Writer, str string) {
	if str == "" {
		return
	}
	writeString(w, "title", str)
}

func writeDesc(w Writer, str string) {
	if str == "" {
		return
	}
	writeString(w, "desc", str)
}

func writeString(w Writer, name, str string) {
	writeOpenElement(w, name, false, nil)
	w.WriteString(str)
	writeCloseElement(w, name)
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
