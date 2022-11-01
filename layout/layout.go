package layout

import (
	"github.com/midbel/svg"
)

type Renderer interface {
	Element() svg.Element
}

type Padding struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}
