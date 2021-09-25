package chart

import (
	"math"

	"github.com/midbel/svg"
)

type CurveStyle uint8

const (
	CurveLinear CurveStyle = iota
	CurveStep
	CurveStepBefore
	CurveStepAfter
	CurveCubic
	CurveQuadratic
)

type Curver interface {
	Draw(Chart, XYSerie, Pair, Pair) svg.Element
}

type CurveFunc func(Chart, XYSerie, Pair, Pair) svg.Element

func (c CurveFunc) Draw(ch Chart, serie XYSerie, px, py Pair) svg.Element {
	return c(ch, serie, px, py)
}

type linearcurve struct {
	width  float64
	height float64
}

func LinearCurve() Curver {
	return linearcurve{}
}

func (c linearcurve) Draw(ch Chart, serie XYSerie, px, py Pair) svg.Element {
	var (
		dx  = ch.GetAreaWidth() / px.Diff()
		dy  = ch.GetAreaHeight() / py.Diff()
		pat = svg.NewPath(serie.GetStroke().Option(), nonefill.Option())
		pos svg.Pos
		ori svg.Pos
		pt  = serie.At(0)
	)
	pos.X = (pt.X - px.First()) * dx
	pos.Y = ch.GetAreaHeight() - (pt.Y * dy)
	if py.First() < 0 {
		pos.Y -= math.Abs(py.First()) * dy
	}
	pat.AbsMoveTo(pos)
	ori = pos
	for i := 1; i < serie.Len(); i++ {
		pt = serie.At(i)

		pos.X = (pt.X - px.First()) * dx
		pos.Y = ch.GetAreaHeight() - (pt.Y * dy)
		if py.First() < 0 {
			pos.Y -= math.Abs(py.First()) * dy
		}
		pat.AbsLineTo(pos)
	}
	return closePath(ch, pos, ori, serie.GetFill(), pat)
}

type stepaftercurve struct {
	width  float64
	height float64
}

func StepAfterCurve() Curver {
	return stepaftercurve{}
}

func (c stepaftercurve) Draw(ch Chart, serie XYSerie, px, py Pair) svg.Element {
	var (
		dx  = ch.GetAreaWidth() / px.Diff()
		dy  = ch.GetAreaHeight() / py.Diff()
		pat = svg.NewPath(serie.GetStroke().Option(), nonefill.Option())
		pos svg.Pos
		ori svg.Pos
		pt  = serie.At(0)
	)
	pos.X = (pt.X - px.First()) * dx
	pos.Y = ch.GetAreaHeight() - (pt.Y * dy)
	if py.First() < 0 {
		pos.Y -= math.Abs(py.First()) * dy
	}
	pat.AbsMoveTo(pos)
	ori = pos
	for i := 1; i < serie.Len(); i++ {
		pt = serie.At(i)

		pos.X += (pt.X - serie.At(i-1).X) * dx
		pat.AbsLineTo(pos)

		pos.Y = ch.GetAreaHeight() - (pt.Y * dy)
		if py.First() < 0 {
			pos.Y -= math.Abs(py.First()) * dy
		}
		pat.AbsLineTo(pos)
	}
	return closePath(ch, pos, ori, serie.GetFill(), pat)
}

type stepbeforecurve struct{}

func StepBeforeCurve() Curver {
	return stepbeforecurve{}
}

func (c stepbeforecurve) Draw(ch Chart, serie XYSerie, px, py Pair) svg.Element {
	var (
		dx  = ch.GetAreaWidth() / px.Diff()
		dy  = ch.GetAreaHeight() / py.Diff()
		pat = svg.NewPath(serie.GetStroke().Option(), nonefill.Option())
		pos svg.Pos
		ori svg.Pos
		pt  = serie.At(0)
	)
	pos.X = (pt.X - px.First()) * dx
	pos.Y = ch.GetAreaHeight() - (pt.Y * dy)
	if py.First() < 0 {
		pos.Y -= math.Abs(py.First()) * dy
	}
	pat.AbsMoveTo(pos)
	ori = pos
	for i := 1; i < serie.Len(); i++ {
		pt = serie.At(i)
		pos.Y = ch.GetAreaHeight() - (pt.Y * dy)
		if py.First() < 0 {
			pos.Y -= math.Abs(py.First()) * dy
		}
		pat.AbsLineTo(pos)
		pos.X += (pt.X - serie.At(i-1).X) * dx
		pat.AbsLineTo(pos)
	}
	return closePath(ch, pos, ori, serie.GetFill(), pat)
}

type stepcurve struct{}

func StepCurve() Curver {
	return stepcurve{}
}

func (c stepcurve) Draw(ch Chart, serie XYSerie, px, py Pair) svg.Element {
	var (
		dx  = ch.GetAreaWidth() / px.Diff()
		dy  = ch.GetAreaHeight() / py.Diff()
		pat = svg.NewPath(serie.GetStroke().Option(), nonefill.Option())
		pos svg.Pos
		ori svg.Pos
		pt  = serie.At(0)
	)
	pos.X = (pt.X - px.First()) * dx
	pos.Y = ch.GetAreaHeight() - (pt.Y * dy)
	if py.First() < 0 {
		pos.Y -= math.Abs(py.First()) * dy
	}
	pat.AbsMoveTo(pos)
	ori = pos
	for i := 1; i < serie.Len(); i++ {
		pt = serie.At(i)
		delta := (pt.X - serie.At(i-1).X) / 2
		pos.X += delta * dx
		pat.AbsLineTo(pos)

		pos.Y = ch.GetAreaHeight() - (pt.Y * dy)
		if py.First() < 0 {
			pos.Y -= math.Abs(py.First()) * dy
		}
		pat.AbsLineTo(pos)
		pos.X += delta * dx
		pat.AbsLineTo(pos)
	}
	return closePath(ch, pos, ori, serie.GetFill(), pat)
}

type cubiccurve struct {
	stretch float64
}

func CubicCurve(stretch float64) Curver {
	return cubiccurve{
		stretch: stretch,
	}
}

func (c cubiccurve) Draw(ch Chart, serie XYSerie, px, py Pair) svg.Element {
	var (
		dx  = ch.GetAreaWidth() / px.Diff()
		dy  = ch.GetAreaHeight() / py.Diff()
		pat = svg.NewPath(serie.GetStroke().Option(), nonefill.Option())
		pos svg.Pos
		ori svg.Pos
		pt  = serie.At(0)
	)

	pos.X = (pt.X - px.First()) * dx
	pos.Y = ch.GetAreaHeight() - (pt.Y * dy)
	if py.First() < 0 {
		pos.Y -= math.Abs(py.First()) * dy
	}
	pat.AbsMoveTo(pos)
	ori = pos
	for i := 1; i < serie.Len(); i++ {
		var (
			ctrl = pos
			old  = pos
		)
		pt = serie.At(i)
		pos.Y = ch.GetAreaHeight() - (pt.Y * dy)
		if py.First() < 0 {
			pos.Y -= math.Abs(py.First()) * dy
		}
		pos.X = (pt.X - px.First()) * dx
		ctrl.X = old.X - (old.X-pos.X)*c.stretch
		ctrl.Y = pos.Y
		pat.AbsCubicCurveSimple(pos, ctrl)
	}
	return closePath(ch, pos, ori, serie.GetFill(), pat)
}

type quadraticcurve struct {
	stretch float64
}

func QuadraticCurve(stretch float64) Curver {
	return quadraticcurve{
		stretch: stretch,
	}
}

func (c quadraticcurve) Draw(ch Chart, serie XYSerie, px, py Pair) svg.Element {
	var (
		dx   = ch.GetAreaWidth() / px.Diff()
		dy   = ch.GetAreaHeight() / py.Diff()
		pat  = svg.NewPath(serie.GetStroke().Option(), nonefill.Option())
		pos  svg.Pos
		ori  svg.Pos
		old  svg.Pos
		ctrl svg.Pos
		pt   = serie.At(0)
	)
	pos.X = (pt.X - px.First()) * dx
	pos.Y = ch.GetAreaHeight() - (pt.Y * dy)
	if py.First() < 0 {
		pos.Y -= math.Abs(py.First()) * dy
	}
	pat.AbsMoveTo(pos)
	ori = pos
	for i := 1; i < serie.Len(); i++ {
		pt = serie.At(i)

		old = pos
		pos.X = (pt.X - px.First()) * dx
		pos.Y = ch.GetAreaHeight() - (pt.Y * dy)
		if py.First() < 0 {
			pos.Y -= math.Abs(py.First()) * dy
		}
		ctrl.X = old.X
		ctrl.Y = pos.Y
		pat.AbsQuadraticCurve(pos, ctrl)
	}
	return closePath(ch, pos, ori, serie.GetFill(), pat)
}

func closePath(ch Chart, pos, end svg.Pos, fill svg.Fill, pat svg.Path) svg.Element {
	if !fill.IsZero() {
		pos.Y = ch.GetAreaHeight()
		pat.AbsLineTo(pos)
		pos.X = end.X
		pat.AbsLineTo(pos)
		pat.AbsLineTo(end)
		pat.Fill = fill
	}
	return pat.AsElement()
}
