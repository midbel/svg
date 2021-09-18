package main

import (
	"bufio"
	"math/rand"
	"os"
	// "time"

	"github.com/midbel/svg"
	"github.com/midbel/svg/chart"
)

const limit = 200

// func init() {
// 	rand.Seed(time.Now().Unix())
// }

func main() {
	var (
		c1 = getChart(chart.CurveLinear, "salmon")
		c2 = getChart(chart.CurveCubic, "olive")
		c3 = getChart(chart.CurveStepBefore, "steelblue")
		c4 = getChart(chart.CurveStepAfter, "orchid")
		c5 = getChart(chart.CurveQuadratic, "orange")
		c6 = getChart(chart.CurveStep, "teal")
	)
	area := svg.NewSVG(svg.WithDimension(1920, 960))
	gp1 := svg.NewGroup(svg.WithTranslate(0, 0))
	gp1.Append(c1)
	area.Append(gp1.AsElement())
	gp2 := svg.NewGroup(svg.WithTranslate(640, 0))
	gp2.Append(c2)
	area.Append(gp2.AsElement())
	gp3 := svg.NewGroup(svg.WithTranslate(0, 480))
	gp3.Append(c3)
	area.Append(gp3.AsElement())
	gp4 := svg.NewGroup(svg.WithTranslate(640, 480))
	gp4.Append(c4)
	area.Append(gp4.AsElement())
	gp5 := svg.NewGroup(svg.WithTranslate(1280, 0))
	gp5.Append(c5)
	area.Append(gp5.AsElement())
	gp6 := svg.NewGroup(svg.WithTranslate(1280, 480))
	gp6.Append(c6)
	area.Append(gp6.AsElement())

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	area.Render(w)
}

func getSerie(curve chart.CurveStyle, color string) chart.LineSerie {
	sr := chart.NewLineSerie(color)
	sr.Curve = curve
	sr.Stroke = svg.NewStroke(color, 2)
	for i := -100; i < 100; i++ {
		c := rand.Intn(10)
		if c == 0 {
			continue
		}
		i += c
		sr.Add(float64(i), float64(-20+rand.Intn(41)))
	}
	return sr
}

func getChart(curve chart.CurveStyle, color string) svg.Element {
	var (
		c chart.LineChart
		s = getSerie(curve, color)
		f = svg.NewFill(color)
	)
	f.Opacity = 0.6
	s.Fill = f
	c.Padding = chart.Padding{
		Top:    10,
		Left:   60,
		Bottom: 60,
		Right:  10,
	}
	c.Width = 640
	c.Height = 480
	c.Axis.Left = chart.CreateNumberAxis(chart.WithTicks(7))
	c.Axis.Bottom = chart.CreateNumberAxis(chart.WithTicks(5))
	return c.RenderElement([]chart.LineSerie{s})
}
