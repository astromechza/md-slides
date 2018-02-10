package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/russross/blackfriday"
)

const scriptHeader = `
<script>
var prevSlide = "/slides/%d";
var nextSlide = "/slides/%d";

document.onkeydown = function(evt) {
	evt = evt || window.event
	if (evt.keyCode == 39) {
		window.location = nextSlide;
	}
	if (evt.keyCode == 37) {
		window.location = prevSlide;
	}
}


window.onresize = function(event) {
	var el = document.getElementById("body-inner");
	var wi = 1600 + 70;
	var hi = 900 + 40;

	var ws = window.innerWidth / wi;
	var hs = window.innerHeight / hi;
	var ss = Math.min(ws, hs);
	el.style.transform = "scale(" + ss + ")";
};

document.addEventListener("DOMContentLoaded", function(event) {
    window.onresize(null);
});

</script>
`

const styleHeader = `
<style>
html {
	height: 100%;
	font-size: 20px;
}

body {
	height: 100%;
    display: flex;
    flex-flow: column;
	background-color: grey;
	justify-content: center;
}

#body-inner {
	display: flex;
	flex-flow: column;
	align-self: center;

	background-color: #fffff8;
	padding: 10px;
    border-radius: 0.1rem;
	box-shadow: 0px 0.2rem 0.6rem black;
	padding-left: 25px;
    padding-right: 25px;
	position: absolute;
	overflow: hidden;

	width: 1600px;
    height: 900px;
}

#body-inner.centered {
	justify-content: center;
}
</style>
`

func Serve(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return fmt.Errorf("expected a single source file as argument")
	}

	content, err := ioutil.ReadFile(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("failed to read '%s': %s", fs.Arg(0), err)
	}

	var documents []blackfriday.Node
	node := blackfriday.New(
		blackfriday.WithExtensions(blackfriday.CommonExtensions),
	).Parse(content)
	currentNode := node.FirstChild

	var currentDoc *blackfriday.Node

	for currentNode != nil {
		nextNode := currentNode.Next
		if currentNode.Type == blackfriday.HorizontalRule {
			if currentDoc != nil {
				documents = append(documents, *currentDoc)
				currentDoc = nil
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
		documents = append(documents, *currentDoc)
	}

	fmt.Printf("num slides %d\n", len(documents))
	http.HandleFunc("/slides/", func(rw http.ResponseWriter, req *http.Request) {
		snRaw := req.URL.Path[len("/slides/"):]
		if snRaw == "" {
			snRaw = "0"
		}
		sn, err := strconv.ParseInt(snRaw, 10, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(http.StatusText(http.StatusBadRequest)))
			return
		}
		if sn < 0 || sn >= int64(len(documents)) {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte(http.StatusText(http.StatusNotFound)))
			return
		}
		nextSlide, prevSlide := sn, sn
		if sn > 0 {
			prevSlide--
		}
		if nextSlide < int64(len(documents))-1 {
			nextSlide++
		}

		doc := documents[sn]
		rndr := blackfriday.Renderer(blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			Flags: blackfriday.CompletePage,
		}))
		rndr = &CustomHTMLRenderer{Renderer: rndr}

		rndr.RenderHeader(rw, nil)
		rw.Write([]byte(fmt.Sprintf(scriptHeader, prevSlide, nextSlide)))
		rw.Write([]byte(normalizeCSS))
		rw.Write([]byte(styleHeader))
		rw.Write([]byte(markdownCSS))
		rw.Write([]byte(`<div id="body-inner">`))
		rw.Write([]byte(`<div class="markdown-body">`))
		doc.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
			return rndr.RenderNode(rw, node, entering)
		})
		rw.Write([]byte(`</div>`))
		rw.Write([]byte(`</div>`))
		rndr.RenderFooter(rw, nil)
	})
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		return err
	}

	return nil
}
