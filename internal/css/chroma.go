package css

import (
	"bytes"
	"html/template"
	"log"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/styles"
)

func AddChromaStyleTemplate(root *template.Template) {
	var b bytes.Buffer
	err := html.New(html.WithClasses()).WriteCSS(&b, styles.BlackWhite)
	if err != nil {
		log.Fatal("failed to generate code highlighter")
	}
	if _, err := root.New("style.chroma").Parse(b.String()); err != nil {
		log.Fatalf("failed to parse: %s", err)
	}
}
