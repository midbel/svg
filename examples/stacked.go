package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/midbel/svg"
	"github.com/midbel/svg/chart"
)

var pad chart.Padding

func init() {
	rand.Seed(time.Now().Unix())
}

var labels = []string{
	"xxh", "rustine", "cli", "svg",
	"cbor", "toml", "fig", "ber", "ldap",
	"fetch", "sdp", "ini", "ipaddr",
}

var axisFill = svg.NewStroke("lightgrey", 1)

func main() {
	var sc chart.StackedChart
	sc.Padding = pad
	var (
		count = flag.Int("c", 5, "count")
	)
	flag.Float64Var(&sc.Width, "x", chart.DefaultWidth, "width")
	flag.Float64Var(&sc.Height, "y", chart.DefaultHeight, "height")
	flag.Parse()
	var xs []chart.StackedSerie
	for i := 0; i < *count; i++ {
		vs := getSeries()
		sr := chart.NewStackedSerie(fmt.Sprintf("serie-%d", i))
		for _, s := range vs {
			sr.Append(s)
		}
		xs = append(xs, sr)
	}
	sc.Render(os.Stdout, xs)
}

func AxisX(xs []string, width, height float64) svg.Element {
	var (
		trst = svg.WithTranslate(pad.Left, height-pad.Bottom)
		grp  = svg.NewGroup(svg.WithID("x-axis"), trst)
		pos1 = svg.NewPos(0, 0)
		pos2 = svg.NewPos(width, 0)
		axis = svg.NewLine(pos1, pos2, axisFill.Option())
		step = width / float64(len(xs))
	)
	grp.Append(axis.AsElement())
	for i := range xs {
		g := svg.NewGroup()
		x := float64(i)*step + step/2
		y := 0.0
		pos := svg.NewPos(x-12, y+24)
		t := svg.NewText(xs[i], pos.Option())
		g.Append(t.AsElement())

		line := svg.NewLine(svg.NewPos(x, 0), svg.NewPos(x, 8), axisFill.Option())
		g.Append(line.AsElement())
		grp.Append(g.AsElement())
	}
	return grp.AsElement()
}

func AxisY(max, width, height, ticks float64) svg.Element {
	var (
		pos1  = svg.NewPos(0, 0)
		pos2  = svg.NewPos(0, height)
		trsl  = svg.WithTranslate(pad.Left, pad.Top)
		grp   = svg.NewGroup(svg.WithID("y-axis"), trsl)
		axis  = svg.NewLine(pos1, pos2, axisFill.Option())
		step  = height / max
		coeff = max / ticks
	)
	grp.Append(axis.AsElement())
	for i := ticks; i >= 0; i-- {
		g := svg.NewGroup()

		v := coeff * float64(i)
		x := -35.0
		y := height - (step * v)
		pos := svg.NewPos(x, y+4)
		anc := svg.WithAnchor("start")
		t := svg.NewText(fmt.Sprintf("%.0f", v), anc, pos.Option())
		g.Append(t.AsElement())

		tick := svg.NewLine(svg.NewPos(0, y), svg.NewPos(-8, y), axisFill.Option())
		g.Append(tick.AsElement())

		if i > 0 {
			stroke := axisFill
			stroke.Dash.Array = []int{5}
			line := svg.NewLine(svg.NewPos(0, y), svg.NewPos(width, y), stroke.Option())
			g.Append(line.AsElement())
		}
		grp.Append(g.AsElement())
	}
	return grp.AsElement()
}

var fills = map[string]svg.Fill{
	"failure": svg.NewFill("#7fc97f"),
	"success": svg.NewFill("#beaed4"),
	"missed":  svg.NewFill("#fdc086"),
}

func getSeries() []chart.Serie {
	// var (
	//   vs = []string{"audio", "video", "data", "other", "intercom"}
	//   cs []CategorySerie
	// )
	// for _, v := range cs {
	//   c := NewCategorySerie(v)
	//   for i := 0; i < 4; i++ {
	//
	//   }
	// }
	r1 := chart.NewSerie("audio")
	r1.Add("failure", randValue())
	r1.Add("success", randValue())
	r1.Add("missed", randValue())
	r2 := chart.NewSerie("video")
	r2.Add("failure", randValue())
	r2.Add("success", randValue())
	r3 := chart.NewSerie("data")
	r3.Add("failure", randValue())
	r3.Add("success", randValue())
	r3.Add("missed", randValue())
	r4 := chart.NewSerie("other")
	r4.Add("failure", randValue())
	r4.Add("success", randValue())
	return []chart.Serie{r1, r2, r3, r4}
}

func randValue() float64 {
	i := rand.Intn(100)
	return float64(i)
}
