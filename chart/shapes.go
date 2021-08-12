package chart

import (
  "github.com/midbel/svg"
)

type ShapeType uint8

const (
	ShapeDefault ShapeType = iota
	ShapeSquare
	ShapeCircle
	ShapeTriangle
	ShapeStar
	ShapeDiamond
)

func (s ShapeType) Draw(rad float64, options ...svg.Option) svg.Element {
  var elem svg.Element
  switch s {
  case ShapeCircle:
    elem = getCircle(rad, options...)
  case ShapeTriangle:
    elem = getTriangle(rad, options...)
  case ShapeStar:
    elem = getStar(rad, options...)
  case ShapeDiamond:
    elem = getDiamond(rad, options...)
  case ShapeSquare:
    elem = getSquare(rad, options...)
  default:
  }
  return elem
}

func getDiamond(rad float64, options ...svg.Option) svg.Element {
	options = append(options, svg.WithDimension(rad, rad))
	i := svg.NewRect(options...)
	return i.AsElement()
}

func getSquare(rad float64, options ...svg.Option) svg.Element {
	options = append(options, svg.WithDimension(rad, rad))
	i := svg.NewRect(options...)
	return i.AsElement()
}

func getTriangle(rad float64, options ...svg.Option) svg.Element {
	points := []svg.Pos{
		svg.NewPos(0, rad),
		svg.NewPos(rad/2, 0),
		svg.NewPos(rad, rad),
	}
	i := svg.NewPolygon(points, options...)
	return i.AsElement()
}

func getCircle(rad float64, options ...svg.Option) svg.Element {
	options = append(options, svg.WithRadius(rad/2))
	i := svg.NewCircle(options...)
	return i.AsElement()
}

func getStar(rad float64, options ...svg.Option) svg.Element {
	rad *= 2
	var (
		onerad   = rad / 5
		tworad   = onerad * 2
		threerad = onerad * 3
		fourrad  = onerad * 4
		halfrad  = rad / 2
	)
	points := []svg.Pos{
		svg.NewPos(onerad, rad),
		svg.NewPos(tworad, halfrad),
		svg.NewPos(0, tworad),
		svg.NewPos(tworad, tworad),
		svg.NewPos(halfrad, 0),
		svg.NewPos(threerad, tworad),
		svg.NewPos(rad, tworad),
		svg.NewPos(threerad, halfrad),
		svg.NewPos(fourrad, rad),
		svg.NewPos(halfrad, threerad),
		svg.NewPos(onerad, rad),
	}
	i := svg.NewPolygon(points, options...)
	return i.AsElement()
}
