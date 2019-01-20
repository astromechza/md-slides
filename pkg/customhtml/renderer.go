package customhtml

import (
	"encoding/json"
	"io"
	"log"
	"net/url"
	"os/exec"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/russross/blackfriday"
)

type CustomRenderer struct {
	blackfriday.Renderer
	CWD string
	InSkip bool
}

const uncheckedSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path d="M5 2c-1.654 0-3 1.346-3 3v14c0 1.654 1.346 3 3 3h14c1.654 0 3-1.346 3-3v-14c0-1.654-1.346-3-3-3h-14zm19 3v14c0 2.761-2.238 5-5 5h-14c-2.762 0-5-2.239-5-5v-14c0-2.761 2.238-5 5-5h14c2.762 0 5 2.239 5 5z"/></svg>`
const checkedSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path d="M10.041 17l-4.5-4.319 1.395-1.435 3.08 2.937 7.021-7.183 1.422 1.409-8.418 8.591zm-5.041-15c-1.654 0-3 1.346-3 3v14c0 1.654 1.346 3 3 3h14c1.654 0 3-1.346 3-3v-14c0-1.654-1.346-3-3-3h-14zm19 3v14c0 2.761-2.238 5-5 5h-14c-2.762 0-5-2.239-5-5v-14c0-2.761 2.238-5 5-5h14c2.762 0 5 2.239 5 5z"/></svg>`

func (r *CustomRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {

	if entering && node.Type == blackfriday.HTMLSpan && strings.HasPrefix(string(node.Literal), "<meta norender>") {
		node.Unlink()
		return blackfriday.SkipChildren
	}

	if entering &&
		node.Type == blackfriday.Item &&
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

	if entering &&
		(node.Type == blackfriday.Text || node.Type == blackfriday.CodeBlock) {
		re := regexp.MustCompile(`{embedcommand: (.*?)}`)
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
				log.Fatalf("failed to execute embedcommand: %s", err)
			}
			cmdOutStr := strings.TrimSpace(string(cmdOut))
			node.Literal = []byte(strings.Replace(string(node.Literal), match[0], cmdOutStr, 1))
		}
	}

	if !entering &&
		node.Type == blackfriday.Image {
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

	if entering &&
		node.Type == blackfriday.CodeBlock {
		var lexer chroma.Lexer
		if len(node.CodeBlockData.Info) > 0 {
			lexer = lexers.Get(string(node.CodeBlockData.Info))
		}
		if lexer == nil {
			lexer = lexers.Fallback
		}
		lexer = chroma.Coalesce(lexer)

		// Tokenize the code
		iterator, err := lexer.Tokenise(nil, string(node.Literal))
		if err == nil {
			if html.New(html.WithClasses()).Format(w, styles.BlackWhite, iterator) == nil {
				return blackfriday.SkipChildren
			}
		}
		log.Printf("code formatter problem: %s", err)
	}

	return r.Renderer.RenderNode(w, node, entering)
}
