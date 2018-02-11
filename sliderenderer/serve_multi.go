package sliderenderer

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday"
)

func (sr *SlideRenderer) InstallMultiSlideHandler(mux *http.ServeMux) {
	mux.HandleFunc("/_multislide/", func(rw http.ResponseWriter, req *http.Request) {
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
		rw.Write([]byte(normalizeCSS))
		rw.Write([]byte(fmt.Sprintf(styleHeader, sr.BGCSS)))
		rw.Write([]byte(styleMultiHeader))
		rw.Write([]byte(markdownCSS))
		for _, doc := range sr.CachedSlides {
			bodyClasses := []string{"body-inner", "body-inner-multipage"}
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
				`<div id="body-inner" class="%s" style="width: %dpx; min-height: %dpx;">`,
				strings.Join(bodyClasses, " "), sr.XRes, sr.YRes,
			)))
			rw.Write([]byte(`<div class="markdown-body">`))
			doc.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
				return rndr.RenderNode(rw, node, entering)
			})
			rw.Write([]byte(`</div>`))
			rw.Write([]byte(`</div>`))
		}

		rndr.RenderFooter(rw, nil)

	})
}
