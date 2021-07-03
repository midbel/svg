package svg

type Option func(Element) error

func WithID(id string) Option {
	return func(e Element) error {
		e.setId(id)
		return nil
	}
}

func WithClass(class ...string) Option {
	return func(e Element) error {
		e.setClass(class)
		return nil
	}
}

func WithStyle(prop string, values ...string) Option {
	return func(e Element) error {
		e.setStyle(prop, values)
		return nil
	}
}

func WithList(es ...Element) Option {
	return func(e Element) error {
		list := NewList(es...)
		switch e := e.(type) {
		case *SVG:
			e.List = list
		case *Defs:
			e.List = list
		case *Group:
			e.List = list
		case *Rect:
			e.List = list
		case *Ellipse:
			e.List = list
		case *Circle:
			e.List = list
		case *Polygon:
			e.List = list
		case *ClipPath:
			e.List = list
		case *Mask:
			e.List = list
		case *Text:
			e.List = list
		default:
		}
		return nil
	}
}

func WithTranslate(x, y float64) Option {
	type translater interface {
		Translate(float64, float64)
	}
	return func(e Element) error {
		if e, ok := e.(translater); ok {
			e.Translate(x, y)
		}
		return nil
	}
}

func WithRotate(a, x, y float64) Option {
	type rotater interface {
		Rotate(float64, float64)
	}
	return func(e Element) error {
		if e, ok := e.(rotater); ok {
			e.Rotate(x, y)
		}
		return nil
	}
}

func WithScale(x, y float64) Option {
	type scaler interface {
		Scale(float64, float64)
	}
	return func(e Element) error {
		if e, ok := e.(scaler); ok {
			e.Scale(x, y)
		}
		return nil
	}
}

func WithSkewX(x float64) Option {
	type skewer interface {
		SkewX(float64)
	}
	return func(e Element) error {
		if e, ok := e.(skewer); ok {
			e.SkewX(x)
		}
		return nil
	}
}

func WithSkewY(y float64) Option {
	type skewer interface {
		SkewY(float64)
	}
	return func(e Element) error {
		if e, ok := e.(skewer); ok {
			e.SkewY(y)
		}
		return nil
	}
}

func WithAnchor(anchor string) Option {
	return func(e Element) error {
		if e, ok := e.(*Text); ok {
			e.Anchor = anchor
		}
		return nil
	}
}

func WithFill(fill Fill) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Line:
			e.Fill = fill
		case *PolyLine:
			e.Fill = fill
		case *Path:
			e.Fill = fill
		case *Rect:
			e.Fill = fill
		case *Circle:
			e.Fill = fill
		case *Text:
			e.Fill = fill
		case *Group:
			e.Fill = fill
		case *Use:
			e.Fill = fill
		case *SVG:
			e.Fill = fill
		case *Polygon:
			e.Fill = fill
		case *Ellipse:
			e.Fill = fill
		case *Mask:
			e.Fill = fill
		case *ClipPath:
			e.Fill = fill
		case *TextPath:
			e.Fill = fill
		default:
		}
		return nil
	}
}

func WithFont(f Font) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Text:
			e.Font = f
		default:
		}
		return nil
	}
}

func WithDim(d Dim) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Rect:
			e.Dim = d
		case *SVG:
			e.Dim = d
		case *Use:
			e.Dim = d
		case *Image:
			e.Dim = d
		case *Mask:
			e.Dim = d
		default:
		}
		return nil
	}
}

func WithDimension(x, y float64) Option {
	return WithDim(NewDim(x, y))
}

func WithRadius(r float64) Option {
	return func(e Element) error {
		if e, ok := e.(*Circle); ok {
			e.Radius = r
		}
		return nil
	}
}

func WithPos(p Pos) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Rect:
			e.Pos = p
		case *Circle:
			e.Pos = p
		case *Text:
			e.Pos = p
		case *Use:
			e.Pos = p
		case *SVG:
			e.Pos = p
		case *Ellipse:
			e.Pos = p
		case *Image:
			e.Pos = p
		case *Mask:
			e.Pos = p
		}
		return nil
	}
}

func WithPosition(x, y float64) Option {
	return WithPos(NewPos(x, y))
}

func WithStroke(s Stroke) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Rect:
			e.Stroke = s
		case *Group:
			e.Stroke = s
		case *Text:
			e.Stroke = s
		case *Use:
			e.Stroke = s
		case *SVG:
			e.Stroke = s
		case *Polygon:
			e.Stroke = s
		case *Ellipse:
			e.Stroke = s
		case *Circle:
			e.Stroke = s
		case *Line:
			e.Stroke = s
		case *PolyLine:
			e.Stroke = s
		case *Path:
			e.Stroke = s
		case *Mask:
			e.Stroke = s
		case *ClipPath:
			e.Stroke = s
		case *TextPath:
			e.Stroke = s
		default:
		}
		return nil
	}
}
