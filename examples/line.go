package main

import (
	"math/rand"
	"os"

	"github.com/midbel/svg/chart"
)

const limit = 200

func main() {
	var xs []chart.LineSerie
	for i := 0; i < 2; i++ {
		var sr chart.LineSerie
		for i := -1000; i <= 1000; i += 50 {
			sr.Add(float64(i), float64(-100+rand.Intn(limit)))
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
	c.Width = 1516
	c.Height = 770
	c.TicksY = 15
	c.TicksX = 15

	c.Render(os.Stdout, xs)
}
