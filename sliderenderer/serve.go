package sliderenderer

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

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
