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
	defaultWidth  = 800
	defaultHeight = 600
)

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

type Literal string

func NewLiteral(str string) Literal {
	return Literal(str)
}

func (i Literal) Render(w Writer) {
	w.WriteString(string(i))
}

func (_ Literal) setId(_ string)                {}
func (_ Literal) setClass(_ []string)           {}
func (_ Literal) setStyle(_ string, _ []string) {}

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
	if e == nil {
		return
	}
	if e, ok := e.(*SVG); ok {
		e.OmitProlog = ok
	}
	i.List = append(i.List, e)
}

func (i *List) Render(w Writer) {
	for _, e := range i.List {
		if e == nil {
			continue
		}
		e.Render(w)
	}
}

type Defs struct {
	List

	Transform
}

func (d *Defs) Render(w Writer) {
	attrs := d.Transform.Attributes()
	writeElement(w, "defs", attrs, func() {
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
	Transform
}

func (u *Use) Render(w Writer) {
	if u.Ref == "" {
		return
	}
	var list List
	u.render(w, "use", list, u, u.Pos, u.Dim, u.Fill, u.Stroke, u.Transform)
}

func (u *Use) AsElement() Element {
	return u
}

func (u *Use) Attributes() []string {
	a := appendString("href", u.Ref)
	return []string{a}
}

type SVG struct {
	node
	List

	OmitProlog    bool
	PreserveRatio []string
	Box
	Pos
	Dim
	Fill
	Stroke
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
	if !s.OmitProlog {
		w.WriteString(prolog)
	}
	s.render(w, "svg", s.List, s, s.Pos, s.Dim, s.Box)
}

func (s *SVG) AsElement() Element {
	return s
}

func (s *SVG) Attributes() []string {
	var attrs []string
	attrs = append(attrs, appendString("xmlns", namespace))
	if len(s.PreserveRatio) == 0 {
		return attrs
	}
	attrs = append(attrs, appendStringArray("preserveAspectRatio", s.PreserveRatio, space))
	return attrs
}

type ClipPath struct {
	node
	List

	Fill
	Stroke
	Transform
}

func (c *ClipPath) Render(w Writer) {
	c.render(w, "clipPath", c.List, c.Fill, c.Stroke, c.Transform)
}

func (c *ClipPath) AsElement() Element {
	return c
}

type TextPath struct {
	node
	Literal string

	Path    string
	Adjust  string
	Method  string
	Side    string
	Spacing string
	Offset  float64
	Length  float64
	Fill
	Stroke
	Transform
}

func NewTextPath(literal, path string, options ...Option) TextPath {
	t := TextPath{
		Literal: literal,
		Path:    path,
		Fill:    DefaultFill,
	}
	for _, o := range options {
		o(&t)
	}
	return t
}

func (t *TextPath) Render(w Writer) {
	if t.Path == "" {
		return
	}
	writeElement(w, "text", nil, func() {
		list := NewList(Literal(t.Literal))
		t.render(w, "textPath", list, t, t.Fill, t.Stroke, t.Transform)
	})
}

func (t *TextPath) AsElement() Element {
	return t
}

func (t *TextPath) Attributes() []string {
	var attrs []string
	attrs = append(attrs, appendString("href", "#"+t.Path))
	if t.Method != "" {
		attrs = append(attrs, appendString("method", t.Method))
	}
	if t.Adjust != "" {
		attrs = append(attrs, appendString("lengthAdjust", t.Adjust))
	}
	if t.Side != "" {
		attrs = append(attrs, appendString("side", t.Side))
	}
	if t.Spacing != "" {
		attrs = append(attrs, appendString("spacing", t.Spacing))
	}
	if t.Length != 0 {
		attrs = append(attrs, appendFloat("textLength", t.Length))
	}
	if t.Offset != 0 {
		attrs = append(attrs, appendFloat("startOffset", t.Offset))
	}
	return attrs
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
	g.render(w, "g", g.List, g.Stroke, g.Fill, g.Transform)
}

func (g *Group) AsElement() Element {
	return g
}

type Image struct {
	node

	Ref           string
	PreserveRatio []string
	Pos
	Dim
}

func NewImage(ref string, options ...Option) Image {
	i := Image{Ref: ref}
	for _, o := range options {
		o(&i)
	}
	return i
}

func (i *Image) Render(w Writer) {
	if i.Ref == "" {
		return
	}
	var list List
	i.render(w, "image", list, i, i.Pos, i.Dim)
}

func (i *Image) AsElement() Element {
	return i
}

func (i *Image) Attributes() []string {
	var attrs []string
	attrs = append(attrs, appendString("href", i.Ref))
	if len(i.PreserveRatio) > 0 {
		a := appendStringArray("preserveAspectRatio", i.PreserveRatio, space)
		attrs = append(attrs, a)
	}
	return attrs
}

type Mask struct {
	node
	List

	Pos
	Dim
	Fill
	Stroke
	Transform
}

func (m *Mask) Render(w Writer) {
	m.render(w, "mask", m.List, m.Pos, m.Dim, m.Fill, m.Stroke, m.Transform)
}

func (m *Mask) AsElement() Element {
	return m
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
	r.Fill = NewFill("none")
	for _, o := range options {
		o(&r)
	}
	return r
}

func (r *Rect) Render(w Writer) {
	r.render(w, "rect", r.List, r, r.Dim, r.Pos, r.Stroke, r.Transform, r.Fill)
}

func (r *Rect) AsElement() Element {
	return r
}

func (r *Rect) Attributes() []string {
	var attrs []string
	if r.RX != 0 {
		attrs = append(attrs, appendFloat("rx", r.RX))
	}
	if r.RY != 0 {
		attrs = append(attrs, appendFloat("ry", r.RY))
	}
	return attrs
}

type Polygon struct {
	node
	List

	Points []Pos
	Fill
	Stroke
	Transform
}

func NewPolygon(pos []Pos, options ...Option) Polygon {
	p := Polygon{Points: pos}
	for _, o := range options {
		o(&p)
	}
	return p
}

func (p *Polygon) Render(w Writer) {
	p.render(w, "polygon", p.List, p, p.Fill, p.Stroke, p.Transform)
}

func (p *Polygon) AsElement() Element {
	return p
}

func (p *Polygon) Attributes() []string {
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

func NewEllipse(options ...Option) Ellipse {
	var e Ellipse
	for _, o := range options {
		o(&e)
	}
	return e
}

func (e *Ellipse) Render(w Writer) {
	e.render(w, "ellipse", e.List, e, e.Stroke, e.Fill, e.Transform)
}

func (e *Ellipse) AsElement() Element {
	return e
}

func (e *Ellipse) Attributes() []string {
	var attrs []string
	attrs = append(attrs, appendFloat("rx", e.RX))
	attrs = append(attrs, appendFloat("ry", e.RY))
	attrs = append(attrs, e.Pos.Center()...)
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
	c.Fill = NewFill("none")
	for _, o := range options {
		o(&c)
	}
	return c
}

func (c *Circle) Render(w Writer) {
	c.render(w, "circle", c.List, c, c.Fill, c.Stroke, c.Transform)
}

func (c *Circle) AsElement() Element {
	return c
}

func (c *Circle) Attributes() []string {
	a := appendFloat("r", c.Radius)
	attrs := []string{a}
	return append(attrs, c.Pos.Center()...)
}

type Text struct {
	node
	List

	Shift  Pos
	Anchor string
	Adjust string
	Length float64
	Fill
	Pos
	Font
	Stroke
	Transform
}

func NewText(str string, options ...Option) Text {
	var t Text
	t.Append(Literal(str))
	for _, o := range options {
		o(&t)
	}
	return t
}

func (t *Text) Render(w Writer) {
	t.render(w, "text", t.List, t, t.Pos, t.Font, t.Fill, t.Stroke, t.Transform)
}

func (t *Text) AsElement() Element {
	return t
}

func (t *Text) Attributes() []string {
	var attrs []string
	attrs = append(attrs, appendFloat("dx", t.Shift.X))
	attrs = append(attrs, appendFloat("dy", t.Shift.Y))
	if t.Anchor != "" {
		attrs = append(attrs, appendString("text-anchor", t.Anchor))
	}
	if t.Adjust != "" {
		attrs = append(attrs, appendString("lengthAdjust", t.Adjust))
	}
	if t.Length != 0 {
		attrs = append(attrs, appendFloat("textLength", t.Length))
	}
	return attrs
}

type TextSpan struct {
	node
	Literal string

	Pos
	Shift  Pos
	Adjust string
	Length float64
	Rotate []float64
}

func (t *TextSpan) Render(w Writer) {
	list := NewList(Literal(t.Literal))
	t.render(w, "tspan", list, t, t.Pos)
}

func (t *TextSpan) AsElement() Element {
	return t
}

func (t *TextSpan) Attributes() []string {
	var attrs []string
	attrs = append(attrs, appendFloat("dx", t.Shift.X))
	attrs = append(attrs, appendFloat("dy", t.Shift.Y))
	if t.Adjust != "" {
		attrs = append(attrs, appendString("lengthAdjust", t.Adjust))
	}
	if t.Length != 0 {
		attrs = append(attrs, appendFloat("textLength", t.Length))
	}
	return attrs
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
	}
	for _, o := range options {
		o(&i)
	}
	return i
}

func (i *Line) Render(w Writer) {
	var list List
	i.render(w, "line", list, i, i.Stroke, i.Fill, i.Transform)
}

func (i *Line) AsElement() Element {
	return i
}

func (i *Line) Attributes() []string {
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

func NewPolyLine(pos []Pos, options ...Option) PolyLine {
	p := PolyLine{Points: pos}
	for _, o := range options {
		o(&p)
	}
	return p
}

func (p *PolyLine) Render(w Writer) {
	var list List
	p.render(w, "polyline", list, p, p.Fill, p.Stroke, p.Transform)
}

func (p *PolyLine) AsElement() Element {
	return p
}

func (p *PolyLine) Attributes() []string {
	var list []float64
	for i := range p.Points {
		list = append(list, p.Points[i].array()...)
	}
	a := appendFloatPair("points", list)
	return []string{a}
}

type Marker struct {
	node
}

func (m *Marker) Render(w Writer) {

}

func (m *Marker) AsElement() Element {
	return m
}

type Symbol struct {
	node
}

func (s *Symbol) Render(w Writer) {

}

func (s *Symbol) AsElement() Element {
	return s
}

type Switch struct {
	node
}

func (s *Switch) Render(w Writer) {

}

func (s *Switch) AsElement() Element {
	return s
}

type Style struct {
	node
}

func (s *Style) Render(w Writer) {

}

func (s *Style) AsElement() Element {
	return s
}

type Script struct {
	node
}

func (s *Script) Render(w Writer) {

}

func (s *Script) AsElement() Element {
	return s
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
	p.Fill = NewFill("none")
	for _, o := range options {
		o(&p)
	}
	return p
}

func (p *Path) Render(w Writer) {
	var list List
	p.render(w, "path", list, p, p.Fill, p.Stroke, p.Transform)
}

func (p *Path) AsElement() Element {
	return p
}

func (p *Path) Attributes() []string {
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

func (p *Path) AbsArcTo(pos Pos, rx, ry, rot float64, large, sweep bool) {
	args := []float64{
		rx,
		ry,
		rot,
	}
	if large {
		args = append(args, 1)
	} else {
		args = append(args, 0)
	}
	if sweep {
		args = append(args, 1)
	} else {
		args = append(args, 0)
	}
	args = append(args, pos.array()...)
	c := makeCommand(cmdArcAbs, args)
	p.commands = append(p.commands, c)
}

func (p *Path) RelArcTo(pos Pos, rx, ry, rot float64, large, sweep bool) {
	args := []float64{
		rx,
		ry,
		rot,
	}
	if large {
		args = append(args, 1)
	} else {
		args = append(args, 0)
	}
	if sweep {
		args = append(args, 1)
	} else {
		args = append(args, 0)
	}
	args = append(args, pos.array()...)
	c := makeCommand(cmdArcAbs, args)
	p.commands = append(p.commands, c)
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
			v := c.values[i][j]
			buf = strconv.AppendFloat(buf, v, 'f', getPrecision(v), 64)
		}
	}
	return string(buf)
}

type node struct {
	Title string
	Desc  string
	Data  []Datum

	Display    string
	Visibility string

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

func (n *node) Attributes() []string {
	var attrs []string
	if n.Id != "" {
		attrs = append(attrs, appendString("id", n.Id))
	}
	if len(n.Class) > 0 {
		attrs = append(attrs, appendStringArray("class", n.Class, space))
	}
	if len(n.Data) > 0 {
		for i := range n.Data {
			attrs = append(attrs, n.Data[i].Attributes()...)
		}
	}
	return attrs
}

func (n *node) render(w Writer, name string, list List, attrs ...Attribute) {
	var as []string
	attrs = append(attrs, n)
	for _, a := range attrs {
		as = append(as, a.Attributes()...)
	}
	writeElement(w, name, as, func() {
		writeTitle(w, n.Title)
		writeDesc(w, n.Desc)
		list.Render(w)
	})
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
