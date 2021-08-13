package main

import (
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
	var c chart.SunburstChart
	var w io.Writer = os.Stdout
	c.Padding = chart.CreatePadding(20, 20)
	c.Width = 720
	c.Height = 720
	c.OuterRadius = 320
	c.InnerRadius = 60

	hs := getHierarchy()
	c.Render(w, hs)
}

func getHierarchy() []chart.Hierarchy {
	return []chart.Hierarchy{
		{
			Label: "A",
			Value: float64(rand.Intn(limit)),
			Sub: []chart.Hierarchy{
				{
					Label: "D",
					Value: float64(rand.Intn(limit)),
					Sub: []chart.Hierarchy{
						{
							Label: "L",
							Value: float64(rand.Intn(limit)),
						},
						{
							Label: "M",
							Value: float64(rand.Intn(limit)),
						},
						{
							Label: "N",
							Value: float64(rand.Intn(limit)),
						},
					},
				},
				{
					Label: "E",
					Value: float64(rand.Intn(limit)),
				},
				{
					Label: "F",
					Value: float64(rand.Intn(limit)),
				},
			},
		},
		{
			Label: "B",
			Value: float64(rand.Intn(limit)),
			Sub: []chart.Hierarchy{
				{
					Label: "G",
					Value: float64(rand.Intn(limit)),
					Sub: []chart.Hierarchy{
						{
							Label: "X1",
							Value: float64(rand.Intn(limit)),
						},
						{
							Label: "Y1",
							Value: float64(rand.Intn(limit)),
						},
					},
				},
				{
					Label: "H",
					Value: float64(rand.Intn(limit)),
					Sub: []chart.Hierarchy{
						{
							Label: "X2",
							Value: float64(rand.Intn(limit)),
						},
						{
							Label: "Y2",
							Value: float64(rand.Intn(limit)),
						},
					},
				},
			},
		},
		{
			Label: "C",
			Value: float64(rand.Intn(limit)),
			Sub: []chart.Hierarchy{
				{
					Label: "I",
					Value: float64(rand.Intn(limit)),
					Sub: []chart.Hierarchy{
						{
							Label: "X1",
							Value: float64(rand.Intn(limit)),
						},
						{
							Label: "Y1",
							Value: float64(rand.Intn(limit)),
						},
					},
				},
				{
					Label: "J",
					Value: float64(rand.Intn(limit)),
				},
				{
					Label: "K",
					Value: float64(rand.Intn(limit)),
					Sub: []chart.Hierarchy{
						{
							Label: "X2",
							Value: float64(rand.Intn(limit)),
						},
						{
							Label: "Y2",
							Value: float64(rand.Intn(limit)),
						},
					},
				},
			},
		},
	}
}
