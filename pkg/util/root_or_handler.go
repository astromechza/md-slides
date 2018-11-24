package util

import (
	"net/http"
	"path"
)

type rootOrOtherHandler struct {
	rootHandler  http.Handler
	otherHandler http.Handler
}

func RootOrHandler(rootHandler http.Handler, otherHandler http.Handler) http.Handler {
	return &rootOrOtherHandler{
		rootHandler:  rootHandler,
		otherHandler: otherHandler,
	}
}

func (h *rootOrOtherHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	p := path.Clean(req.URL.Path)
	if p == "/" {
		h.rootHandler.ServeHTTP(rw, req)
	} else {
		h.otherHandler.ServeHTTP(rw, req)
	}
}
