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
	draw := svg.NewSVG(svg.WithDimension(650, 340))

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

	grp5 := svg.NewGroup(svg.WithID("line-group"))

	line1 := svg.NewLine(svg.NewPos(125, 41), svg.NewPos(175, 41), svg.WithID("line1"))
	grp5.Append(line1.AsElement())
	line2 := svg.NewLine(svg.NewPos(125, 161), svg.NewPos(175, 161), svg.WithID("line2"))
	grp5.Append(line2.AsElement())
	line5 := svg.NewLine(svg.NewPos(285, 173), svg.NewPos(335, 173), svg.WithID("line5"))
	grp5.Append(line5.AsElement())
	line6 := svg.NewLine(svg.NewPos(285, 283), svg.NewPos(335, 283), svg.WithID("line6"))
	grp5.Append(line6.AsElement())
	line7 := svg.NewLine(svg.NewPos(285, 41), svg.NewPos(495, 41), svg.WithID("line7"))
	grp5.Append(line7.AsElement())

	var path4 svg.Path
	path4.Stroke = svg.DefaultStroke
	path4.AbsMoveTo(svg.NewPos(175, 173))
	path4.RelHorizontalLine(-15)
	path4.RelVerticalLine(98)
	path4.RelHorizontalLine(15)
	grp5.Append(path4.AsElement())

	var path5 svg.Path
	path5.Stroke = svg.DefaultStroke
	path5.AbsMoveTo(svg.NewPos(335, 185))
	path5.RelHorizontalLine(-15)
	path5.RelVerticalLine(87)
	path5.RelHorizontalLine(15)
	grp5.Append(path5.AsElement())

	var path1 svg.Path
	path1.Id = "FH"
	path1.Stroke = svg.DefaultStroke
	path1.AbsMoveTo(svg.NewPos(445, 162))
	path1.RelHorizontalLine(175)
	path1.RelVerticalLine(-109)
	path1.RelHorizontalLine(-15)

	tp1 := svg.NewTextPath("F -> H", "FH")
	tp1.Offset = 40

	gt1 := svg.NewGroup(svg.WithList(tp1.AsElement(), path1.AsElement()))
	grp5.Append(gt1.AsElement())

	var path2 svg.Path
	path2.Id = "GH"
	path2.Stroke = svg.DefaultStroke
	path2.AbsMoveTo(svg.NewPos(445, 271))
	path2.RelHorizontalLine(185)
	path2.RelVerticalLine(-230)
	path2.RelHorizontalLine(-25)

	tp2 := svg.NewTextPath("G -> H", "GH")
	tp2.Offset = 40

	gt2 := svg.NewGroup(svg.WithList(tp2.AsElement(), path2.AsElement()))
	grp5.Append(gt2.AsElement())

	var path3 svg.Path
	path3.Stroke = svg.DefaultStroke
	path3.AbsMoveTo(svg.NewPos(285, 162))
	path3.RelHorizontalLine(25)
	path3.RelVerticalLine(-97)
	path3.RelHorizontalLine(185)
	grp5.Append(path3.AsElement())

	draw.Append(grp1.AsElement())
	draw.Append(grp2.AsElement())
	draw.Append(grp3.AsElement())
	draw.Append(grp4.AsElement())
	draw.Append(grp5.AsElement())

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
