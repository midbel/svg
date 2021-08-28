package main

import (
	"math/rand"
	"os"

	"github.com/midbel/svg"
	"github.com/midbel/svg/chart"
)

func main() {
	var c chart.ScatterChart
	c.Padding = chart.Padding{
		Left:   60,
		Right:  20,
		Bottom: 40,
		Top:    20,
	}
	c.Width = 800
	c.Height = 600
	c.LineAxis = chart.NewLineAxis(10, true, true)

	sr1 := chart.NewScatterSerie("triangle")
	sr1.Shape = chart.ShapeTriangle
	sr1.Size = 8
	sr1.Fill = svg.NewFill("salmon")
	for i := -100; i < 100; i += 5 {
		sr1.Add(float64(i), float64(-20+rand.Intn(41)))
	}
	sr2 := chart.NewScatterSerie("diamond")
	sr2.Shape = chart.ShapeDiamond
	sr2.Size = 8
	sr2.Fill = svg.NewFill("steelblue")
	sr2.Stroke = svg.NewStroke("lightblue", 1)
	sr2.Highlight = true
	for i := -100; i < 100; i += 5 {
		sr2.Add(float64(i), float64(40+rand.Intn(21)))
	}
	sr3 := chart.NewScatterSerie("square")
	sr3.Shape = chart.ShapeCircle
	sr3.Size = 8
	sr3.Fill = svg.NewFill("olive")
	for i := -100; i < 100; i += 10 {
		sr3.Add(float64(i), float64(70+rand.Intn(61)))
	}
	sr4 := chart.NewScatterSerie("star")
	sr4.Shape = chart.ShapeStar
	sr4.Size = 8
	sr4.Fill = svg.NewFill("orchid")
	for i := 100; i < 200; i += 5 {
		sr4.Add(float64(i), float64(-20+rand.Intn(161)))
	}
	c.Render(os.Stdout, []chart.ScatterSerie{sr1, sr2, sr3, sr4})
}
