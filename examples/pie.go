package main

import (
	"os"
  "io"

	"github.com/midbel/svg/chart"
)

func main() {
	var (
		c chart.PieChart
		s chart.Serie
	)
  var w io.Writer = os.Stdout
  // w = ioutil.Discard
  c.Padding = chart.CreatePadding(20, 20)
  c.Width = 400
  c.Height = 400
	c.MaxRadius = 150
	c.MinRadius = 50
	s.Add("toml", 150)
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
	c.Render(w, s)
}
