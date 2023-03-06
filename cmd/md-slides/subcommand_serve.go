package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/astromechza/md-slides/internal/slides"
	"github.com/astromechza/md-slides/internal/slides/parse"
	"github.com/astromechza/md-slides/internal/slides/renderers/paged"
	"github.com/astromechza/md-slides/internal/slides/renderers/scrolling"
	"github.com/astromechza/md-slides/internal/util"
)

const serveUsage = `serve the slides

Usage:
  md-slides serve [options...] --source <file>

`

func SubcommandServe(args []string) error {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	sourceFlag := fs.String("source", "", "Path to the source markdown file")
	hotFlag := fs.Bool("hot", false, "Reload slides from disk on each request")
	modeFlag := fs.String("mode", "paged", "Mode to serve slides in (choose from paged, scrolling)")
	noStaticFlag := fs.Bool("no-statics", false, "Don't serve static files")
	listenAddressFlag := fs.String("listen", "127.0.0.1:8080", "Address to listen on")

	fs.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, serveUsage)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *sourceFlag == "" {
		return fmt.Errorf("expected a value for --source")
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("expected 0 positional arguments")
	}

	var slideSource slides.SlideSource
	slideSource = &parse.FileSource{Path: *sourceFlag}
	if *hotFlag == false {
		slideSource = &parse.CachedSource{Inner: slideSource}
	}

	var handler http.Handler
	var err error
	switch *modeFlag {
	case "paged":
		handler, err = paged.New("/", slideSource)
		if err != nil {
			return fmt.Errorf("failed to construct handler: %s", err)
		}
	case "scrolling":
		handler, err = scrolling.New("/", slideSource)
		if err != nil {
			return fmt.Errorf("failed to construct handler: %s", err)
		}
	default:
		return fmt.Errorf("unknown mode '%s'", *modeFlag)
	}

	http.Handle("/favicon.ico", http.NotFoundHandler())
	if !*noStaticFlag {
		log.Printf("Setting up static file server on / for %s", filepath.Dir(*sourceFlag))
		http.Handle("/", util.RootOrHandler(handler, http.FileServer(util.CustomDirFS{Directory: filepath.Dir(*sourceFlag)})))
	} else {
		http.Handle("/", handler)
	}

	log.Printf("Ready to serve on http://%s", *listenAddressFlag)
	if err := http.ListenAndServe(*listenAddressFlag, http.DefaultServeMux); err != nil {
		return err
	}

	return nil
}
