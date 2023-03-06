package slides

import (
	"html/template"

	"github.com/russross/blackfriday/v2"
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
