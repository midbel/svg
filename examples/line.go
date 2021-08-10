package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/midbel/svg"
	"github.com/midbel/svg/chart"
)

const limit = 200

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	var (
		xs []chart.LineSerie
		cs = []string{"red", "blue", "green"}
	)
	for i := 0; i < 2; i++ {
		var (
			s  = fmt.Sprint("serie-%d", i)
			sr = chart.NewLineSerieWithColor(s, cs[i])
		)
		for i := -100; i <= 100; i += 5 + rand.Intn(5) {
			sr.Add(float64(i), float64(-15+rand.Intn(200)))
		}
		xs = append(xs, sr)
	}

	var (
		c1 = getChart(chart.CurveLinear, xs)
		c2 = getChart(chart.CurveCubic, xs)
		c3 = getChart(chart.CurveStepBefore, xs)
		c4 = getChart(chart.CurveStepAfter, xs)
		c5 = getChart(chart.CurveQuadratic, xs)
		c6 = getChart(chart.CurveStep, xs)
	)
	area := svg.NewSVG(svg.WithDimension(1440, 720))
	gp1 := svg.NewGroup(svg.WithTranslate(0, 0))
	gp1.Append(c1)
	area.Append(gp1.AsElement())
	gp2 := svg.NewGroup(svg.WithTranslate(480, 0))
	gp2.Append(c2)
	area.Append(gp2.AsElement())
	gp3 := svg.NewGroup(svg.WithTranslate(0, 360))
	gp3.Append(c3)
	area.Append(gp3.AsElement())
	gp4 := svg.NewGroup(svg.WithTranslate(480, 360))
	gp4.Append(c4)
	area.Append(gp4.AsElement())
	gp5 := svg.NewGroup(svg.WithTranslate(960, 0))
	gp5.Append(c5)
	area.Append(gp5.AsElement())
	gp6 := svg.NewGroup(svg.WithTranslate(960, 360))
	gp6.Append(c6)
	area.Append(gp6.AsElement())

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	area.Render(w)
}

func getChart(curve chart.CurveStyle, data []chart.LineSerie) svg.Element {
	var c chart.LineChart
	c.Padding = chart.Padding{
		Top:    20,
		Left:   60,
		Bottom: 60,
		Right:  30,
	}
	c.Curve = curve
	c.Width = 480
	c.Height = 360
	c.TicksY = 7
	c.TicksX = 7

	e := c.RenderElement(data)
	if e, ok := e.(*svg.SVG); ok {
		e.OmitProlog = ok
	}
	return e
}
