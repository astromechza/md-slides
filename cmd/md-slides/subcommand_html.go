package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/astromechza/md-slides/internal/slides/parse"
	"github.com/astromechza/md-slides/internal/slides/renderers/scrolling"
)

const htmlUsage = `Render the slides to html as an index.html file in the target directory.

Usage:
  md-slides html [options...] --source <file> --target-dir <directory>

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

func exportSlidesHTML(source *parse.CachedSource, targetDirectory string, noStatics bool) error {
	log.Printf("Generating html output to %s..", targetDirectory)
	handler, err := scrolling.New("/", source)
	if err != nil {
		return fmt.Errorf("failed to construct handler: %s", err)
	}

	f, err := os.Create(filepath.Join(targetDirectory, "index.html"))
	if err != nil {
		return fmt.Errorf("failed to open target file: %s", err)
	}
	defer f.Close()

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest("GET", "/", nil))
	if recorder.Code != 200 {
		return fmt.Errorf("failed to handle internal request: code %d", recorder.Code)
	}
	_, err = io.Copy(f, recorder.Body)
	if err != nil {
		return fmt.Errorf("failed to copy all bytes: %s", err)
	}

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
	sourceFlag := fs.String("source", "", "Path to the source markdown file")
	destFlag := fs.String("target-dir", "", "Path to the directory to write output files to")
	noStaticFlag := fs.Bool("no-statics", false, "Don't serve static files")

	fs.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, htmlUsage)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *sourceFlag == "" {
		return fmt.Errorf("expected a value for --source")
	}
	if *destFlag == "" {
		return fmt.Errorf("expected a value for --target-dir")
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("expected 0 positional arguments")
	}

	slideSource := &parse.CachedSource{Inner: &parse.FileSource{Path: *sourceFlag}}
	if err := exportSlidesHTML(slideSource, *destFlag, *noStaticFlag); err != nil {
		return err
	}

	return nil
}
