package svg

import (
	"fmt"
	"io"
)

const (
	namespace = `http://www.w3.org/2000/svg`
	prolog    = `<?xml version="1.0" encoding="utf-8"?>`
)

type Writer interface {
	io.ByteWriter
	io.StringWriter
}

type Element interface {
	Render(Writer)
}

type node struct {
	Title string
	Desc  string
	Data  []Datum

	Display    string
	Visibility string

	Id     string
	Class  []string
	Styles map[string][]string

	Clip      string
	Rendering string
}

func (n *node) Attributes() []string {
	var attrs []string
	if n.Id != "" {
		attrs = append(attrs, appendString("id", n.Id))
	}
	if len(n.Class) > 0 {
		attrs = append(attrs, appendStringArray("class", n.Class, space))
	}
	if n.Clip != "" {
		url := fmt.Sprintf("url(#%s)", n.Clip)
		attrs = append(attrs, appendString("clip-path", url))
	}
	if n.Rendering != "" {
		attrs = append(attrs, appendString("shape-rendering", n.Rendering))
	}
	if len(n.Data) > 0 {
		for i := range n.Data {
			attrs = append(attrs, n.Data[i].Attributes()...)
		}
	}
	return attrs
}

func (n *node) render(w Writer, name string, list List, attrs ...Attribute) {
	var as []string
	attrs = append(attrs, n)
	for _, a := range attrs {
		as = append(as, a.Attributes()...)
	}
	writeElement(w, name, as, func() {
		writeTitle(w, n.Title)
		writeDesc(w, n.Desc)
		list.Render(w)
	})
}

func writeElement(w Writer, name string, attrs []string, inner func()) {
	closed := inner == nil
	writeOpenElement(w, name, closed, attrs)
	if !closed {
		inner()
		writeCloseElement(w, name)
	}
}

func writeTitle(w Writer, str string) {
	if str == "" {
		return
	}
	writeString(w, "title", str)
}

func writeDesc(w Writer, str string) {
	if str == "" {
		return
	}
	writeString(w, "desc", str)
}

func writeString(w Writer, name, str string) {
	writeOpenElement(w, name, false, nil)
	w.WriteString(str)
	writeCloseElement(w, name)
}

func writeOpenElement(w Writer, name string, closed bool, attrs []string) {
	w.WriteByte(langle)
	w.WriteString(name)
	for i := range attrs {
		w.WriteByte(space)
		w.WriteString(attrs[i])
	}
	if closed {
		w.WriteByte(space)
		w.WriteByte(slash)
	}
	w.WriteByte(rangle)
}

func writeCloseElement(w Writer, name string) {
	w.WriteByte(langle)
	w.WriteByte(slash)
	w.WriteString(name)
	w.WriteByte(rangle)
}
