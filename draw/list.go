package draw

import (
	"fmt"
	"sort"
	"sync"
)

type Item struct {
	Drawer
	Label string
}

func MakeItem(label string) *Item {
	return &Item{Label: strings.ToUpper(label)}
}

type VisitFunc func(*Item) bool

type Set []*Item

func (s Set) String() string {
	var str []string
	for i := range s {
		str = append(str, s[i].String())
	}
	return strings.Join(str, " -> ")
}

func MakeItem(label string) *Item {
	return &Item{Label: strings.ToUpper(label)}
}

type Point struct {
	X, Y  int
	Label string
}

func makePoint(label string, x, y int) Point {
	return Point{
		X:     x,
		Y:     y,
		Label: label,
	}
}

type List struct {
	mu    sync.RWMutex
	nodes []*Item
	edges map[int][]cost
}

func NewList() *List {
	i := List{
		edges: make(map[int][]cost),
	}
	return &i
}

func (i *List) Split() []*List {
	return nil
}

func (i *List) EdgesCount() int {
	return len(i.nodes)
}

func (i *List) LinksCount() int {
	var c int
	for _, cs := range i.edges {
		c += len(cs)
	}
	return c
}

func (i *List) Render() []Point {
	return nil
}

func (i *List) CountFrom(it *Item) int {
	x, ok := i.indexOf(it)
	if !ok {
		return 0
	}
	return len(i.edges[x])
}

func (i *List) CountTo(it *Item) int {
	x, ok := i.indexOf(it)
	if !ok {
		return 0
	}
	var c int
	for _, vs := range i.edges {
		j := sort.Search(len(vs), func(j int) bool {
			return x <= vs[j].index
		})
		if j < len(vs) && vs[j].index == x {
			c++
		}
	}
	return c
}

func (i *List) IsAdjacent(fst, snd *Item) bool {
	x, ok := i.indexOf(fst)
	if !ok {
		return false
	}
	y, ok := i.indexOf(snd)
	if !ok {
		return false
	}
	j := sort.Search(len(i.edges[x]), func(j int) bool {
		return y <= i.edges[x][j].index
	})
	return j < len(i.edges[x]) && i.edges[x][j].index == y
}

func (i *List) ShortestPath(from, to *Item) []*Item {
	if i.IsAdjacent(from, to) {
		return []*Item{from, to}
	}
	return []*Item{from, to}
}

func (i *List) Walk(it *Item) []Set {
	i.mu.RLock()
	defer i.mu.RUnlock()

	var (
		is   [][]*Item
		sets []Set
	)
	is = i.walk(it, is)

	for _, s := range is {
		sets = append(sets, Set(s))
	}
	return sets
}

func (i *List) BFS(it *Item, visit VisitFunc) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	i.visitBFS(it, visit)
}

func (i *List) DFS(it *Item, visit VisitFunc) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	seen := make(map[int]struct{})
	i.visitDFS(it, seen, visit)
}

func (i *List) Neighbors(it *Item) []*Item {
	i.mu.RLock()
	defer i.mu.RUnlock()

	j, ok := i.indexOf(it)
	if !ok {
		return nil
	}
	var ns []*Item
	for _, x := range i.edges[j] {
		ns = append(ns, i.nodes[x.index])
	}
	return ns
}

func (i *List) AddNode(it *Item) {
	i.mu.Lock()
	defer i.mu.Unlock()
	j, ok := i.indexOf(it)
	if ok {
		return
	}
	i.nodes = append(i.nodes[:j], append([]*Item{it}, i.nodes[j:]...)...)
}

func (i *List) DelNode(it *Item) {
	i.mu.Lock()
	defer i.mu.Unlock()
	j, ok := i.indexOf(it)
	if !ok {
		return
	}
	delete(i.edges, j)
	i.nodes = append(i.nodes[:j], i.nodes[j+1:]...)
}

func (i *List) Order() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return len(i.nodes)
}

func (i *List) Exists(from, to *Item) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	var (
		x  int
		y  int
		ok bool
	)
	if x, ok = i.indexOf(from); !ok {
		return false
	}
	if y, ok = i.indexOf(to); !ok {
		return false
	}
	if len(i.edges[x]) == 0 {
		return false
	}
	j := sort.Search(len(i.edges[x]), func(j int) bool {
		return y <= i.edges[x][j].index
	})
	return j < len(i.edges[x]) && i.edges[x][j].index == y
}

func (i *List) AddLinkWithWeight(from, to *Item, weight int) {
	i.mu.Lock()
	defer i.mu.Unlock()
	var (
		x  int
		y  int
		ok bool
	)
	if x, ok = i.indexOf(from); !ok {
		i.nodes = append(i.nodes[:x], append([]*Item{from}, i.nodes[x:]...)...)
	}
	if y, ok = i.indexOf(to); !ok {
		i.nodes = append(i.nodes[:y], append([]*Item{to}, i.nodes[y:]...)...)
		i.edges[x] = append(i.edges[x], makeCost(y, weight))
		return
	}
	if len(i.edges[x]) == 0 {
		i.edges[x] = append(i.edges[x], makeCost(y, weight))
		return
	}
	j := sort.Search(len(i.edges[x]), func(j int) bool {
		return y <= i.edges[x][j].index
	})
	if j < len(i.edges[x]) && i.edges[x][j].index == y {
		return
	}
	c := makeCost(y, weight)
	i.edges[x] = append(i.edges[x][:j], append([]cost{c}, i.edges[x][j:]...)...)
}

func (i *List) AddLink(from, to *Item) {
	i.AddLinkWithWeight(from, to, 0)
}

func (i *List) AddLinks(from *Item, to ...*Item) {
	for _, t := range to {
		i.AddLink(from, t)
	}
}

func (i *List) DelLink(from, to *Item) {
	i.mu.Lock()
	defer i.mu.Unlock()
	var (
		x  int
		y  int
		ok bool
	)
	if x, ok = i.indexOf(from); !ok {
		return
	}
	if y, ok = i.indexOf(to); !ok {
		return
	}
	j := sort.Search(len(i.edges[x]), func(j int) bool {
		return y <= i.edges[x][j].index
	})
	if j < len(i.edges[x]) && i.edges[x][j].index == y {
		i.edges[x] = append(i.edges[x][:j], i.edges[x][j+1:]...)
	}
}

func (i *List) DelLinks(from *Item, to ...*Item) {
	for _, t := range to {
		i.DelLink(from, t)
	}
}

func (i *List) walk(it *Item, is [][]*Item) [][]*Item {
	j, ok := i.indexOf(it)
	if !ok {
		return nil
	}
	if len(i.edges[j]) == 0 {
		for x := range is {
			is[x] = append(is[x], it)
		}
		return is
	}
	var (
		z  = len(is) - 1
		ns [][]*Item
	)
	if z >= 0 && hasCycle(it, is[z]) {
		return is
	}
	for _, x := range i.edges[j] {
		var xs []*Item
		if z >= 0 {
			xs = make([]*Item, len(is[z]))
			copy(xs, is[z])
		}
		xs = append(xs, it)

		ns = append(ns, xs)
		ns = i.walk(i.nodes[x.index], ns[len(ns)-1:])
		is = append(is, ns...)
	}
	if z == 0 {
		return is[1:]
	}
	return is
}

func (i *List) visitDFS(it *Item, seen map[int]struct{}, visit VisitFunc) bool {
	j, ok := i.indexOf(it)
	if !ok {
		return true
	}
	if ok = visit(it); ok {
		return ok
	}
	for _, x := range i.edges[j] {
		if _, ok := seen[x.index]; ok {
			continue
		}
		seen[x.index] = struct{}{}
		if ok = visit != nil && i.visitDFS(i.nodes[x.index], seen, visit); ok {
			break
		}
	}
	return ok
}

func (i *List) visitBFS(it *Item, visit VisitFunc) {
	j, ok := i.indexOf(it)
	if !ok {
		return
	}
	var (
		seen = make(map[int]struct{})
		rest = []int{j}
	)
	for len(rest) > 0 {
		x := rest[0]
		if _, ok := seen[x]; ok {
			rest = rest[1:]
			continue
		}

		seen[x] = struct{}{}
		if visit != nil && visit(i.nodes[x]) {
			return
		}
		for _, c := range i.edges[x] {
			rest = append(rest, c.index)
		}
	}
}

func (i *List) indexOf(it *Item) (int, bool) {
	j := sort.Search(len(i.nodes), func(j int) bool {
		return it.Label <= i.nodes[j].Label
	})
	ok := j < len(i.nodes) && i.nodes[j].Label == it.Label
	return j, ok
}

func hasCycle(it *Item, is []*Item) bool {
	for i := range is {
		if is[i].Label == it.Label {
			return true
		}
	}
	return false
}

type cost struct {
	index  int
	weight int
}

func makeCost(index, weight int) cost {
	return cost{
		index:  index,
		weight: weight,
	}
}

type edge struct {
	From int
	To   int
}

func makeEdge(x, y int) edge {
	return edge{
		From: x,
		To:   y,
	}
}

type edgeset map[edge]struct{}

func (e edgest) Seen(f, t int) ok {
	g := makeEdge(f, t)
	_, ok := e[g]
	if !ok {
		e[g] = struct{}{}
	}
	return ok
}
