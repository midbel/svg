package main

import (
	"math/rand"
	"os"

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
	c.Radius = 7
	c.Shape = chart.ShapeStar
	c.Width = 640
	c.Height = 480
	c.TicksY = 10
	c.TicksX = 5

	sr1 := chart.NewLineSerieWithColor("triangle", "salmon")
	sr1.Shape = chart.ShapeTriangle
	for i := -100; i < 100; i += 5 {
		sr1.Add(float64(i), float64(-20+rand.Intn(41)))
	}
	sr2 := chart.NewLineSerieWithColor("diamond", "steelblue")
	sr2.Shape = chart.ShapeDiamond
	for i := -100; i < 100; i += 5 {
		sr2.Add(float64(i), float64(40+rand.Intn(21)))
	}
	c.Render(os.Stdout, []chart.LineSerie{sr1, sr2})
}
