package sliderenderer

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/russross/blackfriday"
)

func (sr *SlideRenderer) MultiServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if sr.ShouldRecache() {
		if err := sr.RecacheSlides(); err != nil {
			log.Printf("Error: %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
	}

	rndr := blackfriday.Renderer(blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Title: filepath.Base(sr.Filename),
		Flags: blackfriday.CompletePage | blackfriday.HrefTargetBlank,
	}))
	rndr = &CustomHTMLRenderer{Renderer: rndr, CWD: filepath.Dir(sr.Filename)}
	rndr.RenderHeader(rw, nil)

	sr.Templates.ExecuteTemplate(rw, "normalize-css", sr)
	sr.Templates.ExecuteTemplate(rw, "markdown-css", sr)
	sr.Templates.ExecuteTemplate(rw, "other-css", sr)
	sr.Templates.ExecuteTemplate(rw, "multipage-css", sr)

	for i, doc := range sr.CachedSlides {
		ctx := struct {
			SlideRenderer
			PageNum      int
			SlideClasses string
		}{
			SlideRenderer: *sr, PageNum: i + 1, SlideClasses: doc.SlideClasses(),
		}

		sr.Templates.ExecuteTemplate(rw, "slide-prefix", ctx)
		doc.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
			return rndr.RenderNode(rw, node, entering)
		})
		sr.Templates.ExecuteTemplate(rw, "slide-suffix", ctx)
	}

	rndr.RenderFooter(rw, nil)
}
