package renderers

import "github.com/AstromechZA/md-slides/pkg/slide"

type SlideSource interface {
	Load() (*slide.Collection, error)
}
