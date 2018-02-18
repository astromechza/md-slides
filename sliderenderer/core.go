package sliderenderer

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"

	"github.com/russross/blackfriday"
)

type SlideRenderer struct {
	Filename     string
	Templates    *template.Template
	CachedSlides []*DocumentNode
	Hot          bool
	XRes         int
	YRes         int
	FontSize     int
	BGCSS        string
	URLPath      string
}

func (sr *SlideRenderer) Init() (err error) {
	sr.Templates, err = LoadTemplates()
	return
}

func (sr *SlideRenderer) ShouldRecache() bool {
	return sr.CachedSlides == nil || sr.Hot
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
