package layout

import (
	"bufio"
	"io"

	"github.com/midbel/svg"
)

type Position int

type Border struct {
	Width  float64
	Height float64

	Central Renderer
	North   []Renderer
	South   []Renderer
	East    []Renderer
	West    []Renderer
}

func (b Border) Element() svg.Element {
	var (
		grid svg.SVG
	)
	grid.Dim = svg.NewDim(b.Width, b.Height)
	return grid.AsElement()	
}

func (b Border) Render(w io.Writer) error {
	ws := bufio.NewWriter(w)
	defer ws.Flush()

	b.Element().Render(ws)
	return nil
}
