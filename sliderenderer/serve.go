package sliderenderer

import (
	"fmt"
	"net/http"
	"strconv"
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

func (sr *SlideRenderer) NthSlidePath(n int) string {
	return fmt.Sprintf("%s?page=%d", sr.URLPath, n)
}

func (sr *SlideRenderer) FirstSlidePath() string {
	return sr.NthSlidePath(1)
}
