package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/midbel/svg/chart"
)

func init() {
	rand.Seed(time.Now().Unix())
}

const limit = 100

func main() {
	var (
		width   = flag.Float64("x", 960, "width")
		height  = flag.Float64("y", 960, "height")
		discard = flag.Bool("d", false, "discard")
	)
	flag.Parse()
	var (
		hs = load(flag.Arg(0), 1+rand.Intn(5))
		c  chart.TreemapChart
	)
	c.Padding = chart.CreatePadding(10, 10)
	c.Tiling = chart.TilingSquarify
	c.Width = *width
	c.Height = *height

	var w io.Writer = os.Stdout
	if *discard {
		w = ioutil.Discard
	}
	c.Render(w, hs)
}

var letters = "ABCDEFGHIJKLMNOPQRSTUVXYZ"

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
