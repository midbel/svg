package main

import (
	"fmt"
	"io"
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
	var c chart.TreemapChart
	var w io.Writer = os.Stdout
	c.Padding = chart.CreatePadding(20, 20)
	c.Tiling = chart.TilingAlternate
	c.Width = 960
	c.Height = 480

	hs := getHierarchy(1 + rand.Intn(5))
	// hs := getHierarchy(4)
	c.Render(w, hs)
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
