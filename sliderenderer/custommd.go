package sliderenderer

import (
	"encoding/json"
	"io"
	"log"
	"net/url"
	"os/exec"
	"regexp"
	"strings"

	"github.com/russross/blackfriday"
)

type CustomHTMLRenderer struct {
	blackfriday.Renderer
	CWD string
}

const uncheckedSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path d="M5 2c-1.654 0-3 1.346-3 3v14c0 1.654 1.346 3 3 3h14c1.654 0 3-1.346 3-3v-14c0-1.654-1.346-3-3-3h-14zm19 3v14c0 2.761-2.238 5-5 5h-14c-2.762 0-5-2.239-5-5v-14c0-2.761 2.238-5 5-5h14c2.762 0 5 2.239 5 5z"/></svg>`
const checkedSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path d="M10.041 17l-4.5-4.319 1.395-1.435 3.08 2.937 7.021-7.183 1.422 1.409-8.418 8.591zm-5.041-15c-1.654 0-3 1.346-3 3v14c0 1.654 1.346 3 3 3h14c1.654 0 3-1.346 3-3v-14c0-1.654-1.346-3-3-3h-14zm19 3v14c0 2.761-2.238 5-5 5h-14c-2.762 0-5-2.239-5-5v-14c0-2.761 2.238-5 5-5h14c2.762 0 5 2.239 5 5z"/></svg>`

func (r *CustomHTMLRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {

	if node.Type == blackfriday.Item &&
		node.FirstChild != nil &&
		node.FirstChild.Type == blackfriday.Paragraph &&
		node.FirstChild.FirstChild != nil &&
		node.FirstChild.FirstChild.Type == blackfriday.Text {
		if strings.HasPrefix(string(node.FirstChild.FirstChild.Literal), "[ ]") {
			node.FirstChild.FirstChild.Literal = node.FirstChild.FirstChild.Literal[3:]
			nn := blackfriday.NewNode(blackfriday.HTMLSpan)
			nn.Literal = []byte(uncheckedSVG)
			nn.Next = node.FirstChild
			node.FirstChild.Prev = nn
			node.FirstChild = nn
		} else if strings.HasPrefix(string(node.FirstChild.FirstChild.Literal), "[x]") {
			node.FirstChild.FirstChild.Literal = node.FirstChild.FirstChild.Literal[3:]
			nn := blackfriday.NewNode(blackfriday.HTMLSpan)
			nn.Literal = []byte(checkedSVG)
			nn.Next = node.FirstChild
			node.FirstChild.Prev = nn
			node.FirstChild = nn
		}
	}

	if node.Type == blackfriday.Text || node.Type == blackfriday.CodeBlock {
		re := regexp.MustCompile(`\{embedcommand: (.*?)\}`)
		matches := re.FindAllStringSubmatch(string(node.Literal), -1)
		for _, match := range matches {
			var cmd []string
			if err := json.Unmarshal([]byte(match[1]), &cmd); err != nil {
				panic(err)
			}
			log.Printf("Attempting to execute '%v'", cmd)
			c := exec.Command(cmd[0], cmd[1:]...)
			c.Dir = r.CWD
			cmdOut, err := c.CombinedOutput()
			if err != nil {
				panic(err)
			}
			node.Literal = []byte(strings.Replace(string(node.Literal), match[0], string(cmdOut), 1))
		}
	}

	if node.Type == blackfriday.Image && entering == false {
		if u, err := url.Parse(string(node.LinkData.Destination)); err == nil {
			extra, _ := url.ParseQuery(u.Fragment)
			w.Write([]byte(`" style="`))
			if extra.Get("height") != "" {
				w.Write([]byte(`height: ` + extra.Get("height") + ";"))
			}
			if extra.Get("width") != "" {
				w.Write([]byte(`width: ` + extra.Get("width") + ";"))
			}
		}
	}

	return r.Renderer.RenderNode(w, node, entering)
}
