package slide

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/russross/blackfriday"
)

type BufferSource struct {
	bytes.Buffer
}

func (bs *BufferSource) Load() (*Collection, error) {
	node := blackfriday.New(
		blackfriday.WithExtensions(
			blackfriday.CommonExtensions |
				blackfriday.Footnotes |
				blackfriday.NoEmptyLineBeforeBlock,
		),
	).Parse(bs.Bytes())
	slides, err := ParseSlidesFromDom(node, nil)
	if err != nil {
		return nil, err
	}
	collection := &Collection{
		Slides: slides,
	}
	return collection, nil
}

type FileSource struct {
	Path string
}

func (fs *FileSource) Load() (*Collection, error) {
	f, err := os.Open(fs.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var b bytes.Buffer
	if _, err := io.Copy(&b, f); err != nil {
		return nil, err
	}
	bs := &BufferSource{b}
	collection, err := bs.Load()
	if err != nil {
		return nil, err
	}
	collection.WorkingDirectory = filepath.Dir(fs.Path)
	collection.Title = filepath.Base(fs.Path)
	return collection, nil
}

type innerIFace interface {
	Load() (*Collection, error)
}

type CachedSource struct {
	Inner      innerIFace
	collection *Collection
}

func (cs *CachedSource) Load() (*Collection, error) {
	var err error
	if cs.collection == nil {
		if cs.collection, err = cs.Inner.Load(); err != nil {
			return nil, err
		}
	}
	return cs.collection, nil
}
