package slide

import (
	"bytes"
	"log"

	"github.com/AstromechZA/md-slides/pkg/util"

	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Settings struct {
	FooterText string
	XResPX     int
	YResPX     int
	HAlign     string
	VAlign     string
	TAlign     string
	FontColor  string
}

type Slide struct {
	blackfriday.Node
	Settings
}

func NewSlide(dom blackfriday.Node, previous *Slide) *Slide {
	n := &Slide{Node: dom}

	// first copy the values from the previous node
	n.Settings = previous.Settings

	// and fill in any values from the meta
	meta := getMetaFromDom(n.Node)
	if v, ok := meta["halign"]; ok {
		n.HAlign = v
	}
	if v, ok := meta["valign"]; ok {
		n.VAlign = v
	}
	if v, ok := meta["talign"]; ok {
		n.TAlign = v
	}
	if v, ok := meta["res"]; ok {
		if v == "" {
			n.XResPX = 0
			n.YResPX = 0
		} else {
			x, y, err := util.ParseXYResString(v)
			if err != nil {
				log.Fatalf("failed to parse res '%s': %s", v, err)
			}
			if x <= 0 {
				log.Fatalf("invalid xres '%s': too small", x)
			}
			n.XResPX = int(x)
			if y <= 0 {
				log.Fatalf("invalid y res '%s': too small", y)
			}
			n.YResPX = int(y)
		}
	}
	if v, ok := meta["footer"]; ok {
		n.FooterText = v
	}
	if v, ok := meta["fontcolor"]; ok {
		n.FontColor = v
	}

	// TODO validate values?

	// if any values are zero set them to defaults
	if n.HAlign == "" {
		n.HAlign = "left"
	}
	if n.VAlign == "" {
		n.VAlign = "top"
	}
	if n.TAlign == "" {
		n.TAlign = "left"
	}
	if n.XResPX == 0 || n.YResPX == 0 {
		n.XResPX = 1366
		n.YResPX = 768
	}
	if n.FontColor == "" {
		n.FontColor = "#111111"
	}

	return n
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
