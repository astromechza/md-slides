package parse

import (
	"bytes"
	"html/template"
	"log"

	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/astromechza/md-slides/internal/slides"
	"github.com/astromechza/md-slides/internal/util"
)

func SlidesFromDom(node *blackfriday.Node, previousSettings *slides.Settings) ([]*slides.Slide, error) {
	var currentSlides []*slides.Slide
	var currentDoc *blackfriday.Node

	currentNode := node.FirstChild
	for currentNode != nil {
		nextNode := currentNode.Next

		if currentNode.Type == blackfriday.HorizontalRule {
			if currentDoc != nil {
				newSettings := calculateSettings(*currentDoc, previousSettings)
				currentSlide := &slides.Slide{Node: *currentDoc, Settings: newSettings}
				currentSlides = append(currentSlides, currentSlide)
				currentDoc = nil
				previousSettings = &currentSlide.Settings
			}
		} else {
			if currentDoc == nil {
				currentDoc = &blackfriday.Node{Type: blackfriday.Document}
			}
			currentDoc.AppendChild(currentNode)
		}

		currentNode = nextNode
	}

	if currentDoc != nil {
		newSettings := calculateSettings(*currentDoc, previousSettings)
		currentSlide := &slides.Slide{Node: *currentDoc, Settings: newSettings}
		currentSlides = append(currentSlides, currentSlide)
	}
	return currentSlides, nil
}


func calculateSettings(dom blackfriday.Node, previous *slides.Settings) slides.Settings {
	if previous == nil {
		previous = new(slides.Settings)
	}
	output := *previous

	// and fill in any values from the meta
	meta := getMetaFromDom(dom)
	if v, ok := meta["halign"]; ok {
		output.HAlign = v
	}
	if v, ok := meta["valign"]; ok {
		output.VAlign = v
	}
	if v, ok := meta["talign"]; ok {
		output.TAlign = v
	}
	if v, ok := meta["res"]; ok {
		if v == "" {
			output.XResPX = 0
			output.YResPX = 0
		} else {
			x, y, err := util.ParseXYResString(v)
			if err != nil {
				log.Fatalf("failed to parse res '%s': %s", v, err)
			}
			output.XResPX = int(x)
			output.YResPX = int(y)
		}
	}
	if v, ok := meta["footer"]; ok {
		output.FooterText = v
	}
	if v, ok := meta["fontcolor"]; ok {
		output.FontColor = v
	}
	if v, ok := meta["background"]; ok {
		output.Background = template.CSS(v)
	}

	// TODO validate values?

	// if any values are zero set them to defaults
	if output.HAlign == "" {
		output.HAlign = "left"
	}
	if output.VAlign == "" {
		output.VAlign = "top"
	}
	if output.TAlign == "" {
		output.TAlign = "left"
	}
	if output.XResPX == 0 || output.YResPX == 0 {
		output.XResPX = 1366
		output.YResPX = 768
	}
	if output.FontColor == "" {
		output.FontColor = "#111111"
	}
	if output.Background == "" {
		output.Background = "#fffff8"
	}
	return output
}

func getMetaFromDom(n blackfriday.Node) map[string]string {
	out := make(map[string]string)
	c := n.FirstChild
	for c != nil {
		if c.Type == blackfriday.Paragraph {
			fc := c.FirstChild
			for fc != nil {
				if fc.Type == blackfriday.HTMLSpan {
					htmlNodes, _ := html.ParseFragment(bytes.NewReader(fc.Literal),
						&html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body},
					)
					for _, hn := range htmlNodes {
						if hn.Data != "meta" {
							continue
						}
						for _, a := range hn.Attr {
							out[a.Key] = a.Val
						}
					}
				}
				fc = fc.Next
			}
		}
		c = c.Next
	}
	return out
}
