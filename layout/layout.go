package layout

import (
	"github.com/midbel/svg"
)

type Renderer interface {
	Element() svg.Element
}
