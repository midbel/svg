package chart

type Hierarchy struct {
	Label string
	Value float64
	Sub   []Hierarchy
}

func (h Hierarchy) Sum() float64 {
	if h.isLeaf() {
		return h.Value
	}
	var s float64
	for i := range h.Sub {
		s += h.Sub[i].Value
	}
	return s
}

func (h Hierarchy) Depth() int {
	if h.isLeaf() {
		return 1
	}
	var d int
	for i := range h.Sub {
		x := h.Sub[i].Depth()
		if x > d {
			d = x
		}
	}
	return d + 1
}

func (h Hierarchy) Len() int {
	return len(h.Sub)
}

func (h Hierarchy) isLeaf() bool {
	return h.Len() == 0
}

type HeatmapChart struct {
	Chart
}

type TreemapChart struct {
	Chart
}

func getDepth(series []Hierarchy) float64 {
	var d int
	for i := range series {
		x := series[i].Depth()
		if x > d {
			d = x
		}
	}
	return float64(d)
}

func getSum(series []Hierarchy) float64 {
	var sum float64
	for i := range series {
		sum += series[i].Value
	}
	if sum == 0 {
		return 1
	}
	return sum
}
