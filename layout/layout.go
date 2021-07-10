package layout

import (
  "github.com/midbel/svg"
)

type Drawer interface {
  Draw() svg.Element
}

type Node struct {
  Drawer
  Nodes []Node
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

func (n Node) Count() int {
  return len(n.Nodes)
}

func (n Node) IsLeaf() bool {
	return n.Count() == 0
}

func Basic(nodes []Node, dim, siz svg.Dim, options ...svg.Option) svg.Element {
  return nil
}

const gridSplit = 2

func Grid(root Node, split int, dim svg.Dim, options ...svg.Option) svg.Element {
  return nil
}

func Tree(root Node, dim svg.Dim, options ...svg.Option) svg.Element {
  return nil
}

type appender interface {
	Append(svg.Element)
}
