package renderers

import "github.com/astromechza/md-slides/pkg/slide"

type SlideSource interface {
	Load() (*slide.Collection, error)
}
