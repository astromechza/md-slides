package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AstromechZA/md-slides/sliderenderer"
)

func parseResString(i string) (int, int, error) {
	i = strings.TrimSpace(strings.ToLower(i))
	parts := strings.Split(i, "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("res string '%s' did not contain one 'x'", i)
	}
	xres, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse x value of res string '%s': %s", i, err)
	}
	yres, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse y value of res string '%s': %s", i, err)
	}
	if xres <= 0 {
		return 0, 0, fmt.Errorf("x value of rest string '%s' should be > 0", i)
	}
	if yres <= 0 {
		return 0, 0, fmt.Errorf("y value of rest string '%s' should be > 0", i)
	}
	return int(xres), int(yres), nil
}

const slidesPath = "/_slides/"

func redirectToFirstSlide(rw http.ResponseWriter) {
	rw.Header().Set("location", fmt.Sprintf("%s0", slidesPath))
	rw.WriteHeader(http.StatusTemporaryRedirect)
	return
}

func Serve(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	hotFlag := fs.Bool("hot", false, "reload, reparse, and regenerate slides on each refresh")
	resFlag := fs.String("res", "1600x900", "set render aspect ratio and zoom for rendering")
	listenFlag := fs.String("listen", ":8080", "interface:port to listen on")
	backgroundCSS := fs.String("css-background", "#fffff8", "slide background css")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fs.Usage()
		fmt.Fprintf(os.Stderr, "\n")
		return fmt.Errorf("expected a single source file as argument")
	}
	filename := fs.Arg(0)

	xres, yres, err := parseResString(*resFlag)
	if err != nil {
		return fmt.Errorf("bad res string: %s", err)
	}
	mux := http.NewServeMux()

	sr := sliderenderer.SlideRenderer{Filename: filename, Hot: *hotFlag, XRes: xres, YRes: yres, BGCSS: *backgroundCSS}
	mux.HandleFunc(slidesPath, func(rw http.ResponseWriter, req *http.Request) {
		snRaw := req.URL.Path[len(slidesPath):]
		if snRaw == "" {
			redirectToFirstSlide(rw)
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

	statics := http.FileServer(CustomDirFS{Directory: filepath.Dir(filename)})
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			redirectToFirstSlide(rw)
			return
		}
		statics.ServeHTTP(rw, req)
	})

	log.Printf("Ready to serve on http://%s", *listenFlag)
	if err := http.ListenAndServe(*listenFlag, mux); err != nil {
		return err
	}
	return nil
}
