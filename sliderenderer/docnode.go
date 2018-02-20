package sliderenderer

import (
	"bytes"
	"net/url"
	"strings"

	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type DocumentNode struct {
	blackfriday.Node
	Settings url.Values
}

func (n *DocumentNode) SlideClasses() string {
	bodyClasses := []string{"slide-wrap"}
	if n.Settings.Get("halign") != "" {
		bodyClasses = append(bodyClasses, "slide-wrap-halign-"+n.Settings.Get("halign"))
	}
	if n.Settings.Get("valign") != "" {
		bodyClasses = append(bodyClasses, "slide-wrap-valign-"+n.Settings.Get("valign"))
	}
	if n.Settings.Get("talign") != "" {
		bodyClasses = append(bodyClasses, "slide-wrap-talign-"+n.Settings.Get("talign"))
	}
	return strings.Join(bodyClasses, " ")
}

func (n *DocumentNode) HasFooter() bool {
	_, ok := n.Settings["footer"]
	return ok
}

func (n *DocumentNode) Footer() string {
	return n.Settings.Get("footer")
}

func breakIntoDocumentNodes(node *blackfriday.Node) []*DocumentNode {
	var documents []*DocumentNode
	var currentDoc *DocumentNode
	currentNode := node.FirstChild
	for currentNode != nil {
		nextNode := currentNode.Next
		if currentNode.Type == blackfriday.HorizontalRule {
			if currentDoc != nil {
				fillDocumentSettings(currentDoc)
				documents = append(documents, currentDoc)
				currentDoc = nil
			}
		} else {
			if currentDoc == nil {
				currentDoc = &DocumentNode{
					Node: blackfriday.Node{Type: blackfriday.Document},
				}
			}
			currentDoc.AppendChild(currentNode)
		}
		currentNode = nextNode
	}
	if currentDoc != nil {
		fillDocumentSettings(currentDoc)
		documents = append(documents, currentDoc)
	}
	cascadeFooters(documents)
	return documents
}

func cascadeFooters(nodes []*DocumentNode) {
	currentFooter := ""
	for _, n := range nodes {
		if n.HasFooter() {
			currentFooter = n.Footer()
		} else {
			n.Settings.Set("footer", currentFooter)
		}
	}
}

func fillDocumentSettings(node *DocumentNode) {
	node.Settings = url.Values{}
	c := node.FirstChild
	for c != nil {
		if c.Type == blackfriday.Paragraph {
			fc := c.FirstChild
			for fc != nil {
				if fc.Type == blackfriday.HTMLSpan {
					htmlNodes, _ := html.ParseFragment(bytes.NewReader(fc.Literal),
						&html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body},
					)
					for _, n := range htmlNodes {
						if n.Data != "meta" {
							continue
						}
						for _, a := range n.Attr {
							node.Settings.Set(a.Key, a.Val)
						}
					}
				}
				fc = fc.Next
			}
		}
		c = c.Next
	}
}
