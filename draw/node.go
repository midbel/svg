package draw

type Node struct {
	Drawer
	Nodes []Node

	idx int
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
