package svg

import (
	"math"
	"strconv"
	"strings"
)

const defaultFontSize = 14

var (
	DefaultStroke = NewStroke("black", 1)
	DefaultFill   = NewFill("black")
	DefaultFont   = NewFont(defaultFontSize)
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
		{Attr: "fill", Value: f.Fill},
		{Attr: "font-style", Value: f.Style},
		{Attr: "font-weight", Value: f.Weight},
		{Attr: "font-variant", Value: f.Variant},
		{Attr: "font-stretch", Value: f.Stretch},
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
	X    float64
	Y    float64
	unit string
}

func NewPos(x, y float64) Pos {
	return Pos{
		X: x,
		Y: y,
	}
}

func (p Pos) Attributes() []string {
	var attrs []string
	attrs = append(attrs, appendFloat("x", p.X))
	attrs = append(attrs, appendFloat("y", p.Y))
	return attrs
}

func (p Pos) Center() []string {
	var attrs []string
	attrs = append(attrs, appendFloat("cx", p.X))
	attrs = append(attrs, appendFloat("cy", p.Y))
	return attrs
}

func (p Pos) InPercent() Pos {
	p.unit = UnitPer
	return p
}

func (p Pos) InPX() Pos {
	p.unit = UnitPX
	return p
}

func (p Pos) InEM() Pos {
	p.unit = UnitEM
	return p
}

func (p Pos) InEX() Pos {
	p.unit = UnitEX
	return p
}

func (p Pos) InPT() Pos {
	p.unit = UnitPT
	return p
}

func (p Pos) InPC() Pos {
	p.unit = UnitPC
	return p
}

func (p Pos) InCM() Pos {
	p.unit = UnitCM
	return p
}

func (p Pos) InMM() Pos {
	p.unit = UnitMM
	return p
}

func (p Pos) InIN() Pos {
	p.unit = UnitIN
	return p
}

func (p Pos) array() []float64 {
	return []float64{p.X, p.Y}
}

type Dim struct {
	W    float64
	H    float64
	unit string
}

func NewDim(w, h float64) Dim {
	return Dim{
		W: w,
		H: h,
	}
}

func (d Dim) Attributes() []string {
	var attrs []string
	attrs = append(attrs, appendFloat("width", d.W))
	attrs = append(attrs, appendFloat("height", d.H))
	return attrs
}

func (d Dim) InPercent() Dim {
	d.unit = UnitPer
	return d
}

func (d Dim) InPX() Dim {
	d.unit = UnitPX
	return d
}

func (d Dim) InEM() Dim {
	d.unit = UnitEM
	return d
}

func (d Dim) InEX() Dim {
	d.unit = UnitEX
	return d
}

func (d Dim) InPT() Dim {
	d.unit = UnitPT
	return d
}

func (d Dim) InPC() Dim {
	d.unit = UnitPC
	return d
}

func (d Dim) InCM() Dim {
	d.unit = UnitCM
	return d
}

func (d Dim) InMM() Dim {
	d.unit = UnitMM
	return d
}

func (d Dim) InIN() Dim {
	d.unit = UnitIN
	return d
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
	Dash struct {
		Array  []int
		Offset []int
	}
	Line struct {
		Cap  string
		Join string
	}
	Width   float64
	Opacity float64
	Miter   float64
	Fill    string
}

func NewStroke(fill string, width int) Stroke {
	return Stroke{
		Fill:  fill,
		Width: float64(width),
	}
}

func (s Stroke) Attributes() []string {
	var attrs []string
	if len(s.Dash.Array) > 0 {
		attrs = append(attrs, appendIntArray("stroke-dasharray", s.Dash.Array, space))
	}
	if len(s.Dash.Offset) > 0 {
		attrs = append(attrs, appendIntArray("stroke-dashoffset", s.Dash.Offset, space))
	}
	if s.Line.Cap != "" {
		attrs = append(attrs, appendString("stroke-linecap", s.Line.Cap))
	}
	if s.Line.Join != "" {
		attrs = append(attrs, appendString("stroke-linejoin", s.Line.Join))
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
	if s.Fill != "" {
		attrs = append(attrs, appendString("stroke", s.Fill))
	}
	return attrs
}

type Fill struct {
	Color   string
	Rule    string
	Opacity float64
}

func NewFill(color string) Fill {
	return Fill{Color: color, Opacity: 100}
}

func (f Fill) Attributes() []string {
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
