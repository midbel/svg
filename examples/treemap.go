package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/midbel/svg"
	"github.com/midbel/svg/chart"
)

func init() {
	rand.Seed(time.Now().Unix())
}

const limit = 100

func main() {
	flag.Parse()
	var (
		hs = load(flag.Arg(0), 1+rand.Intn(5))
		c1 = getChart(hs, chart.TilingHorizontal)
		c2 = getChart(hs, chart.TilingVertical)
		c3 = getChart(hs, chart.TilingAlternate)
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

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	area.Render(w)
}

var letters = "ABCDEFGHIJKLMNOPQRSTUVXYZ"

func getChart(hs []chart.Hierarchy, tiling chart.TilingMethod) svg.Element {
	var c chart.TreemapChart
	c.Padding = chart.CreatePadding(20, 20)
	c.Tiling = tiling
	c.Width = 480
	c.Height = 360

	return c.RenderElement(hs)
}

func load(file string, level int) []chart.Hierarchy {
	r, err := os.Open(file)
	if err != nil {
		return getHierarchy(level)
	}
	defer r.Close()
	var h chart.Hierarchy
	if err := json.NewDecoder(r).Decode(&h); err != nil {
		return nil
	}
	return []chart.Hierarchy{h}
}

func getHierarchy(level int) []chart.Hierarchy {
	if level <= 0 {
		return nil
	}
	var xs []chart.Hierarchy
	for i := 0; i < count(); i++ {
		a := letters[rand.Intn(len(letters))]
		n := chart.Hierarchy{
			Label: fmt.Sprintf("%02x-%d-%d", a, i, level),
			Value: float64(rand.Intn(limit)),
			Sub:   getHierarchy(level - 1),
		}
		xs = append(xs, n)
	}
	return xs
}

func count() int {
	return 1 + rand.Intn(10)
}
