package parse

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/russross/blackfriday"

	"github.com/astromechza/md-slides/internal/slides"
)

type FileSource struct {
	Path string
}

func (fs *FileSource) Load() (*slides.Collection, error) {
	f, err := os.Open(fs.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var b bytes.Buffer
	if _, err := io.Copy(&b, f); err != nil {
		return nil, err
	}
	node := blackfriday.New(
		blackfriday.WithExtensions(
			blackfriday.CommonExtensions |
				blackfriday.Footnotes |
				blackfriday.NoEmptyLineBeforeBlock,
		),
	).Parse(b.Bytes())
	parsed, err := SlidesFromDom(node, nil)
	if err != nil {
		return nil, err
	}
	return &slides.Collection{
		WorkingDirectory: filepath.Dir(fs.Path),
		Title: filepath.Base(fs.Path),
		Slides: parsed,
	}, nil
}
