package sliderenderer

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/russross/blackfriday"
)

func (sr *SlideRenderer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	snRaw := req.URL.Query().Get("page")
	if snRaw == "" {
		rw.Header().Set("location", sr.FirstSlidePath())
		rw.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	sn, err := strconv.Atoi(snRaw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	sr.Serve(int(sn), rw, req)
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

	ctx := struct {
		SlideRenderer
		PageNum      int
		PageNext     int
		PagePrev     int
		SlideClasses string
	}{
		*sr, i, nextSlide, prevSlide, doc.SlideClasses(),
	}

	sr.Templates.ExecuteTemplate(rw, "scripts", ctx)
	sr.Templates.ExecuteTemplate(rw, "normalize-css", ctx)
	sr.Templates.ExecuteTemplate(rw, "markdown-css", ctx)
	sr.Templates.ExecuteTemplate(rw, "chroma-css", ctx)
	sr.Templates.ExecuteTemplate(rw, "other-css", ctx)
	sr.Templates.ExecuteTemplate(rw, "slide-prefix", ctx)
	doc.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return rndr.RenderNode(rw, node, entering)
	})
	sr.Templates.ExecuteTemplate(rw, "slide-suffix", ctx)
	rndr.RenderFooter(rw, nil)
}
