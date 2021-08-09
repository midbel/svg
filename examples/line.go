package main

import (
	"fmt"
	"math/rand"
	"os"
	// "time"

	"github.com/midbel/svg/chart"
)

const limit = 200

// func init() {
// 	rand.Seed(time.Now().Unix())
// }

func main() {
	var (
		xs []chart.LineSerie
		cs = []string{"red", "blue", "green"}
	)
	for i := 0; i < 1; i++ {
		var (
			s    = fmt.Sprint("serie-%d", i)
			sr   = chart.NewLineSerieWithColor(s, cs[i])
			step = 5 + rand.Intn(5)
		)
		for i := -100; i <= 100; i += step {
			sr.Add(float64(i), float64(-20+rand.Intn(41)))
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
	c.Curve = chart.CurveStepBefore
	c.Width = 1516
	c.Height = 770
	c.TicksY = 15
	c.TicksX = 15

	c.Render(os.Stdout, xs)
}
