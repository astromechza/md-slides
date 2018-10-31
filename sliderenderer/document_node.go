package sliderenderer

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type DocumentNode struct {
	blackfriday.Node
	Settings map[string]string
}

func (n *DocumentNode) SlideClasses() string {
	bodyClasses := []string{"slide-wrap"}
	for _, m := range []string{"halign", "valign", "talign"} {
		if k, ok := n.Settings[m]; ok && k != "" {
			bodyClasses = append(bodyClasses, fmt.Sprintf("slide-wrap-%s-%s", m, k))
		}
	}
	return strings.Join(bodyClasses, " ")
}

func (n *DocumentNode) Footer() string {
	return n.Settings["footer"]
}

func (n *DocumentNode) FillMetaSettingsFromSelf() {
	n.Settings = make(map[string]string)
	c := n.FirstChild
	for c != nil {
		if c.Type == blackfriday.Paragraph {
			fc := c.FirstChild
			for fc != nil {
				if fc.Type == blackfriday.HTMLSpan {
					htmlNodes, _ := html.ParseFragment(bytes.NewReader(fc.Literal),
						&html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body},
					)
					for _, hn := range htmlNodes {
						if hn.Data != "meta" {
							continue
						}
						for _, a := range hn.Attr {
							n.Settings[a.Key] = a.Val
						}
					}
				}
				fc = fc.Next
			}
		}
		c = c.Next
	}
}

func (n *DocumentNode) FillMetaSettingsFromParent(other *DocumentNode) {
	// loop through parents settings
	for k, v := range other.Settings {
		// check if we have a value
		_, ok := n.Settings[k]

		// if we don't, then copy the last value to ourselves
		if !ok {
			n.Settings[k] = v
		}

		// if our value was or is now empty we can delete it entirely
		if n.Settings[k] == "" {
			delete(n.Settings, k)
		}
	}
}

func ConvertRootIntoDocumentNodes(node *blackfriday.Node) []*DocumentNode {
	var documents []*DocumentNode
	var currentDoc *DocumentNode
	var previousDoc *DocumentNode

	currentNode := node.FirstChild
	for currentNode != nil {
		nextNode := currentNode.Next
		if currentNode.Type == blackfriday.HorizontalRule {
			if currentDoc != nil {
				currentDoc.FillMetaSettingsFromSelf()
				if previousDoc != nil {
					currentDoc.FillMetaSettingsFromParent(previousDoc)
				}
				documents = append(documents, currentDoc)
				previousDoc = currentDoc
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
		currentDoc.FillMetaSettingsFromSelf()
		if previousDoc != nil {
			currentDoc.FillMetaSettingsFromParent(previousDoc)
		}
		documents = append(documents, currentDoc)
	}
	return documents
}
