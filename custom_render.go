package main

import (
	"io"
	"strings"

	"github.com/russross/blackfriday"
)

type CustomHTMLRenderer struct {
	blackfriday.Renderer
}

func (r *CustomHTMLRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	if entering && node.Type == blackfriday.Text {
		if node.Parent != nil &&
			node.Parent.Type == blackfriday.Paragraph &&
			node.Parent.Parent != nil &&
			node.Parent.Parent.Type == blackfriday.Item {
			if strings.HasPrefix(string(node.Literal), "[ ] ") {
				w.Write([]byte(`<input type="checkbox" disabled="">`))
				node.Literal = node.Literal[4:]
			} else if strings.HasPrefix(string(node.Literal), "[x] ") || strings.HasSuffix(string(node.Literal), "[X] ") {
				w.Write([]byte(`<input type="checkbox" checked disabled="">`))
				node.Literal = node.Literal[4:]
			}
		}
	}
	return r.Renderer.RenderNode(w, node, entering)
}
