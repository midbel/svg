package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/midbel/svg"
	"github.com/midbel/svg/chart"
)

func main() {
	var c chart.TimeChart
	c.Width = 1200
	c.Height = 800
	// c.Title = "time chart demo"
	c.Axis.Bottom = chart.CreateTimeAxis(chart.WithTicks(10))
	// c.Axis.Top = chart.CreateTimeAxis(chart.WithTicks(10))
	c.Axis.Left = chart.CreateNumberAxis(chart.WithTicks(12))
	// c.Axis.Right = chart.CreateNumberAxis(chart.WithTicks(4))
	c.Padding = chart.Padding{
		Left:   80,
		Right:  60,
		Bottom: 40,
		Top:    40,
	}
	sr1 := getSerie()
	c.Render(os.Stdout, []chart.TimeSerie{sr1})
}

func getSerie() chart.TimeSerie {
	var (
		serie = chart.NewTimeSerie("time serie")
		now   = time.Now()
		delta = time.Hour * 24
	)
	serie.Stroke = svg.NewStroke("olive", 1)
	for i := 0; i < 100; i++ {
		c := rand.Intn(5)
		if c == 0 {
			continue
		}
		now = now.Add(time.Duration(c) * delta)
		serie.Add(now, float64(-100+rand.Intn(200)))
	}
	return serie
}
