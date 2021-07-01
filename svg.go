package svg

import (
	"io"
	"math"
	"strconv"
	"strings"
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

type Option func(Element) error

func WithID(id string) Option {
	return func(e Element) error {
		e.setId(id)
		return nil
	}
}

func WithClass(class ...string) Option {
	return func(e Element) error {
		e.setClass(class)
		return nil
	}
}

func WithStyle(prop string, values ...string) Option {
	return func(e Element) error {
		e.setStyle(prop, values)
		return nil
	}
}

func WithTranslate(x, y float64) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Rect:
			e.Transform.TX, e.Transform.TY = x, y
		case *Text:
			e.Transform.TX, e.Transform.TY = x, y
		case *Group:
			e.Transform.TX, e.Transform.TY = x, y
		}
		return nil
	}
}

func WithRotate(a, x, y float64) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Rect:
			e.Transform.RA = a
			e.Transform.RX, e.Transform.RY = x, y
		case *Text:
			e.Transform.RA = a
			e.Transform.RX, e.Transform.RY = x, y
		case *Group:
			e.Transform.RA = a
			e.Transform.RX, e.Transform.RY = x, y
		}
		return nil
	}
}

func WithScale(x, y float64) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Rect:
			e.Transform.SX, e.Transform.SY = x, y
		case *Text:
			e.Transform.SX, e.Transform.SY = x, y
		case *Group:
			e.Transform.SX, e.Transform.SY = x, y
		}
		return nil
	}
}

func WithSkewX(x float64) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Rect:
			e.Transform.KX = x
		case *Text:
			e.Transform.KX = x
		case *Group:
			e.Transform.KX = x
		}
		return nil
	}
}

func WithSkewY(y float64) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Rect:
			e.Transform.KY = y
		case *Text:
			e.Transform.KY = y
		case *Group:
			e.Transform.KY = y
		}
		return nil
	}
}

func WithAnchor(anchor string) Option {
	return func(e Element) error {
		if e, ok := e.(*Text); ok {
			e.Anchor = anchor
		}
		return nil
	}
}

func WithFill(fill string) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Line:
			e.Fill = fill
		case *Path:
			e.Fill = fill
		case *Rect:
			e.Fill = fill
		case *Circle:
			e.Fill = fill
		case *Text:
			e.Fill = fill
		case *Group:
			e.Fill = fill
		case *SVG:
		default:
		}
		return nil
	}
}

func WithFont(f Font) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Text:
			e.Font = f
		default:
		}
		return nil
	}
}

func WithDim(d Dim) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Rect:
			e.Dim = d
		case *SVG:
			e.Dim = d
		default:
		}
		return nil
	}
}

func WithDimension(x, y float64) Option {
	return WithDim(NewDim(x, y))
}

func WithRadius(r float64) Option {
	return func(e Element) error {
		if e, ok := e.(*Circle); ok {
			e.Radius = r
		}
		return nil
	}
}

func WithPos(p Pos) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Rect:
			e.Pos = p
		case *Circle:
			e.Pos = p
		case *Text:
			e.Pos = p
		}
		return nil
	}
}

func WithPosition(x, y float64) Option {
	return WithPos(NewPos(x, y))
}

func WithStroke(s Stroke) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Rect:
			e.Stroke = s
		case *Group:
			e.Stroke = s
		case *Text:
			e.Stroke = s
		default:
		}
		return nil
	}
}

type Font struct {
	Family []string
	Fill   string
	Size   float64
}

func NewFont(size float64) Font {
	return Font{Size: size}
}

func (f Font) List() []string {
	var attrs []string
	if len(f.Family) > 0 {
		attrs = append(attrs, appendStringArray("font-family", f.Family, comma))
	}
	if f.Fill != "" {
		attrs = append(attrs, appendString("fill", f.Fill))
	}
	if f.Size > 0 {
		attrs = append(attrs, appendFloat("font-size", f.Size))
	}
	return attrs
}

type Pos struct {
	X float64
	Y float64
}

func NewPos(x, y float64) Pos {
	return Pos{
		X: x,
		Y: y,
	}
}

func (p Pos) List() []string {
	var attrs []string
	if p.X != 0 {
		attrs = append(attrs, appendFloat("x", p.X))
	}
	if p.Y != 0 {
		attrs = append(attrs, appendFloat("y", p.Y))
	}
	return attrs
}

func (p Pos) array() []float64 {
	return []float64{p.X, p.Y}
}

type Dim struct {
	W float64
	H float64
}

func NewDim(w, h float64) Dim {
	return Dim{
		W: w,
		H: h,
	}
}

func (d Dim) List() []string {
	var attrs []string
	attrs = append(attrs, appendFloat("width", d.W))
	attrs = append(attrs, appendFloat("height", d.H))
	return attrs
}

type Stroke struct {
	Dash  []int
	Width float64
	Color string
}

func NewStroke(fill string, width int) Stroke {
	return Stroke{
		Color: fill,
		Width: float64(width),
	}
}

func (s Stroke) List() []string {
	var attrs []string
	if len(s.Dash) > 0 {
		attrs = append(attrs, appendIntArray("stroke-dasharray", s.Dash, space))
	}
	if s.Width > 0 {
		attrs = append(attrs, appendFloat("stroke-width", s.Width))
	}
	if s.Color != "" {
		attrs = append(attrs, appendString("stroke", s.Color))
	}
	return attrs
}

type Transform struct {
	TX float64
	TY float64

	SX float64
	SY float64

	RA float64
	RX float64
	RY float64

	KX float64
	KY float64
}

func (t Transform) List() []string {
	var attrs []string
	if t.TX != 0 || t.TY != 0 {
		attrs = append(attrs, appendFunc("translate", t.TX, t.TY))
	}
	if t.SX != 0 || t.SY != 0 {
		attrs = append(attrs, appendFunc("scale", t.SX, t.SY))
	}
	if t.RA != 0 {
		attrs = append(attrs, appendFunc("rotate", t.RA, t.RX, t.RY))
	}
	if t.KX != 0 {
		attrs = append(attrs, appendFunc("skewX", t.KX))
	}
	if t.KY != 0 {
		attrs = append(attrs, appendFunc("skewY", t.KY))
	}
	if len(attrs) == 0 {
		return nil
	}
	a := appendString("transform", strings.Join(attrs, " "))
	return []string{a}
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
	attrs = append(attrs, c.attrs()...)
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

func (c *Circle) attrs() []string {
	var attrs []string
	attrs = append(attrs, appendFloat("cx", c.Pos.X))
	attrs = append(attrs, appendFloat("cy", c.Pos.Y))
	return attrs
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

type PolyLine struct{}

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

const (
	quote  = '"'
	space  = ' '
	comma  = ','
	equal  = '='
	slash  = '/'
	langle = '<'
	rangle = '>'
	lparen = '('
	rparen = ')'
)

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

const defaultPrecision = 2

func appendFunc(name string, list ...float64) string {
	buf := []byte(name)
	buf = append(buf, lparen)
	for i := range list {
		if i > 0 {
			buf = append(buf, comma)
		}
		p := defaultPrecision
		if math.Ceil(list[i]) == list[i] {
			p = 0
		}
		buf = strconv.AppendFloat(buf, list[i], 'f', p, 64)
	}
	buf = append(buf, rparen)
	return string(buf)
}

func appendFloat(attr string, v float64) string {
	var (
		buf  = []byte(attr)
		prec = defaultPrecision
	)
	if math.Ceil(v) == v {
		prec = 0
	}
	buf = append(buf, equal, quote)
	buf = strconv.AppendFloat(buf, v, 'f', prec, 64)
	buf = append(buf, quote)
	return string(buf)
}

func appendString(attr, v string) string {
	buf := []byte(attr)
	buf = append(buf, equal)
	buf = strconv.AppendQuote(buf, v)
	return string(buf)
}

func appendStringArray(attr string, list []string, sep byte) string {
	buf := []byte(attr)
	buf = append(buf, equal, quote)
	for i := range list {
		if i > 0 {
			buf = append(buf, sep)
			if sep != space {
				buf = append(buf, space)
			}
		}
		buf = append(buf, []byte(list[i])...)
	}
	buf = append(buf, quote)
	return string(buf)
}

func appendIntArray(attr string, list []int, sep byte) string {
	buf := []byte(attr)
	buf = append(buf, equal, quote)
	for i := range list {
		if i > 0 {
			buf = append(buf, sep)
		}
		buf = strconv.AppendInt(buf, int64(list[i]), 10)
	}
	buf = append(buf, quote)
	return string(buf)
}
