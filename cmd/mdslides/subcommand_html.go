package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/AstromechZA/md-slides/pkg/renderers/scrolling"
	"github.com/AstromechZA/md-slides/pkg/slide"
)

const htmlUsage = `Render the slides to html as an index.html file in the target directory.

Usage:
  md-slides html [options...] [file.md] [target directory]

`

func copyFile(srcPath, dstPath string) error {
	d, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %s", err)
	}
	defer d.Close()
	s, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %s", err)
	}
	defer s.Close()
	if _, err := io.Copy(d, s); err != nil {
		return fmt.Errorf("copy failed: %s", err)
	}
	return nil
}

func exportSlidesHTML(source *slide.CachedSource, targetDirectory string, noStatics bool) error {
	log.Printf("Generating html output to %s..", targetDirectory)
	urlPath := "/_slides"
	handler, err := scrolling.New(urlPath, source)
	if err != nil {
		return fmt.Errorf("failed to construct handler: %s", err)
	}

	f, err := os.Create(filepath.Join(targetDirectory, "index.html"))
	if err != nil {
		return fmt.Errorf("failed to open target file: %s", err)
	}
	defer f.Close()

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest("GET", urlPath, nil))
	if recorder.Code != 200 {
		return fmt.Errorf("failed to handle internal request: code %d", recorder.Code)
	}
	io.Copy(f, recorder.Body)

	if !noStatics {
		slides, _ := source.Load()
		statics, err := slides.CollectReferencedStaticFiles()
		if err != nil {
			return fmt.Errorf("failed to load static files")
		}
		for _, s := range statics {
			log.Printf("copying %s to output directory..", s)
			if err := copyFile(filepath.Join(slides.WorkingDirectory, s), filepath.Join(targetDirectory, s)); err != nil {
				return fmt.Errorf("failed to copy %s: %s", s, err)
			}
		}
	}
	log.Printf("Done.")
	return nil
}

func SubcommandHTML(args []string) error {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	noStaticFlag := fs.Bool("no-statics", false, "Don't serve static files")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, htmlUsage)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 2 {
		fs.Usage()
		fmt.Fprintf(os.Stderr, "\n")
		return fmt.Errorf("expected positional arguments")
	}

	sourceFileName, targetDirectoryName := fs.Arg(0), fs.Arg(1)
	slideSource := &slide.CachedSource{Inner: &slide.FileSource{Path: sourceFileName}}

	if err := exportSlidesHTML(slideSource, targetDirectoryName, *noStaticFlag); err != nil {
		return err
	}

	return nil
}
