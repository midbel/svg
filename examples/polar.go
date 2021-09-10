package main

import (
	"math/rand"
	"os"

	"github.com/midbel/svg/chart"
)

func main() {
	var (
		c  chart.PolarChart
		s1 = getSerie()
		s2 = getSerie()
		s3 = getSerie()
		s4 = getSerie()
	)
	c.Padding = chart.CreatePadding(20, 20)
	c.Width = 800
	c.Height = 800
	c.Zone = 5

	c.Render(os.Stdout, []chart.PolarSerie{s1, s2, s3, s4})
}

func getSerie() chart.PolarSerie {
	var sr chart.PolarSerie
	for i := 0; i < 5; i++ {
		c := rand.Intn(10)
		sr.Add(float64(i+c), float64(-20+rand.Intn(41)))
	}
	return sr
}
