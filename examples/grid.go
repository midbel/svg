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

	for i, n := range root.Nodes {
		width := defaultWidth / n.Depth()
		draw(&canvas, n, width, height)
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	canvas.Render(w)
}

func draw(canvas Appender, root Node, width, height int) {

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
