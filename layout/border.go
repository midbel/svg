package layout

import (
	"bufio"
	"io"
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

func (b Border) Render(w io.Writer) error {
	var (
		grid svg.SVG
	)
	grid.Dim = svg.NewDim(g.Width, g.Height)

	ws := bufio.NewWriter(w)
	defer ws.Flush()

	grid.Render(ws)
	return nil
}
