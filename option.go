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

func WithTranslate(x, y float64) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Rect:
			e.Transform.TX, e.Transform.TY = x, y
		case *Text:
			e.Transform.TX, e.Transform.TY = x, y
		case *Group:
			e.Transform.TX, e.Transform.TY = x, y
		}
		return nil
	}
}

func WithRotate(a, x, y float64) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Rect:
			e.Transform.RA = a
			e.Transform.RX, e.Transform.RY = x, y
		case *Text:
			e.Transform.RA = a
			e.Transform.RX, e.Transform.RY = x, y
		case *Group:
			e.Transform.RA = a
			e.Transform.RX, e.Transform.RY = x, y
		}
		return nil
	}
}

func WithScale(x, y float64) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Rect:
			e.Transform.SX, e.Transform.SY = x, y
		case *Text:
			e.Transform.SX, e.Transform.SY = x, y
		case *Group:
			e.Transform.SX, e.Transform.SY = x, y
		}
		return nil
	}
}

func WithSkewX(x float64) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Rect:
			e.Transform.KX = x
		case *Text:
			e.Transform.KX = x
		case *Group:
			e.Transform.KX = x
		}
		return nil
	}
}

func WithSkewY(y float64) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Rect:
			e.Transform.KY = y
		case *Text:
			e.Transform.KY = y
		case *Group:
			e.Transform.KY = y
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

func WithFill(fill string) Option {
	return func(e Element) error {
		switch e := e.(type) {
		case *Line:
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
		case *SVG:
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
		default:
		}
		return nil
	}
}
