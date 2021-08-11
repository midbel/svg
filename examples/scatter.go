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
	c.Width = 640
	c.Height = 480
	c.TicksY = 10
	c.TicksX = 5

	var sr chart.LineSerie
	for i := -100; i < 100; i += 5 {
		sr.Add(float64(i), float64(-20+rand.Intn(41)))
	}
	c.Render(os.Stdout, sr)
}
