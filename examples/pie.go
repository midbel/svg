package main

import (
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

func main() {
	var (
		c chart.PieChart
		s chart.Serie
	)
	var w io.Writer = os.Stdout
	c.Padding = chart.CreatePadding(20, 20)
	c.GetColor = func(_ string, i int) svg.Fill {
		return svg.NewFill(colors.YlGnBu3[i%len(colors.YlGnBu3)])
	}
	c.Width = 720
	c.Height = 720
	c.OuterRadius = 300
	c.InnerRadius = 120

	repos := []string{
		"toml",
		"json",
		"xml",
		"pcap",
		"pdf",
		"try",
		"tail",
		"linewriter",
		"wip",
		"cbor",
		"ldap",
		"transmit",
		"achile",
		"prospect",
		"svg",
		"hadock",
		"assist",
		"fig",
		"dissect",
		"pl",
		"comma",
		"jwt",
		"packit",
		"tape",
		"upifinder",
		"alea",
		"fetch",
		"uuid",
		"cli",
		"xxh",
		"ipaddr",
		"sdp",
		"hexdump",
		"curly",
	}
	for _, str := range repos {
		s.Add(str, float64(rand.Intn(10000)))
	}
	c.Render(w, s)
}
