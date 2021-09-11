package svg

type Stop struct {
	Class   []string
	Color   string
	Opacity float64
	Offset  float64
}

func (s *Stop) Render(w Writer) {

}

func (s *Stop) AsElement() Element {
	return s
}

type Linear struct {
	node
	List

	Pos1   Pos
	Pos2   Pos
	Spread string
}

func (i *Linear) Render(w Writer) {
	c.render(w, "linearGradient", c.List)
}

func (i *Linear) AsElement() Element {
	return i
}

type Radial struct {
	node
	List

	Pos
	Fx     float64
	Fy     float64
	Fr     float64
	Radius float64
	Spread string
}

func (r *Radial) Render(w Writer) {
	c.render(w, "radialGradient", c.List)
}

func (r *Radial) AsElement() Element {
	return r
}

type Pattern struct {
	node
}
