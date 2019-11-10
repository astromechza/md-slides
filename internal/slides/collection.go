package slides

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday"
)

type Collection struct {
	Title            string
	WorkingDirectory string
	Slides           []*Slide
}

func buildLinkChecker(workingDirectory string) func(string) (string, error) {
	return func(href string) (string, error) {
		hrefParsed, err := url.Parse(href)
		if err != nil {
			return "", fmt.Errorf("invalid url")
		}
		if hrefParsed.Scheme != "" {
			return "", fmt.Errorf("has a scheme")
		}
		fullpath := filepath.Clean(filepath.Join(workingDirectory, hrefParsed.Path))
		relpath, err := filepath.Rel(workingDirectory, fullpath)
		if err != nil {
			return "", fmt.Errorf("could not be made relative")
		}
		if strings.HasPrefix(relpath, "..") {
			return "", fmt.Errorf("is in a parent directory")
		}
		s, err := os.Stat(fullpath)
		if err != nil {
			return "", fmt.Errorf("does not exist")
		}
		if s.IsDir() {
			return "", fmt.Errorf("is a directory")
		}
		return relpath, nil
	}
}

func (c *Collection) CollectReferencedStaticFiles() ([]string, error) {
	seen := make(map[string]bool)
	f := buildLinkChecker(c.WorkingDirectory)
	for _, s := range c.Slides {
		s.Node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
			switch node.Type {
			case blackfriday.Image:
				href := string(node.LinkData.Destination)
				if p, err := f(href); err == nil {
					seen[p] = true
				}
			default:
			}
			return blackfriday.GoToNext
		})
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out, nil
}
