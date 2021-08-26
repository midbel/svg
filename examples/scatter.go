package main

import (
	"math/rand"
	"os"

	"github.com/midbel/svg"
	"github.com/midbel/svg/chart"
)

var colors = []string{"salmon", "steelblue", "olive", "orchid"}

func main() {
	var c chart.ScatterChart
	c.Padding = chart.Padding{
		Left:   60,
		Right:  20,
		Bottom: 40,
		Top:    20,
	}
	c.Radius = 8
	c.Shape = chart.ShapeStar
	c.Width = 800
	c.Height = 600
	c.LineAxis = chart.NewLineAxisWith(10, true, true)
	c.GetColor = func(_ string, i int) svg.Fill {
		return svg.NewFill(colors[i%len(colors)])
	}

	sr1 := chart.NewLineSerie("triangle")
	sr1.Shape = chart.ShapeTriangle
	for i := -100; i < 100; i += 5 {
		sr1.Add(float64(i), float64(-20+rand.Intn(41)))
	}
	sr2 := chart.NewLineSerie("diamond")
	sr2.Shape = chart.ShapeDiamond
	for i := -100; i < 100; i += 5 {
		sr2.Add(float64(i), float64(40+rand.Intn(21)))
	}
	sr3 := chart.NewLineSerie("square")
	for i := -100; i < 100; i += 10 {
		sr3.Add(float64(i), float64(70+rand.Intn(61)))
	}
	sr4 := chart.NewLineSerie("star")
	sr4.Shape = chart.ShapeStar
	for i := 100; i < 200; i += 5 {
		sr4.Add(float64(i), float64(-20+rand.Intn(161)))
	}
	c.Render(os.Stdout, []chart.LineSerie{sr1, sr2, sr3, sr4})
}
