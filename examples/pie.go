package main

import (
	"os"

	"github.com/midbel/svg/chart"
)

func main() {
	var (
		c chart.PieChart
		s chart.Serie
	)
	c.Radius = 300
	s.Add("toml", 100)
	s.Add("json", 78)
	s.Add("xml", 12)
	s.Add("ber", 70)
	s.Add("cbor", 80)
	s.Add("ldap", 98)
	s.Add("transmit", 200)
	s.Add("achile", 100)
	s.Add("prospect", 50)
	s.Add("svg", 178)
	s.Add("hadock", 500)
	s.Add("assist", 280)
	s.Add("inspect", 270)
	c.Render(os.Stdout, s)
}
