package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/russross/blackfriday"
)

type SlideRenderer struct {
	Filename     string
	CachedSlides []*blackfriday.Node
	Hot          bool
	XRes         int
	YRes         int
}

func (sr *SlideRenderer) ShouldRecache() bool {
	return sr.CachedSlides == nil || sr.Hot
}

func breakIntoDocumentNodes(node *blackfriday.Node) []*blackfriday.Node {
	var documents []*blackfriday.Node
	var currentDoc *blackfriday.Node
	currentNode := node.FirstChild
	for currentNode != nil {
		nextNode := currentNode.Next
		if currentNode.Type == blackfriday.HorizontalRule {
			if currentDoc != nil {
				documents = append(documents, currentDoc)
				currentDoc = nil
			}
		} else {
			if currentDoc == nil {
				currentDoc = &blackfriday.Node{Type: blackfriday.Document}
			}
			currentDoc.AppendChild(currentNode)
		}
		currentNode = nextNode
	}
	if currentDoc != nil {
		documents = append(documents, currentDoc)
	}
	return documents
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
			blackfriday.CommonExtensions | blackfriday.NoEmptyLineBeforeBlock,
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
	if i < 0 || i >= len(sr.CachedSlides) {
		rw.Header().Set("location", "/slides/0")
		rw.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	nextSlide, prevSlide := i, i
	if prevSlide > 0 {
		prevSlide--
	}
	if nextSlide < len(sr.CachedSlides)-1 {
		nextSlide++
	}

	doc := sr.CachedSlides[i]
	rndr := blackfriday.Renderer(blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Title: filepath.Base(sr.Filename),
		Flags: blackfriday.CompletePage | blackfriday.HrefTargetBlank,
	}))
	rndr = &CustomHTMLRenderer{Renderer: rndr}

	rndr.RenderHeader(rw, nil)
	rw.Write([]byte(fmt.Sprintf(scriptHeader, prevSlide, nextSlide)))
	rw.Write([]byte(normalizeCSS))
	rw.Write([]byte(styleHeader))
	rw.Write([]byte(markdownCSS))
	rw.Write([]byte(fmt.Sprintf(`<div id="body-inner" style="width: %dpx; height: %dpx;">`, sr.XRes, sr.YRes)))
	rw.Write([]byte(`<div class="markdown-body">`))
	doc.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return rndr.RenderNode(rw, node, entering)
	})
	rw.Write([]byte(`</div>`))
	rw.Write([]byte(`</div>`))
	rndr.RenderFooter(rw, nil)
}
