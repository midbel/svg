package main

import (
	"math/rand"
	"os"

	"github.com/midbel/svg"
	"github.com/midbel/svg/chart"
)

func main() {
	var (
		c  chart.PolarChart
		s1 = getSerie("olive")
		// s2 = getSerie()
		// s3 = getSerie()
		// s4 = getSerie()
	)
	c.Padding = chart.CreatePadding(20, 20)
	c.Width = 800
	c.Height = 800
	c.Zone = 5

	c.Render(os.Stdout, s1)
}

func getSerie(color string) chart.PolarSerie {
	var sr chart.PolarSerie
	sr.Fill = svg.NewFill(color)
	sr.Radius = 6
	for i := 0; i < 10; i++ {
		sr.Add(float64(-20 + rand.Intn(41)))
	}
	return sr
}
