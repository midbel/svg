package svg

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const defaultFontSize = 14

var (
	DefaultStroke   = NewStroke("black", 1)
	DefaultFill     = NewFill("black")
	TransparentFill = NewFill("transparent")
	DefaultFont     = NewFont(defaultFontSize)
)

const (
	UnitEM  = "em"
	UnitEX  = "ex"
	UnitPX  = "px"
	UnitPT  = "pt"
	UnitPC  = "pc"
	UnitCM  = "cm"
	UnitMM  = "mm"
	UnitIN  = "in"
	UnitPer = "%"
)

type Number struct {
	Value float64
	unit  string
}

type Attribute interface {
	Attributes() []string
}

type Clipping struct {
	Path  string
	Box   string
	Rule  string
	Units string
}

func ClipURL(url string) Clipping {
	return Clipping{Path: url}
}

func ClipShape(shape, box string) Clipping {
	return Clipping{
		Path: shape,
		Box:  box,
	}
}

func (c Clipping) Attributes() []string {
	if c.Path == "" {
		return nil
	}
	var attrs []string
	if c.Box != "" {
		attrs = append(attrs, appendString("clip-path", "")) // TBD
	} else {
		attrs = append(attrs, appendString("clip-path", "")) // TBD
	}
	if c.Rule != "" {
		attrs = append(attrs, appendString("clip-rule", c.Rule))
	}
	if c.Units != "" {
		attrs = append(attrs, appendString("clipPathUnits", c.Units))
	}
	return attrs
}

type Datum struct {
	Name  string
	Value interface{}
}

func (d Datum) Attributes() []string {
	var a string
	switch v := d.Value.(type) {
	case float64:
		a = appendFloat(fmt.Sprintf("data-%s", d.Name), v)
	case string:
		a = appendString(fmt.Sprintf("data-%s", d.Name), v)
	case int:
		a = appendInt(fmt.Sprintf("data-%s", d.Name), int64(v))
	case int64:
		a = appendInt(fmt.Sprintf("data-%s", d.Name), v)
	default:
		return nil
	}
	return []string{a}
}

type Font struct {
	Family  []string
	Style   string
	Weight  string
	Variant string
	Stretch string
	Fill    string
	Size    float64
	Adjust  float64
}

func NewFont(size float64, families ...string) Font {
	return Font{
		Size:   size,
		Family: families,
		Fill:   "black",
	}
}

func (f Font) Attributes() []string {
	var attrs []string
	values := []struct {
		Attr  string
		Value string
	}{
		{Attr: "font-style", Value: f.Style},
		{Attr: "font-weight", Value: f.Weight},
		{Attr: "font-variant", Value: f.Variant},
		{Attr: "font-stretch", Value: f.Stretch},
		{Attr: "fill", Value: f.Fill},
	}
	for _, v := range values {
		if v.Value == "" {
			continue
		}
		attrs = append(attrs, appendString(v.Attr, v.Value))
	}
	if len(f.Family) > 0 {
		attrs = append(attrs, appendStringArray("font-family", f.Family, comma))
	}
	attrs = append(attrs, appendFloat("font-size", f.Size))
	attrs = append(attrs, appendFloat("font-size-adjust", f.Adjust))
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

func (p Pos) Adjust(x, y float64) Pos {
	p.X += x
	p.Y += y
	return p
}

func (p Pos) Attributes() []string {
	var attrs []string
	if p.X != 0 {
		attrs = append(attrs, appendFloat("x", p.X))
	}
	if p.Y != 0 {
		attrs = append(attrs, appendFloat("y", p.Y))
	}
	return attrs
}

func (p Pos) Center() []string {
	var attrs []string
	if p.X != 0 {
		attrs = append(attrs, appendFloat("cx", p.X))
	}
	if p.Y != 0 {
		attrs = append(attrs, appendFloat("cy", p.Y))
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

func (d Dim) Attributes() []string {
	var attrs []string
	if d.W != 0 {
		attrs = append(attrs, appendFloat("width", d.W))
	}
	if d.H != 0 {
		attrs = append(attrs, appendFloat("height", d.H))
	}
	return attrs
}

func (d Dim) array() []float64 {
	return []float64{d.W, d.H}
}

type Box struct {
	Pos
	Dim
}

func (b Box) Attributes() []string {
	if b.W <= 0 || b.H <= 0 {
		return nil
	}
	arr := append([]float64{}, b.Pos.array()...)
	arr = append(arr, b.Dim.array()...)
	a := appendFloatArray("viewBox", arr, space)
	return []string{a}
}

type Stroke struct {
	DashArray  []int
	DashOffset []int
	LineCap    string
	LineJoin   string
	Width      float64
	Opacity    float64
	Miter      float64
	Color      string
}

func NewStroke(fill string, width float64) Stroke {
	return Stroke{
		Color: fill,
		Width: width,
	}
}

func (s Stroke) Fill() Fill {
	return NewFill(s.Color)
}

func (s Stroke) Attributes() []string {
	if s.IsZero() {
		return nil
	}
	var attrs []string
	attrs = append(attrs, appendString("stroke", s.Color))

	if len(s.DashArray) > 0 {
		attrs = append(attrs, appendIntArray("stroke-dasharray", s.DashArray, space))
	}
	if len(s.DashOffset) > 0 {
		attrs = append(attrs, appendIntArray("stroke-dashoffset", s.DashOffset, space))
	}
	if s.LineCap != "" {
		attrs = append(attrs, appendString("stroke-linecap", s.LineCap))
	}
	if s.LineJoin != "" {
		attrs = append(attrs, appendString("stroke-linejoin", s.LineJoin))
	}
	if s.Width > 0 {
		attrs = append(attrs, appendFloat("stroke-width", s.Width))
	}
	if s.Opacity > 0 {
		attrs = append(attrs, appendFloat("stroke-opacity", s.Opacity))
	}
	if s.Miter > 0 {
		attrs = append(attrs, appendFloat("stroke-miterlimit", s.Miter))
	}
	return attrs
}

func (s Stroke) IsZero() bool {
	return s.Color == ""
}

type Fill struct {
	Color   string
	Rule    string
	Opacity float64
}

func NewFill(color string) Fill {
	if color == "" {
		color = "none"
	}
	return Fill{Color: color, Opacity: 100}
}

func (f Fill) Stroke() Stroke {
	return NewStroke(f.Color, 1)
}

func (f Fill) Attributes() []string {
	if f.IsZero() {
		return nil
	}
	var attrs []string
	if f.Color != "" {
		attrs = append(attrs, appendString("fill", f.Color))
	}
	if f.Rule != "" {
		attrs = append(attrs, appendString("fill-rule", f.Rule))
	}
	attrs = append(attrs, appendFloat("fill-opacity", f.Opacity))
	return attrs
}

func (f Fill) IsZero() bool {
	return f.Color == ""
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

func Translate(left, top float64) Transform {
	return Transform{
		TX: left,
		TY: top,
	}
}

func (t *Transform) SkewX(x float64) {
	t.KX = x
}

func (t *Transform) SkewY(y float64) {
	t.KY = y
}

func (t *Transform) Rotate(a, x, y float64) {
	t.RA = a
	t.RX = x
	t.RY = y
}

func (t *Transform) Scale(x, y float64) {
	t.SX = x
	t.SY = y
}

func (t *Transform) Translate(x, y float64) {
	t.TX = x
	t.TY = y
}

func (t Transform) Attributes() []string {
	var attrs []string
	if t.TX != 0 || t.TY != 0 {
		attrs = append(attrs, appendFunc("translate", t.TX, t.TY))
	}
	if t.SX != 0 || t.SY != 0 {
		attrs = append(attrs, appendFunc("scale", t.SX, t.SY))
	}
	if t.RA != 0 {
		var args []float64
		args = append(args, t.RA)
		if t.RX != 0 && t.RY != 0 {
			args = append(args, t.RX, t.RY)
		}
		attrs = append(attrs, appendFunc("rotate", args...))
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

const defaultPrecision = 2

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

func appendFunc(name string, list ...float64) string {
	buf := []byte(name)
	buf = append(buf, lparen)
	for i := range list {
		if i > 0 {
			buf = append(buf, comma)
		}
		buf = strconv.AppendFloat(buf, list[i], 'f', getPrecision(list[i]), 64)
	}
	buf = append(buf, rparen)
	return string(buf)
}

func appendInt(attr string, v int64) string {
	buf := []byte(attr)
	buf = append(buf, equal, quote)
	buf = strconv.AppendInt(buf, v, 10)
	buf = append(buf, quote)
	return string(buf)
}

func appendFloat(attr string, v float64) string {
	buf := []byte(attr)
	buf = append(buf, equal, quote)
	buf = strconv.AppendFloat(buf, v, 'f', getPrecision(v), 64)
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

func appendFloatArray(attr string, list []float64, sep byte) string {
	buf := []byte(attr)
	buf = append(buf, equal, quote)
	for i := range list {
		if i > 0 {
			buf = append(buf, sep)
		}
		buf = strconv.AppendFloat(buf, list[i], 'f', getPrecision(list[i]), 64)
	}
	buf = append(buf, quote)
	return string(buf)
}

func appendFloatPair(attr string, list []float64) string {
	buf := []byte(attr)
	buf = append(buf, equal, quote)
	for i := 0; i < len(list); i += 2 {
		if i > 0 {
			buf = append(buf, space)
		}
		buf = strconv.AppendFloat(buf, list[i], 'f', getPrecision(list[i]), 64)
		buf = append(buf, comma)
		buf = strconv.AppendFloat(buf, list[i+1], 'f', getPrecision(list[i+1]), 64)
	}
	buf = append(buf, quote)
	return string(buf)
}

func getPrecision(f float64) int {
	if math.Ceil(f) == f {
		return 0
	}
	return defaultPrecision
}
