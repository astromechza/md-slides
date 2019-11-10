package parse

import (
	"github.com/astromechza/md-slides/internal/slides"
)

type CachedSource struct {
	Inner      slides.SlideSource
	collection *slides.Collection
}

func (cs *CachedSource) Load() (*slides.Collection, error) {
	var err error
	if cs.collection == nil {
		if cs.collection, err = cs.Inner.Load(); err != nil {
			return nil, err
		}
	}
	return cs.collection, nil
}
