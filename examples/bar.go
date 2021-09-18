package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/midbel/svg"
	"github.com/midbel/svg/chart"
	"github.com/midbel/svg/colors"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	var c chart.BarChart
	c.Title = "bar chart demo"
	c.Width = 800
	c.Height = 640
	c.Padding = chart.CreatePadding(60, 40)
	c.Axis.Bottom = chart.CreateLabelAxis()
	c.Axis.Left = chart.CreateNumberAxis(chart.WithTicks(10))

	var fill []svg.Fill
	for i := range colors.RdYlBu4 {
		fill = append(fill, svg.NewFill(colors.RdYlBu4[i]))
	}

	c.Render(os.Stdout, getSeries(fill))
}

func getSeries(fills []svg.Fill) []chart.BarSerie {
	r1 := chart.NewBarSerie("rustine")
	r1.Fill = append(r1.Fill, fills...)
	r1.Add("code", randValue())
	r1.Add("bug", randValue())
	r1.Add("ticket", randValue())
	r1.Add("repo", randValue())

	r2 := chart.NewBarSerie("midbel")
	r2.Fill = append(r1.Fill, fills...)
	r2.Add("code", randValue())
	r2.Add("bug", randValue())
	r2.Add("ticket", randValue())
	r2.Add("repo", randValue())

	r3 := chart.NewBarSerie("hadock")
	r3.Fill = append(r1.Fill, fills...)
	r3.Add("code", randValue())
	r3.Add("bug", randValue())
	r3.Add("ticket", randValue())
	r3.Add("repo", randValue())

	r4 := chart.NewBarSerie("assist")
	r4.Fill = append(r1.Fill, fills...)
	r4.Add("code", randValue())
	r4.Add("bug", randValue())
	r4.Add("ticket", randValue())
	r4.Add("repo", randValue())

	return []chart.BarSerie{r1, r2, r3, r4}
}

func randValue() float64 {
	i := rand.Intn(100)
	return float64(i)
}
