package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/midbel/svg"
	"github.com/midbel/svg/chart"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	var c chart.AreaChart
	c.Width = 800
	c.Height = 640
	c.Padding = chart.CreatePadding(60, 40)
	c.Title = "area chart demo"
	c.Axis.Bottom = chart.CreateNumberAxis(chart.WithTicks(10))
	c.Axis.Left = chart.CreateNumberAxis(chart.WithTicks(10))

	a := chart.NewAreaSerie("serie", getSerie("1", 20), getSerie("2", 40))
	a.Fill = svg.NewFill("steelblue")
	a.Fill.Opacity = 0.6
	a.Stroke = svg.NewStroke("steelblue", 2)

	c.Render(os.Stdout, a)
}

func getSerie(title string, limit int) chart.LineSerie {
	s := chart.NewLineSerie(title)
	for i := -100; i < 100; i += 1 + rand.Intn(10) {
		c := rand.Intn(10)
		if c == 0 {
			continue
		}
		i += c
		s.Add(float64(i), float64(limit-rand.Intn(10)))
	}
	return s
}

func randValue() float64 {
	i := rand.Intn(100)
	return float64(i)
}
