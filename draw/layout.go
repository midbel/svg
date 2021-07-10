package draw

import (
	"github.com/midbel/svg"
)

type Context struct {
	svg.Pos
	svg.Dim
	Depth int
	Nth   int
	Leaf  bool
}

type Drawer interface {
	Draw(Context) svg.Element
}

func Basic(nodes []Node, dim, siz svg.Dim, options ...svg.Option) svg.Element {
	return nil
}

const gridSplit = 2

func Grid(root Node, split int, dim svg.Dim, options ...svg.Option) svg.Element {
	canvas := svg.NewSVG(append(options, dim.Option())...)
	return canvas.AsElement()
}

func Tree(root Node, dim svg.Dim, options ...svg.Option) svg.Element {
	var (
		canvas = svg.NewSVG(append(options, dim.Option())...)
		group  svg.Group
		state  treestate
	)

	state.H = dim.H / float64(root.Leaf())
	state.Width = dim.W
	for i := range root.Nodes {
		root.Nodes[i].idx = i

		state.W = dim.W / float64(root.Nodes[i].Depth())
		state.Draw(&group, root.Nodes[i])
		state.Y += state.H * float64(root.Nodes[i].Leaf())
	}
	canvas.Append(group.AsElement())
	return canvas.AsElement()
}

type appender interface {
	Append(svg.Element)
}

type treestate struct {
	svg.Dim
	svg.Pos
	Level int
	Width float64
}

func (s treestate) Draw(app appender, root Node) {
	if !root.IsLeaf() {
		s.H *= float64(root.Leaf())
	}
	if s.Width > 0 && root.IsLeaf() {
		s.W = s.Width
	}
	ctx := Context{
		Pos:   s.Pos,
		Dim:   s.Dim,
		Depth: s.Level,
		Nth:   root.idx,
		Leaf:  root.IsLeaf(),
	}
	app.Append(root.Draw(ctx))

	s.H /= float64(root.Leaf())
	s.X += s.W
	s.Width -= s.W
	s.Level++
	for i := range root.Nodes {
		root.Nodes[i].idx = i
		s.Draw(app, root.Nodes[i])
		s.Y += s.H * float64(root.Nodes[i].Leaf())
	}
}
