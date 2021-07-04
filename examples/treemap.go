package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"

	"github.com/midbel/svg"
)

type Appender interface {
	Append(svg.Element)
}

var colors = []string{
	"coral",
	"magenta",
	"darkturquoise",
	"yellowgreen",
	"rosybrown",
	"lightgreen",
	"lightskyblue",
	"darksalmon",
	"darkseagreen",
	"slategrey",
	"sienna",
	"sandybrown",
}

type Node struct {
	Label string
	Root  bool
	Nodes []Node
}

func (n Node) IsLeaf() bool {
	return n.Count() == 0
}

func (n Node) Count() int {
	return len(n.Nodes)
}

func (n Node) Depth() int {
	if len(n.Nodes) == 0 {
		return 1
	}
	var d int
	for _, n := range n.Nodes {
		t := n.Depth()
		if d == 0 || t > d {
			d = t
		}
	}
	return d + 1
}

func (n Node) Leaf() int {
	if n.IsLeaf() {
		return 1
	}
	var c int
	for _, n := range n.Nodes {
		c += n.Leaf()
	}
	return c
}

const defaultWidth = 900
const defaultHeight = 600

func main() {
	root := getRoot()
	height := defaultHeight / root.Leaf()

	d := svg.NewDim(defaultWidth, defaultHeight)
	canvas := svg.NewSVG(svg.WithDim(d))

	var offset int
	for i, n := range root.Nodes {
		width := defaultWidth / n.Depth()
		draw(&canvas, n, width, height, 0, i, offset, defaultWidth)
		offset += height * n.Leaf()
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	canvas.Render(w)
}

func draw(canvas Appender, root Node, width, height, levelx, levely, offsetY, available int) {
	offsetX := width * levelx

	nheight := height
	if !root.IsLeaf() {
		nheight *= root.Leaf()
	}

	nwidth := width
	if available > 0 && root.IsLeaf() {
		nwidth = available
	}

	g := svg.NewGroup(svg.WithTranslate(float64(offsetX), float64(offsetY)), svg.WithID(root.Label))
	d := svg.NewDim(float64(nwidth), float64(nheight))
	f := svg.NewFill(colors[rand.Intn(len(colors))])
	r := svg.NewRect(svg.WithDim(d), svg.WithFill(f))
	t := svg.NewText(root.Label, svg.WithPosition(10, 15))
	r.Title = fmt.Sprintf("offsetX: %d, offsetY: %d, width: %d, height: %d", offsetX, offsetY, width, nheight)
	c := svg.NewGroup(svg.WithClass("shape"))
	c.Append(r.AsElement())
	c.Append(t.AsElement())
	g.Append(c.AsElement())
	canvas.Append(g.AsElement())

	if root.Depth() >= 1 {
		levelx = 0
	}
	offsetY = 0
	for i, n := range root.Nodes {
		draw(&g, n, width, height, levelx+1, i, offsetY, available-nwidth)
		offsetY += height * n.Leaf()
	}
}

func getRoot() Node {
	nodes := []Node{
		{
			Label: "A",
			Nodes: []Node{
				{
					Label: "B",
					Nodes: []Node{
						{Label: "X"},
						{Label: "Y"},
						{Label: "Z"},
					},
				},
				{Label: "C"},
			},
		},
		{
			Label: "D",
			Nodes: []Node{
				{Label: "E"},
				{
					Label: "F",
					Nodes: []Node{
						{
							Label: "X",
						},
						{
							Label: "Y",
							Nodes: []Node{
								{Label: "M"},
								{Label: "N"},
								{Label: "O"},
								{Label: "P"},
							},
						},
						{
							Label: "Z",
							Nodes: []Node{
								{Label: "U"},
								{Label: "V"},
							},
						},
					},
				},
				{Label: "G"},
			},
		},
		{
			Label: "H",
			Nodes: []Node{
				{Label: "A"},
				{
					Label: "B",
					Nodes: []Node{
						{Label: "C"},
						{Label: "D"},
					},
				},
			},
		},
	}
	return Node{Label: "root", Root: true, Nodes: nodes}
}
