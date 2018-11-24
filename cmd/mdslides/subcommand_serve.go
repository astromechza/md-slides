package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AstromechZA/md-slides/pkg/renderers/paged"
	"github.com/AstromechZA/md-slides/pkg/renderers/scrolling"

	"github.com/AstromechZA/md-slides/pkg/slide"
	"github.com/AstromechZA/md-slides/pkg/util"
)

const serveUsage = `serve the slides

Usage:
  md-slides serve [options...] [file.md]

`

type slideSourcer interface {
	Load() (*slide.Collection, error)
}

func SubcommandServe(args []string) error {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	hotFlag := fs.Bool("hot", false, "Reload slides from disk on each request")
	modeFlag := fs.String("mode", "paged", "Mode to serve slides in (choose from paged, scrolling)")
	noStaticFlag := fs.Bool("no-statics", false, "Don't serve static files")
	listenAddressFlag := fs.String("listen", "127.0.0.1:8080", "Address to listen on")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, serveUsage)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() == 0 {
		fs.Usage()
		fmt.Fprintf(os.Stderr, "\n")
		return fmt.Errorf("expected slides source file as first positional argument")
	}

	filename := fs.Arg(0)
	var slideSource slideSourcer
	slideSource = &slide.FileSource{Path: filename}
	if *hotFlag == false {
		slideSource = &slide.CachedSource{Inner: slideSource}
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
		log.Printf("Setting up static file server on / for %s", filepath.Dir(filename))
		http.Handle("/", util.RootOrHandler(handler, http.FileServer(util.CustomDirFS{Directory: filepath.Dir(filename)})))
	} else {
		http.Handle("/", handler)
	}

	log.Printf("Ready to serve on %s", *listenAddressFlag)
	if err := http.ListenAndServe(*listenAddressFlag, http.DefaultServeMux); err != nil {
		return err
	}

	return nil
}
