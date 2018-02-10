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

const uncheckedSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path d="M5 2c-1.654 0-3 1.346-3 3v14c0 1.654 1.346 3 3 3h14c1.654 0 3-1.346 3-3v-14c0-1.654-1.346-3-3-3h-14zm19 3v14c0 2.761-2.238 5-5 5h-14c-2.762 0-5-2.239-5-5v-14c0-2.761 2.238-5 5-5h14c2.762 0 5 2.239 5 5z"/></svg>`
const checkedSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path d="M10.041 17l-4.5-4.319 1.395-1.435 3.08 2.937 7.021-7.183 1.422 1.409-8.418 8.591zm-5.041-15c-1.654 0-3 1.346-3 3v14c0 1.654 1.346 3 3 3h14c1.654 0 3-1.346 3-3v-14c0-1.654-1.346-3-3-3h-14zm19 3v14c0 2.761-2.238 5-5 5h-14c-2.762 0-5-2.239-5-5v-14c0-2.761 2.238-5 5-5h14c2.762 0 5 2.239 5 5z"/></svg>`

func (r *CustomHTMLRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {

	if node.Type == blackfriday.Item &&
		node.FirstChild != nil &&
		node.FirstChild.Type == blackfriday.Paragraph &&
		node.FirstChild.FirstChild != nil &&
		node.FirstChild.FirstChild.Type == blackfriday.Text {
		if strings.HasPrefix(string(node.FirstChild.FirstChild.Literal), "[ ]") {
			node.FirstChild.FirstChild.Literal = node.FirstChild.FirstChild.Literal[3:]
			nn := blackfriday.NewNode(blackfriday.HTMLSpan)
			nn.Literal = []byte(uncheckedSVG)
			nn.Next = node.FirstChild
			node.FirstChild.Prev = nn
			node.FirstChild = nn
		} else if strings.HasPrefix(string(node.FirstChild.FirstChild.Literal), "[x]") {
			node.FirstChild.FirstChild.Literal = node.FirstChild.FirstChild.Literal[3:]
			nn := blackfriday.NewNode(blackfriday.HTMLSpan)
			nn.Literal = []byte(checkedSVG)
			nn.Next = node.FirstChild
			node.FirstChild.Prev = nn
			node.FirstChild = nn
		}
	}
	return r.Renderer.RenderNode(w, node, entering)
}
