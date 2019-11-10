package slides

import (
	"html/template"

	"github.com/russross/blackfriday"
)

type Settings struct {
	FooterText string
	XResPX     int
	YResPX     int
	HAlign     string
	VAlign     string
	TAlign     string
	FontColor  string
	Background template.CSS
}

type Slide struct {
	blackfriday.Node
	Settings
}
