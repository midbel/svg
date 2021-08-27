package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
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

const limit = 100

func main() {
	flag.Parse()
	var (
		c chart.SunburstChart
		w io.Writer = os.Stdout
	)
	c.GetColor = func(_ string, i int) svg.Fill {
		return svg.NewFill(colors.Paired12[i%len(colors.Paired12)])
	}
	c.Padding = chart.CreatePadding(20, 20)
	c.Width = 800
	c.Height = 800
	c.OuterRadius = 380
	c.InnerRadius = 0

	hs := load(flag.Arg(0), 1+rand.Intn(5))
	c.Render(w, hs)
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

var letters = "ABCDEFGHIJKLMNOPQRSTUVXYZ"

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
