package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"

	"github.com/midbel/svg"
	"github.com/midbel/svg/chart"
)

const limit = 200

func main() {
	var (
		c1 = getChart(chart.LinearCurve(), "salmon", true)
		c2 = getChart(chart.CubicCurve(0.5), "olive", true)
		c3 = getChart(chart.StepBeforeCurve(), "steelblue", false)
		c4 = getChart(chart.StepAfterCurve(), "orchid", true)
		c5 = getChart(chart.QuadraticCurve(0.5), "orange", false)
		c6 = getChart(chart.StepCurve(), "teal", true)
	)
	area := svg.NewSVG(svg.WithDimension(1920, 480*3))
	gp0 := svg.NewGroup()
	gp0.Append(multiserie())
	area.Append(gp0.AsElement())
	gp1 := svg.NewGroup(svg.WithTranslate(0, 480))
	gp1.Append(c1)
	area.Append(gp1.AsElement())
	gp2 := svg.NewGroup(svg.WithTranslate(640, 480))
	gp2.Append(c2)
	area.Append(gp2.AsElement())
	gp5 := svg.NewGroup(svg.WithTranslate(1280, 480))
	gp5.Append(c5)
	area.Append(gp5.AsElement())
	gp3 := svg.NewGroup(svg.WithTranslate(0, 960))
	gp3.Append(c3)
	area.Append(gp3.AsElement())
	gp4 := svg.NewGroup(svg.WithTranslate(640, 960))
	gp4.Append(c4)
	area.Append(gp4.AsElement())
	gp6 := svg.NewGroup(svg.WithTranslate(1280, 960))
	gp6.Append(c6)
	area.Append(gp6.AsElement())

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	area.Render(w)
}

func multiserie() svg.Element {
	sr1 := getSerie(chart.LinearCurve(), "red")
	sr1.Title = "blue"
	sr1.YAxis = chart.CreateNumberAxis(
		chart.WithTicks(7, true, false),
		chart.WithPosition(chart.Left),
		chart.WithNumberRange(-100, 100),
	)
	sr2 := getSerie(chart.LinearCurve(), "blue")
	sr2.Title = "blue"
	sr2.YAxis = chart.CreateNumberAxis(
		chart.WithTicks(4, true, false),
		chart.WithPosition(chart.Right),
	)

	var c chart.LineChart
	c.Title = "multi-axis red/blue chart"
	c.Padding = chart.CreatePadding(60, 60)
	c.Width = 1920
	c.Height = 480
	c.XAxis = chart.CreateNumberAxis(chart.WithTicks(15, true, true), chart.WithPosition(chart.Bottom))
	return c.RenderElement([]chart.LineSerie{sr1, sr2})
}

func getSerie(curve chart.Curver, color string) chart.LineSerie {
	sr := chart.NewLineSerie(color)
	sr.Curver = curve
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

func getChart(curve chart.Curver, color string, fill bool) svg.Element {
	var (
		c chart.LineChart
		s = getSerie(curve, color)
		f = svg.NewFill(color)
	)
	f.Opacity = 0.6
	if fill {
		s.Fill = f
	}
	c.Title = fmt.Sprintf("line serie (fill: %s)", color)
	c.Padding = chart.Padding{
		Top:    30,
		Left:   60,
		Bottom: 60,
		Right:  10,
	}
	c.Width = 640
	c.Height = 480
	c.YAxis = chart.CreateNumberAxis(chart.WithTicks(7, true, true), chart.WithPosition(chart.Left))
	c.XAxis = chart.CreateNumberAxis(chart.WithTicks(5, true, true), chart.WithPosition(chart.Bottom))
	return c.RenderElement([]chart.LineSerie{s})
}
