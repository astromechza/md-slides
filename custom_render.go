package main

import (
	"io"
	"strings"

	"github.com/russross/blackfriday"
)

type CustomHTMLRenderer struct {
	blackfriday.Renderer
}

/*
- Item
	- Paragraph
		- Text
unfortunatley have to take the approach of
*/

func (r *CustomHTMLRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {

	if node.Type == blackfriday.Item &&
		node.FirstChild != nil &&
		node.FirstChild.Type == blackfriday.Paragraph &&
		node.FirstChild.FirstChild != nil &&
		node.FirstChild.FirstChild.Type == blackfriday.Text {
		if strings.HasPrefix(string(node.FirstChild.FirstChild.Literal), "[ ] ") {
			node.FirstChild.FirstChild.Literal = node.FirstChild.FirstChild.Literal[4:]
			nn := blackfriday.NewNode(blackfriday.HTMLSpan)
			nn.Literal = []byte(`<input type="checkbox" disabled="">`)
			nn.Next = node.FirstChild
			node.FirstChild.Prev = nn
			node.FirstChild = nn
		} else if strings.HasPrefix(string(node.FirstChild.FirstChild.Literal), "[x] ") {
			node.FirstChild.FirstChild.Literal = node.FirstChild.FirstChild.Literal[4:]
			nn := blackfriday.NewNode(blackfriday.HTMLSpan)
			nn.Literal = []byte(`<input type="checkbox" disabled="" checked>`)
			nn.Next = node.FirstChild
			node.FirstChild.Prev = nn
			node.FirstChild = nn
		}
	}
	return r.Renderer.RenderNode(w, node, entering)
}
