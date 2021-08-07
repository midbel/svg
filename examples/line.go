package main

import (
	"math/rand"
	"os"

	"github.com/midbel/svg/chart"
)

const limit = 100

func main() {
	var xs []chart.LineSerie
	for i := 0; i < 2; i++ {
		var sr chart.LineSerie
		for i := 0; i < 20; i++ {
			sr.Add(float64(i), float64(5+rand.Intn(limit)))
		}
		xs = append(xs, sr)
	}
	var c chart.LineChart
	c.Padding = chart.Padding{
		Top:    20,
		Left:   60,
		Bottom: 60,
		Right:  30,
	}
	c.Width = 960
	c.Height = 720
	c.TicksY = 10
	c.TicksX = 5

	c.Render(os.Stdout, xs)
}
