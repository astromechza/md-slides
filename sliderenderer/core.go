package sliderenderer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

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

func (sr *SlideRenderer) Serve(i int, rw http.ResponseWriter, req *http.Request) {
	if sr.ShouldRecache() {
		if err := sr.RecacheSlides(); err != nil {
			log.Printf("Error: %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
	}
	if i < 1 || i > len(sr.CachedSlides) {
		rw.Header().Set("location", sr.FirstSlidePath())
		rw.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	nextSlide, prevSlide := i, i
	if prevSlide > 1 {
		prevSlide--
	}
	if nextSlide < len(sr.CachedSlides) {
		nextSlide++
	}

	doc := sr.CachedSlides[i-1]
	rndr := blackfriday.Renderer(blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Title: fmt.Sprintf("%s %d/%d", filepath.Base(sr.Filename), i, len(sr.CachedSlides)),
		Flags: blackfriday.CompletePage | blackfriday.HrefTargetBlank,
	}))
	rndr = &CustomHTMLRenderer{Renderer: rndr, CWD: filepath.Dir(sr.Filename)}
	rndr.RenderHeader(rw, nil)
	rw.Write([]byte(fmt.Sprintf(scriptHeader, prevSlide, nextSlide)))
	rw.Write([]byte(normalizeCSS))
	rw.Write([]byte(fmt.Sprintf(styleHeader, sr.BGCSS)))
	rw.Write([]byte(markdownCSS))
	bodyClasses := []string{"body-inner"}
	if doc.Settings.Get("halign") != "" {
		bodyClasses = append(bodyClasses, "body-inner-halign-"+doc.Settings.Get("halign"))
	}
	if doc.Settings.Get("valign") != "" {
		bodyClasses = append(bodyClasses, "body-inner-valign-"+doc.Settings.Get("valign"))
	}
	if doc.Settings.Get("talign") != "" {
		bodyClasses = append(bodyClasses, "body-inner-talign-"+doc.Settings.Get("talign"))
	}
	rw.Write([]byte(fmt.Sprintf(
		`<div id="body-inner" class="%s" style="width: %dpx; height: %dpx">`,
		strings.Join(bodyClasses, " "), sr.XRes, sr.YRes,
	)))
	rw.Write([]byte(`<div class="markdown-body">`))
	doc.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return rndr.RenderNode(rw, node, entering)
	})
	rw.Write([]byte(`</div>`))
	rw.Write([]byte(fmt.Sprintf(`<div class="page-number">%d/%d</div>`, i, len(sr.CachedSlides))))
	rw.Write([]byte(`</div>`))
	rndr.RenderFooter(rw, nil)
}
