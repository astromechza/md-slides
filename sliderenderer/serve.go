package sliderenderer

import (
	"fmt"
	"net/http"
	"strconv"
)

const slidesPath = "/_slides/"

func (sr *SlideRenderer) InstallHandler(mux *http.ServeMux) {
	mux.HandleFunc(slidesPath, func(rw http.ResponseWriter, req *http.Request) {
		snRaw := req.URL.Path[len(slidesPath):]
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
	})
}

func (sr *SlideRenderer) FirstSlidePath() string {
	return fmt.Sprintf("%s0", slidesPath)
}
