package sliderenderer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type SlideRenderer struct {
	Filename     string
	CachedSlides []*DocumentNode
	Hot          bool
	XRes         int
	YRes         int
	BGCSS        string
	URLPath      string
}

func (sr *SlideRenderer) ShouldRecache() bool {
	return sr.CachedSlides == nil || sr.Hot
}

type DocumentNode struct {
	blackfriday.Node
	Settings url.Values
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
	return documents
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

func (sr *SlideRenderer) RecacheSlides() error {
	log.Printf("Reading file %s", sr.Filename)
	content, err := ioutil.ReadFile(sr.Filename)
	if err != nil {
		return fmt.Errorf("failed to read '%s': %s", sr.Filename, err)
	}
	log.Printf("Parsing Markdown")
	node := blackfriday.New(
		blackfriday.WithExtensions(
			blackfriday.CommonExtensions | blackfriday.Footnotes | blackfriday.NoEmptyLineBeforeBlock,
		),
	).Parse(content)
	log.Printf("Breaking into slides")
	documents := breakIntoDocumentNodes(node)
	log.Printf("Preprocessing done (%d slides)", len(documents))
	sr.CachedSlides = documents
	return nil
}

func (sr *SlideRenderer) NthSlidePath(n int) string {
	return fmt.Sprintf("%s?page=%d", sr.URLPath, n)
}

func (sr *SlideRenderer) FirstSlidePath() string {
	return sr.NthSlidePath(1)
}
