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
	c.Padding = chart.Padding{
		Left:   80,
		Right:  60,
		Bottom: 40,
		Top:    20,
	}
	c.TimeAxis = chart.NewTimeAxis(7, true, true)
	c.TimeAxis.OuterX = true
	c.TimeAxis.OuterY = true
	c.TimeAxis.FormatTime = func(t time.Time, _ int) string {
		return t.Format("2006-01-02 15:04")
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
