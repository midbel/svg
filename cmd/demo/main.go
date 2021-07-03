package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"

	"github.com/midbel/svg"
)

const (
	rectWidth  = 120.0
	rectHeight = 90.0
)

func main() {
	draw := svg.NewSVG(svg.WithDimension(620, 340))

	grp1 := svg.NewGroup(svg.WithTranslate(10, 0), svg.WithID("level-1"))
	grp1.Append(makeRectangle("A", "192.168.67.1", 0, 10, true))
	grp1.Append(makeRectangle("B", "10.0.0.1", 0, 130, true))

	grp2 := svg.NewGroup(svg.WithTranslate(170, 0), svg.WithID("level-2"))
	grp2.Append(makeRectangle("C", "193.144.97.81", 0, 10, true))
	grp2.Append(makeRectangle("D", "8.8.8.8", 0, 130, false))
	grp2.Append(makeRectangle("E", "172.16.10.11", 0, 240, true))

	grp3 := svg.NewGroup(svg.WithTranslate(330, 0), svg.WithID("level-3"))
	grp3.Append(makeRectangle("F", "172.17.10.14", 0, 130, true))
	grp3.Append(makeRectangle("G", "172.18.10.14", 0, 240, false))

	grp4 := svg.NewGroup(svg.WithTranslate(490, 0), svg.WithID("level-3"))
	grp4.Append(makeRectangle("H", "192.168.127.1", 0, 10, true))

	line1 := svg.NewLine(svg.NewPos(130, 55), svg.NewPos(170, 55))
	draw.Append(line1.AsElement())
	line2 := svg.NewLine(svg.NewPos(130, 175), svg.NewPos(170, 175))
	draw.Append(line2.AsElement())
	line3 := svg.NewLine(svg.NewPos(230, 220), svg.NewPos(230, 240))
	draw.Append(line3.AsElement())
	line4 := svg.NewLine(svg.NewPos(390, 220), svg.NewPos(390, 240))
	draw.Append(line4.AsElement())
	line5 := svg.NewLine(svg.NewPos(290, 175), svg.NewPos(330, 175))
	draw.Append(line5.AsElement())
	line6 := svg.NewLine(svg.NewPos(290, 290), svg.NewPos(330, 290))
	draw.Append(line6.AsElement())
	line7 := svg.NewLine(svg.NewPos(290, 50), svg.NewPos(490, 50), svg.WithID("CH"))
	draw.Append(line7.AsElement())

	tp1 := svg.NewTextPath("F -> H", "FH")
	tp1.Offset = 40
	draw.Append(tp1.AsElement())

	tp2 := svg.NewTextPath("G -> H", "GH")
	tp2.Offset = 40
	draw.Append(tp2.AsElement())

	var path1 svg.Path
	path1.Id = "FH"
	path1.Stroke = svg.DefaultStroke
	path1.AbsMoveTo(svg.NewPos(450, 180))
	path1.RelHorizontalLine(100)
	path1.RelVerticalLine(-80)
	draw.Append(path1.AsElement())

	var path2 svg.Path
	path2.Id = "GH"
	path2.Stroke = svg.DefaultStroke
	path2.AbsMoveTo(svg.NewPos(450, 290))
	path2.RelHorizontalLine(110)
	path2.RelVerticalLine(-190)
	draw.Append(path2.AsElement())

	var path3 svg.Path
	path3.Stroke = svg.DefaultStroke
	path3.AbsMoveTo(svg.NewPos(290, 160))
	path3.RelHorizontalLine(20)
	path3.RelVerticalLine(-100)
	path3.RelHorizontalLine(180)
	draw.Append(path3.AsElement())

	draw.Append(grp1.AsElement())
	draw.Append(grp2.AsElement())
	draw.Append(grp3.AsElement())
	draw.Append(grp4.AsElement())

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	draw.Render(w)
}

func makeRectangle(label, ip string, x, y float64, up bool) svg.Element {
	var (
		height = rectHeight / 2
		tlab   = svg.NewText(label, svg.WithPosition(rectWidth/2, height-10), svg.WithAnchor("middle"))
		ilab   = svg.NewText(ip, svg.WithPosition(rectWidth/2, height+10), svg.WithAnchor("middle"))
	)

	g := svg.NewGroup(svg.WithID(label), svg.WithTranslate(x, y))
	g.Append(tlab.AsElement())
	g.Append(ilab.AsElement())

	options := []svg.Option{
		svg.WithStroke(svg.NewStroke("black", 1)),
		svg.WithDimension(rectWidth, rectHeight),
		svg.WithFill(svg.NewFill("none")),
		svg.WithClass("device"),
	}
	r := svg.NewRect(options...)
	r.Title = fmt.Sprintf("%s - %s (x: %.0f, y: %.0f)", label, ip, x, y)
	g.Append(r.AsElement())

	fill := "green"
	if !up {
		fill = "red"
	}
	options = []svg.Option{
		svg.WithPosition(rectWidth-6, 10),
		svg.WithRadius(4),
		svg.WithFill(svg.NewFill(fill)),
		svg.WithClass("status"),
	}
	c := svg.NewCircle(options...)
	g.Append(c.AsElement())

	var (
		right = svg.NewGroup(svg.WithClass("connector"), svg.WithTranslate(rectWidth-12, 25))
		left  = svg.NewGroup(svg.WithClass("connector"), svg.WithTranslate(0, 25))
	)
	g.Append(makeConnector(right, rand.Intn(4)))
	g.Append(makeConnector(left, rand.Intn(4)))
	return g.AsElement()
}

func makeConnector(grp svg.Group, c int) svg.Element {
	if c == 0 {
		c++
	}
	c += 1
	height := float64(c) * 12
	options := []svg.Option{
		svg.WithDimension(12, height),
	}
	r := svg.NewRect(options...)
	r.RX = 5
	r.RY = 5
	grp.Append(r.AsElement())

	for i := 0; i < c; i++ {
		options = []svg.Option{
			svg.WithRadius(4),
			svg.WithPosition(6, 6+float64(i)*12),
			svg.WithFill(svg.NewFill("black")),
		}
		c := svg.NewCircle(options...)
		grp.Append(c.AsElement())
	}
	return grp.AsElement()
}
