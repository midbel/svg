package svg

import (
	"math"
	"strconv"
	"strings"
)

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
