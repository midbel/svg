package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/midbel/svg/chart"
)

var pad = chart.Padding{
	Top:    20,
	Bottom: 20,
	Left:   60,
	Right:  20,
}

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	var sc chart.StackedChart
	sc.Width = 640
	sc.Height = 480
	sc.Padding = pad //chart.CreatePadding(40, 20)
	sc.Ticks = 10
	sc.BarWidth = 22
	var (
		count = flag.Int("c", 5, "count")
	)
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

func getSeries() []chart.Serie {
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