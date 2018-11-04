package scrolling

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/AstromechZA/md-slides/pkg/slide"

	"github.com/russross/blackfriday"

	"github.com/AstromechZA/md-slides/pkg/css"
	"github.com/AstromechZA/md-slides/pkg/customhtml"
	"github.com/AstromechZA/md-slides/pkg/renderers"
)

type Renderer struct {
	Path      string
	Source    renderers.SlideSource
	Templates *template.Template
}

func New(path string, source renderers.SlideSource) (*Renderer, error) {

	var err error
	root := template.New("")
	root = root.Option("missingkey=error")

	css.AddCommonStyleTemplate(root)
	css.AddMarkdownStyleTemplate(root)
	css.AddNormalizeStyleTemplate(root)
	css.AddChromaStyleTemplate(root)

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

type preparedSlide struct {
	PageNum  int
	Content  template.HTML
	Settings slide.Settings
	PageLeft int
	PageTop  int
	Scale    float32
}

func (sr *Renderer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	collection, err := sr.Source.Load()
	if err != nil {
		log.Printf("Error: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}

	renderer := &customhtml.CustomRenderer{
		CWD: collection.WorkingDirectory,
		Renderer: blackfriday.Renderer(blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			Flags: blackfriday.HrefTargetBlank,
		})),
	}

	var preparedSlides []*preparedSlide
	for i, s := range collection.Slides {
		var b bytes.Buffer
		s.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
			return renderer.RenderNode(&b, node, entering)
		})
		preparedSlides = append(preparedSlides, &preparedSlide{
			PageNum:  i + 1,
			Content:  template.HTML(b.String()),
			Settings: s.Settings,
		})
	}

	pageXRes := preparedSlides[0].Settings.XResPX
	pageYRes := preparedSlides[0].Settings.YResPX

	for _, s := range preparedSlides {
		s.Scale = float32(pageXRes) / float32(s.Settings.XResPX)
		if ys := float32(pageYRes) / float32(s.Settings.YResPX); (s.Scale > 0 && ys < s.Scale) || (s.Scale < 0 && ys > s.Scale) {
			s.Scale = ys
		}
		s.PageLeft = (pageXRes-s.Settings.XResPX)/2 + 20
		s.PageTop = (pageYRes-s.Settings.YResPX)/2 + 20
	}

	if err := sr.Templates.Execute(rw, struct {
		PageCount int

		URLPath            string
		Title              string
		PageXResPX         int
		PageYResPX         int
		AdjustedPageYResPX int
		PreparedSlides     []*preparedSlide
	}{
		PageCount: len(collection.Slides),

		URLPath:            sr.Path,
		Title:              collection.Title,
		PageXResPX:         pageXRes + 40,
		PageYResPX:         pageYRes + 40,
		AdjustedPageYResPX: pageYRes + 40 - 1,
		PreparedSlides:     preparedSlides,
	}); err != nil {
		log.Fatalf("error executing template: %s", err)
	}
}
