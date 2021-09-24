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

func (c CurveStyle) Curve(w, h float64) Curver {
  switch c {
  case CurveLinear:
    return linear(w, h)
  case CurveStep:
    return step(w, h)
  case CurveStepBefore:
    return stepbefore(w, h)
  case CurveStepAfter:
    return stepafter(w, h)
  case CurveCubic:
    return cubic(w, h)
  case CurveQuadratic:
    return quadratic(w, h)
  }
  return nil
}

type Curver interface {
  Draw(Appender, LineSerie, pair, pair)
}

type linearcurve struct {
  width float64
  height float64
}

func linear(w, h float64) Curver {
  return linearcurve{
    width: w,
    height: h,
  }
}

func (c linearcurve) Draw(ap Appender, serie LineSerie, px, py pair) {
  var (
		dx  = c.width / px.Diff()
		dy  = c.height / py.Diff()
		pat = svg.NewPath(serie.Stroke.Option(), nonefill.Option())
		pos svg.Pos
		ori svg.Pos
	)
	pos.X = (serie.values[0].X - px.Min) * dx
	pos.Y = c.height - (serie.values[0].Y * dy)
	if py.Min < 0 {
		pos.Y -= math.Abs(py.Min) * dy
	}
	pat.AbsMoveTo(pos)
	ori = pos
	for i := 1; i < serie.Len(); i++ {
		pos.X = (serie.values[i].X - px.Min) * dx
		pos.Y = c.height - (serie.values[i].Y * dy)
		if py.Min < 0 {
			pos.Y -= math.Abs(py.Min) * dy
		}
		pat.AbsLineTo(pos)
	}
	if !serie.Fill.IsZero() {
		pos.Y = c.height
		pat.AbsLineTo(pos)
		pos.X = ori.X
		pat.AbsLineTo(pos)
		pat.AbsLineTo(ori)
		pat.Fill = serie.Fill
	}
  ap.Append(pat.AsElement())
}

type stepaftercurve struct {
  width float64
  height float64
}

func stepafter(w, h float64) Curver {
  return stepaftercurve{
    width: w,
    height: h,
  }
}

func (c stepaftercurve) Draw(ap Appender, serie LineSerie, px, py pair) {
  var (
    dx  = c.width / px.Diff()
    dy  = c.height / py.Diff()
    pat = svg.NewPath(serie.Stroke.Option(), nonefill.Option())
    pos svg.Pos
    ori svg.Pos
  )
  pos.X = (serie.values[0].X - px.Min) * dx
  pos.Y = c.height - (serie.values[0].Y * dy)
  if py.Min < 0 {
    pos.Y -= math.Abs(py.Min) * dy
  }
  pat.AbsMoveTo(pos)
  ori = pos
  for i := 1; i < serie.Len(); i++ {
    pos.X += (serie.values[i].X - serie.values[i-1].X) * dx
    pat.AbsLineTo(pos)

    pos.Y = c.height - (serie.values[i].Y * dy)
    if py.Min < 0 {
      pos.Y -= math.Abs(py.Min) * dy
    }
    pat.AbsLineTo(pos)
  }
  if !serie.Fill.IsZero() {
    pos.Y = c.height
    pat.AbsLineTo(pos)
    pos.X = ori.X
    pat.AbsLineTo(pos)
    pat.AbsLineTo(ori)
    pat.Fill = serie.Fill
  }
  ap.Append(pat.AsElement())
}

type stepbeforecurve struct {
  width float64
  height float64
}

func stepbefore(w, h float64) Curver {
  return stepbeforecurve{
    width: w,
    height: h,
  }
}

func (c stepbeforecurve) Draw(ap Appender, serie LineSerie, px, py pair) {
  var (
		dx  = c.width / px.Diff()
		dy  = c.height / py.Diff()
		pat = svg.NewPath(serie.Stroke.Option(), nonefill.Option())
		pos svg.Pos
		ori svg.Pos
	)
	pos.X = (serie.values[0].X - px.Min) * dx
	pos.Y = c.height - (serie.values[0].Y * dy)
	if py.Min < 0 {
		pos.Y -= math.Abs(py.Min) * dy
	}
	pat.AbsMoveTo(pos)
	ori = pos
	for i := 1; i < serie.Len(); i++ {
		pos.Y = c.height - (serie.values[i].Y * dy)
		if py.Min < 0 {
			pos.Y -= math.Abs(py.Min) * dy
		}
		pat.AbsLineTo(pos)
		pos.X += (serie.values[i].X - serie.values[i-1].X) * dx
		pat.AbsLineTo(pos)
	}
	if !serie.Fill.IsZero() {
		pos.Y = c.height
		pat.AbsLineTo(pos)
		pos.X = ori.X
		pat.AbsLineTo(pos)
		pat.AbsLineTo(ori)
		pat.Fill = serie.Fill
	}
  ap.Append(pat.AsElement())
}

type stepcurve struct {
  width float64
  height float64
}

func step(w, h float64) Curver {
  return stepcurve{
    width: w,
    height: h,
  }
}

func (c stepcurve) Draw(ap Appender, serie LineSerie, px, py pair) {
  var (
    dx  = c.width / px.Diff()
    dy  = c.height / py.Diff()
    pat = svg.NewPath(serie.Stroke.Option(), nonefill.Option())
    pos svg.Pos
    ori svg.Pos
  )
  pos.X = (serie.values[0].X - px.Min) * dx
  pos.Y = c.height - (serie.values[0].Y * dy)
  if py.Min < 0 {
    pos.Y -= math.Abs(py.Min) * dy
  }
  pat.AbsMoveTo(pos)
  ori = pos
  for i := 1; i < serie.Len(); i++ {
    delta := (serie.values[i].X - serie.values[i-1].X) / 2
    pos.X += delta * dx
    pat.AbsLineTo(pos)

    pos.Y = c.height - (serie.values[i].Y * dy)
    if py.Min < 0 {
      pos.Y -= math.Abs(py.Min) * dy
    }
    pat.AbsLineTo(pos)
    pos.X += delta * dx
    pat.AbsLineTo(pos)
  }
  if !serie.Fill.IsZero() {
    pos.Y = c.height
    pat.AbsLineTo(pos)
    pos.X = ori.X
    pat.AbsLineTo(pos)
    pat.AbsLineTo(ori)
    pat.Fill = serie.Fill
  }
  ap.Append(pat.AsElement())
}


type cubiccurve struct {
  width float64
  height float64
  stretch float64
}

func cubic(w, h float64) Curver {
  return cubiccurve{
    width: w,
    height: h,
    stretch: 0.5,
  }
}

func (c cubiccurve) Draw(ap Appender, serie LineSerie, px, py pair) {
  var (
    dx  = c.width / px.Diff()
    dy  = c.height / py.Diff()
    pat = svg.NewPath(serie.Stroke.Option(), nonefill.Option())
    pos svg.Pos
    ori svg.Pos
  )
  pos.X = (serie.values[0].X - px.Min) * dx
  pos.Y = c.height - (serie.values[0].Y * dy)
  if py.Min < 0 {
    pos.Y -= math.Abs(py.Min) * dy
  }
  pat.AbsMoveTo(pos)
  ori = pos
  for i := 1; i < serie.Len(); i++ {
    var (
      ctrl = pos
      old  = pos
    )
    pos.Y = c.height - (serie.values[i].Y * dy)
    if py.Min < 0 {
      pos.Y -= math.Abs(py.Min) * dy
    }
    pos.X = (serie.values[i].X - px.Min) * dx
    ctrl.X = old.X - (old.X-pos.X)*c.stretch
    ctrl.Y = pos.Y
    pat.AbsCubicCurveSimple(pos, ctrl)
  }
  if !serie.Fill.IsZero() {
    pos.Y = c.height
    pat.AbsLineTo(pos)
    pos.X = ori.X
    pat.AbsLineTo(pos)
    pat.AbsLineTo(ori)
    pat.Fill = serie.Fill
  }
  ap.Append(pat.AsElement())
}

type quadraticcurve struct {
  width float64
  height float64
  stretch float64
}

func quadratic(w, h float64) Curver {
  return quadraticcurve{
    width: w,
    height: h,
    stretch: 0.5,
  }
}

func (c quadraticcurve) Draw(ap Appender, serie LineSerie, px, py pair) {
  var (
    dx   = c.width / px.Diff()
    dy   = c.height / py.Diff()
    pat  = svg.NewPath(serie.Stroke.Option(), nonefill.Option())
    pos  svg.Pos
    ori  svg.Pos
    old  svg.Pos
    ctrl svg.Pos
  )
  pos.X = (serie.values[0].X - px.Min) * dx
  pos.Y = c.height - (serie.values[0].Y * dy)
  if py.Min < 0 {
    pos.Y -= math.Abs(py.Min) * dy
  }
  pat.AbsMoveTo(pos)
  ori = pos
  for i := 1; i < serie.Len(); i++ {
    old = pos
    pos.X = (serie.values[i].X - px.Min) * dx
    pos.Y = c.height - (serie.values[i].Y * dy)
    if py.Min < 0 {
      pos.Y -= math.Abs(py.Min) * dy
    }
    ctrl.X = old.X
    ctrl.Y = pos.Y
    pat.AbsQuadraticCurve(pos, ctrl)
  }
  if !serie.Fill.IsZero() {
    pos.Y = c.height
    pat.AbsLineTo(pos)
    pos.X = ori.X
    pat.AbsLineTo(pos)
    pat.AbsLineTo(ori)
    pat.Fill = serie.Fill
  }
  ap.Append(pat.AsElement())
}
