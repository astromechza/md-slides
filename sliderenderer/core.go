package sliderenderer

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"strconv"

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

func NewSlideRenderer(filename string) (sr *SlideRenderer, err error) {
	sr = &SlideRenderer{
		Filename: filename,
		Hot:      false,
		XRes:     1366,
		YRes:     768,
		FontSize: 18,
		BGCSS:    "#fffff8",
		URLPath:  "/_slides",
	}
	if sr.Templates, err = LoadTemplates(); err != nil {
		return nil, err
	}
	return sr, nil
}

func (sr *SlideRenderer) ShouldRecache() bool {
	return sr.CachedSlides == nil || sr.Hot
}

func (sr *SlideRenderer) RecacheSlides() error {
	log.Printf("Reading file %s..", sr.Filename)
	content, err := ioutil.ReadFile(sr.Filename)
	if err != nil {
		return fmt.Errorf("failed to read '%s': %s", sr.Filename, err)
	}
	log.Printf("Parsing Markdown..")
	node := blackfriday.New(
		blackfriday.WithExtensions(
			blackfriday.CommonExtensions |
				blackfriday.Footnotes |
				blackfriday.NoEmptyLineBeforeBlock,
		),
	).Parse(content)
	log.Printf("Breaking into slides..")
	documents := ConvertRootIntoDocumentNodes(node)
	log.Printf("Loading top level settings..")

	if fs, ok := documents[0].Settings["font-size"]; ok {
		i, err := strconv.ParseInt(fs, 10, 63)
		if err != nil {
			return fmt.Errorf("top level font size from slide 0 is invalid: %s", err)
		}
		sr.FontSize = int(i)
	}

	if r, ok := documents[0].Settings["res"]; ok {
		x, y, err := ParseResString(r)
		if err != nil {
			return fmt.Errorf("top level res from slide 0 is invalid: %s", err)
		}
		sr.XRes = x
		sr.YRes = y
	}

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
