package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/midbel/svg"
	"github.com/midbel/svg/colors"
)

func main() {
	snd := flag.Bool("s", false, "second")
	flag.Parse()

	if *snd {
		graph2()
	} else {
		graph()
	}
}

func example() {
	var (
		list  = NewList()
		zero  = MakeItem("0")
		one   = MakeItem("1")
		two   = MakeItem("2")
		three = MakeItem("3")
		four  = MakeItem("4")
		five  = MakeItem("5")
		six   = MakeItem("6")
	)
	list.AddNode(zero)
	list.AddNode(one)
	list.AddNode(two)
	list.AddNode(three)
	list.AddNode(four)
	list.AddNode(five)
	list.AddNode(six)

	list.AddLinkWithWeight(zero, one, 2)
	list.AddLinkWithWeight(zero, two, 6)
	list.AddLinkWithWeight(one, three, 5)
	list.AddLinkWithWeight(two, three, 8)
	list.AddLinkWithWeight(three, five, 15)
	list.AddLinkWithWeight(three, four, 10)
	list.AddLinkWithWeight(five, four, 6)
	list.AddLinkWithWeight(five, six, 6)
	list.AddLinkWithWeight(four, six, 2)

	canvas := list.Render(zero)
	ws := bufio.NewWriter(os.Stdout)
	defer ws.Flush()
	canvas.Render(ws)
}

func graph2() {
	var (
		list = NewList()
		a    = MakeItem("a")
		b    = MakeItem("b")
		c    = MakeItem("c")
		d    = MakeItem("d")
		e    = MakeItem("e")
		f    = MakeItem("f")
		g    = MakeItem("g")
		h    = MakeItem("h")
		i    = MakeItem("i")
	)

	list.AddNode(a)
	list.AddNode(c)
	list.AddNode(f)
	list.AddNode(b)
	list.AddNode(e)
	list.AddNode(i)
	list.AddNode(g)
	list.AddNode(h)

	list.AddLinks(a, c, d)
	list.AddLinks(b, c, d, e, f)
	list.AddLinks(c, d, e)
	list.AddLinks(d, f)
	list.AddLinks(e, f, g, h)
	list.AddLinks(f, g, i)
	list.AddLinks(h, i)

	var (
		cs = list.Render(a)
		ws = bufio.NewWriter(os.Stdout)
	)
	defer ws.Flush()
	cs.Render(ws)
}

func graph() {
	var (
		list = NewList()
		a    = MakeItem("a")
		b    = MakeItem("b")
		c    = MakeItem("c")
		d    = MakeItem("d")
		e    = MakeItem("e")
		f    = MakeItem("f")
		g    = MakeItem("g")
		h    = MakeItem("h")
		i    = MakeItem("i")
		j    = MakeItem("j")
		k    = MakeItem("k")
		l    = MakeItem("l")
		m    = MakeItem("m")
		n    = MakeItem("n")
		o    = MakeItem("o")
		u    = MakeItem("u")
		v    = MakeItem("v")
		w    = MakeItem("w")
		x    = MakeItem("x")
		y    = MakeItem("y")
		z    = MakeItem("z")
	)
	list.AddNode(a)
	list.AddNode(b)
	list.AddNode(c)
	list.AddNode(d)
	list.AddNode(e)
	list.AddNode(f)
	list.AddNode(g)
	list.AddNode(h)
	list.AddNode(i)
	list.AddNode(j)
	list.AddNode(k)
	list.AddNode(m)
	list.AddNode(n)
	list.AddNode(o)
	list.AddNode(u)
	list.AddNode(v)
	list.AddNode(w)
	list.AddNode(x)
	list.AddNode(y)
	list.AddNode(z)

	list.AddLinkWithWeight(a, b, 5)
	list.AddLinkWithWeight(a, c, 3)
	list.AddLinkWithWeight(b, c, 5)
	list.AddLinkWithWeight(b, d, 5)
	list.AddLinkWithWeight(c, d, 2)
	list.AddLinkWithWeight(c, e, 2)
	list.AddLinkWithWeight(c, f, 2)
	list.AddLinkWithWeight(d, i, 5)
	list.AddLinkWithWeight(e, g, 1)
	list.AddLinkWithWeight(f, h, 5)
	list.AddLinkWithWeight(f, l, 5)
	list.AddLinkWithWeight(g, i, 2)
	list.AddLinkWithWeight(g, j, 2)
	list.AddLinkWithWeight(g, k, 1)
	list.AddLinkWithWeight(j, n, 1)
	list.AddLinkWithWeight(i, j, 5)
	list.AddLinkWithWeight(h, k, 1)
	list.AddLinkWithWeight(k, n, 1)
	list.AddLinkWithWeight(k, o, 1)
	list.AddLinkWithWeight(n, o, 1)
	list.AddLinkWithWeight(n, m, 1)
	list.AddLinkWithWeight(o, m, 1)
	list.AddLinkWithWeight(a, u, 1)
	list.AddLinkWithWeight(u, v, 1)
	list.AddLinkWithWeight(u, x, 1)
	list.AddLinkWithWeight(u, y, 1)
	list.AddLinkWithWeight(v, w, 1)
	list.AddLinkWithWeight(x, z, 1)
	list.AddLinkWithWeight(y, z, 1)
	list.AddLinkWithWeight(w, z, 1)

	var (
		cs = list.Render(a)
		ws = bufio.NewWriter(os.Stdout)
	)
	defer ws.Flush()
	cs.Render(ws)
}

func neighbors(list *List, all []*Item) {
	for _, a := range all {
		fmt.Print(a.Label)
		fmt.Print(": ")
		for x, i := range list.Neighbors(a) {
			if x > 0 {
				fmt.Print(", ")
			}
			fmt.Print(i.Label)
		}
		fmt.Println()
	}
}

func walk(list *List, it *Item) {
	now := time.Now()
	for _, s := range list.Walk(it) {
		fmt.Println("path:", s)
	}
	fmt.Println("walking: elapsed", time.Since(now))
}

type Item struct {
	Label string
	index int
}

func (i *Item) String() string {
	return i.Label
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

type List struct {
	nodes []*Item
	edges map[int][]cost
}

func NewList() *List {
	i := List{
		edges: make(map[int][]cost),
	}
	return &i
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

const (
	defaultWidth  = 1200
	defaultHeight = 720
)

func (i *List) Render(it *Item) svg.Element {
	var (
		dim    = svg.NewDim(defaultWidth, defaultHeight)
		elm    = svg.NewSVG(dim.Option())
		groups = i.groups(it)
		areas  []Area
		width  = defaultWidth / len(groups)
		height int
	)
	const (
		fifty  = 50
		fourty = 30
	)
	for x, gs := range groups {
		if len(gs) == 0 {
			continue
		}
		height = defaultHeight / len(gs)
		var (
			offx = (width / 2) - (fifty / 2)
			offy = (height / 2) - (fourty / 2)
		)
		for y, g := range gs {
			j, _ := i.indexOf(g)
			a := Area{
				Curr:  j,
				Label: g.Label,
				Index: g.index,
				Dim:   svg.NewDim(fifty, fourty),
				Pos:   svg.NewPos(float64((x*width)+offx), float64((y*height)+offy)),
			}
			areas = append(areas, a)
		}
	}
	rel := svg.NewGroup(svg.WithID("edges"))
	for _, f := range areas {
		for _, v := range i.edges[f.Curr] {
			var (
				t = searchAreas(areas, v.index)
				x = float64(fifty / 2)
				y = float64(fourty / 2)
				s = svg.NewStroke("darkgray", 1)
				i = svg.NewLine(f.Pos.Adjust(x, y), t.Pos.Adjust(x, y), s.Option())
			)
			rel.Append(i.AsElement())
		}
	}

	nod := svg.NewGroup(svg.WithID("nodes"))
	for _, p := range areas {
		options := []svg.Option{
			p.Dim.Option(),
			p.Pos.Option(),
			svg.NewFill(randomColor()).Option(),
		}
		r := svg.NewRect(options...)
		r.Title = fmt.Sprintf("%s: %d", p.Label, p.Index)
		nod.Append(r.AsElement())
	}
	elm.Append(rel.AsElement())
	elm.Append(nod.AsElement())
	return elm.AsElement()
}

func (i *List) groups(it *Item) [][]*Item {
	var (
		seen  = make(map[int]struct{})
		list  [][]int
		group [][]*Item
	)
	j, ok := i.indexOf(it)
	if !ok {
		return nil
	}
	list = append(list, []int{j})
	for len(list) > 0 {
		var (
			curr     = list[0]
			children []int
			other    []*Item
		)
		for _, n := range curr {
			if _, ok := seen[n]; ok {
				continue
			}
			seen[n] = struct{}{}
			other = append(other, i.nodes[n])
			for _, c := range i.edges[n] {
				if _, ok := seen[c.index]; ok {
					continue
				}
				children = append(children, c.index)
			}
		}
		list = list[1:]
		if len(children) > 0 {
			list = append(list, children)
		}
		// sort.Slice(other, func(i, j int) bool {
		// 	return other[i].index < other[j].index
		// })
		group = append(group, other)
	}
	return group
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
	i.visitBFS(it, visit)
}

func (i *List) DFS(it *Item, visit VisitFunc) {
	seen := make(map[int]struct{})
	i.visitDFS(it, seen, visit)
}

func (i *List) Neighbors(it *Item) []*Item {
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
	j, ok := i.indexOf(it)
	if ok {
		return
	}
	it.index = len(i.nodes)
	i.nodes = append(i.nodes[:j], append([]*Item{it}, i.nodes[j:]...)...)
}

func (i *List) DelNode(it *Item) {
	j, ok := i.indexOf(it)
	if !ok {
		return
	}
	delete(i.edges, j)
	i.nodes = append(i.nodes[:j], i.nodes[j+1:]...)
}

func (i *List) Order() int {
	return len(i.nodes)
}

func (i *List) Exists(from, to *Item) bool {
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
		i.edges[x] = append(i.edges[x], makeCost(y, weight, 0))
		return
	}
	if len(i.edges[x]) == 0 {
		i.edges[x] = append(i.edges[x], makeCost(y, weight, 0))
		return
	}
	j := sort.Search(len(i.edges[x]), func(j int) bool {
		return y <= i.edges[x][j].index
	})
	if j < len(i.edges[x]) && i.edges[x][j].index == y {
		return
	}
	c := makeCost(y, weight, len(i.edges[x]))
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
	rank   int
}

func makeCost(index, weight, rank int) cost {
	return cost{
		index:  index,
		weight: weight,
		rank:   rank,
	}
}

func randomColor() string {
	return colors.Paired10[rand.Intn(len(colors.Paired10))]
}

type Area struct {
	Label string
	Curr  int
	Index int
	svg.Pos
	svg.Dim
}

func searchAreas(areas []Area, n int) Area {
	var a Area
	for _, i := range areas {
		if i.Curr == n {
			a = i
			break
		}
	}
	return a
}
