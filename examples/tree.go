package main

import (
	"bufio"
	"math/rand"
	"os"

	"github.com/midbel/svg"
	"github.com/midbel/svg/draw"
)

const delta = 3

type Shape struct {
	Label string
}

func (s Shape) Draw(ctx draw.Context) svg.Element {
	options := []svg.Option{
		ctx.Pos.Option(),
		ctx.Dim.Option(),
	}
	r := svg.NewRect(options...)

	ctx.X += ctx.W / 2
	ctx.Y += ctx.H / 2
	if ctx.Leaf {
		var offset float64
		switch mod := (ctx.Nth) % 3; mod {
		case 0:
		case 1:
			offset += ctx.W
		case 2:
			offset -= ctx.W
		}
		ctx.X += offset / delta
	}
	options = []svg.Option{
		svg.WithRadius(12),
		ctx.Pos.Option(),
		svg.DefaultStroke.Option(),
		svg.NewFill(randomColor()).Option(),
	}
	c := svg.NewCircle(options...)
	c.Title = s.Label

	g := svg.NewGroup(svg.WithID("grp-" + s.Label))
	g.Append(r.AsElement())
	g.Append(c.AsElement())
	return g.AsElement()
}

func randomColor() string {
	return svg.Colors[rand.Intn(len(svg.Colors))]
}

func main() {
	root := getRoot()
	canvas := draw.Tree(root, svg.NewDim(900, 600))

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	canvas.Render(w)
}

func getRoot() draw.Node {
	nodes := []draw.Node{
		{
			Drawer: Shape{Label: "A-1-1"},
			Nodes: []draw.Node{
				{
					Drawer: Shape{Label: "B-1-2"},
					Nodes: []draw.Node{
						{Drawer: Shape{Label: "X-1-3"}},
						{Drawer: Shape{Label: "Y-1-3"}},
						{Drawer: Shape{Label: "Z-1-3"}},
					},
				},
				{Drawer: Shape{Label: "C-1-2"}},
			},
		},
		{
			Drawer: Shape{Label: "D-2-1"},
			Nodes: []draw.Node{
				{Drawer: Shape{Label: "E-2-1"}},
				{
					Drawer: Shape{Label: "F-2-2"},
					Nodes: []draw.Node{
						{
							Drawer: Shape{Label: "X-2-3"},
						},
						{
							Drawer: Shape{Label: "Y-2-4"},
							Nodes: []draw.Node{
								{Drawer: Shape{Label: "M-2-5"}},
								{Drawer: Shape{Label: "N-2-5"}},
								{Drawer: Shape{Label: "O-2-5"}},
								{Drawer: Shape{Label: "P-2-5"}},
							},
						},
						{
							Drawer: Shape{Label: "Z-2-3"},
							Nodes: []draw.Node{
								{Drawer: Shape{Label: "U-2-4"}},
								{Drawer: Shape{Label: "V-2-4"}},
							},
						},
					},
				},
				{Drawer: Shape{Label: "G-2-2"}},
			},
		},
		{
			Drawer: Shape{Label: "H-3-1"},
			Nodes: []draw.Node{
				{Drawer: Shape{Label: "A-3-2"}},
				{
					Drawer: Shape{Label: "B-3-3"},
					Nodes: []draw.Node{
						{Drawer: Shape{Label: "C-3-4"}},
						{Drawer: Shape{Label: "D-3-4"}},
					},
				},
			},
		},
	}
	return draw.Node{Nodes: nodes}
}
