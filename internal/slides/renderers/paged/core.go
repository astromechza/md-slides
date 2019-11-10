package paged

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/astromechza/md-slides/internal/slides"

	"github.com/russross/blackfriday"

	"github.com/astromechza/md-slides/internal/css"
	"github.com/astromechza/md-slides/internal/customhtml"
)

type Renderer struct {
	Path      string
	Source    slides.SlideSource
	Templates *template.Template
}

func New(path string, source slides.SlideSource) (*Renderer, error) {

	var err error
	root := template.New("")
	root = root.Option("missingkey=error")

	AddScriptTemplate(root)
	css.AddCommonStyleTemplate(root)
	css.AddMarkdownStyleTemplate(root)
	css.AddNormalizeStyleTemplate(root)
	css.AddChromaStyleTemplate(root)
	AddStyleOverridesTemplate(root)

	root, err = root.Parse(slideTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse root template %s", err)
	}

	return &Renderer{
		Path:      path,
		Source:    source,
		Templates: root,
	}, nil
}

func (sr *Renderer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	snRaw := req.URL.Query().Get("page")
	if snRaw == "" {
		rw.Header().Set("location", sr.FirstSlidePath())
		rw.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	sn, err := strconv.Atoi(snRaw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	sr.Serve(int(sn), rw, req)
}

func (sr *Renderer) NthSlidePath(n int) string {
	return fmt.Sprintf("%s?page=%d", sr.Path, n)
}

func (sr *Renderer) FirstSlidePath() string {
	return sr.NthSlidePath(1)
}

func (sr *Renderer) Serve(i int, rw http.ResponseWriter, req *http.Request) {
	collection, err := sr.Source.Load()
	if err != nil {
		log.Printf("Error: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte(err.Error()))
		return
	}
	if i < 1 || i > len(collection.Slides) {
		rw.Header().Set("location", sr.FirstSlidePath())
		rw.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	nextSlide, prevSlide := i, i
	if prevSlide > 1 {
		prevSlide--
	}
	if nextSlide < len(collection.Slides) {
		nextSlide++
	}
	currentSlide := collection.Slides[i-1]

	renderer := &customhtml.CustomRenderer{
		CWD: collection.WorkingDirectory,
		Renderer: blackfriday.Renderer(blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			Flags: blackfriday.HrefTargetBlank,
		})),
	}
	var b bytes.Buffer
	currentSlide.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return renderer.RenderNode(&b, node, entering)
	})
	if err := sr.Templates.Execute(rw, struct {
		PageNum   int
		PageNext  int
		PagePrev  int
		PageCount int

		URLPath  string
		Title    string
		Settings slides.Settings

		SlideContent template.HTML
	}{
		PageNum:   i,
		PageNext:  nextSlide,
		PagePrev:  prevSlide,
		PageCount: len(collection.Slides),

		URLPath:  sr.Path,
		Title:    collection.Title,
		Settings: currentSlide.Settings,

		SlideContent: template.HTML(b.String()),
	}); err != nil {
		log.Fatalf("error executing template: %s", err)
	}
}
