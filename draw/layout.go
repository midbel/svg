package draw

import (
	"sort"

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

func Grid(root Node, dim svg.Dim, options ...svg.Option) svg.Element {
	var (
		canvas = svg.NewSVG(append(options, dim.Option())...)
		group  svg.Group
		state  gridstate
	)
	state.Dim = dim
	state.Draw(&group, root.Nodes)
	canvas.Append(group.AsElement())
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

type gridstate struct {
	svg.Dim
	svg.Pos
	Level int
	Copy  bool
}

func (g gridstate) Draw(app appender, nodes []Node) {
	var (
		size  = len(nodes)
		step  = size / 2
		horiz = g.Level%2 == 0
	)
	if size <= 4 {
		g.draw(app, nodes)
		return
	}
	if step%2 != 0 {
		step++
	}
	if horiz {
		g.W /= 2
	} else {
		g.H /= 2
	}
	g.Level++
	for i := 0; i < size; i += step {
		j := i + step
		if j > size {
			j = size
		}
		if i > 0 {
			if horiz {
				g.X += g.W
			} else {
				g.X = g.W
				g.Y += g.H
			}
		}
		g.Draw(app, nodes[i:j])
	}
}

func (g gridstate) draw(app appender, nodes []Node) {
	draw := func(n Node, c Context) {
		app.Append(n.Draw(c))
		g := gridstate{
			Dim:  c.Dim,
			Pos:  c.Pos,
			Copy: !g.Copy,
		}
		g.Draw(app, n.Nodes)
	}
	ctx := Context{
		Dim: g.Dim,
		Pos: g.Pos,
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Depth() < nodes[j].Depth()
	})
	switch len(nodes) {
	case 0:
	case 1:
		draw(nodes[0], ctx)
		app.Append(nodes[0].Draw(ctx))
	case 2:
		if g.Copy {
			ctx.W /= 2
		} else {
			ctx.H /= 2
		}
		draw(nodes[0], ctx)
		if g.Copy {
			ctx.X += ctx.W
		} else {
			ctx.Y += ctx.H
		}
		draw(nodes[1], ctx)
	case 3:
		ctx.W /= 2
		ctx.H /= 2
		draw(nodes[0], ctx)
		ctx.X += ctx.W
		draw(nodes[1], ctx)
		ctx.Y += ctx.H
		ctx.X -= ctx.W
		ctx.W += ctx.W
		draw(nodes[2], ctx)
	case 4:
		ctx.W /= 2
		ctx.H /= 2
		draw(nodes[0], ctx)
		ctx.X += ctx.W
		draw(nodes[1], ctx)
		ctx.X -= ctx.W
		ctx.Y += ctx.H
		draw(nodes[2], ctx)
		ctx.X += ctx.W
		draw(nodes[3], ctx)
	}
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
