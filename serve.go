package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
</script>
`

const styleHeader = `
<style>
body {
	background-color: grey;
}

#body-inner {
	background-color: white;
	padding: 1rem;
	height: -webkit-fill-available;
    width: -webkit-fill-available;
    border-radius: 0.1rem;
    box-shadow: 0px 1px 3px black;
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
		log.Printf("Doing node %s (next=%s)", currentNode.Type.String(), nextNode)
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
		log.Printf("Rendering doc %d: %v", sn, doc)
		rndr := blackfriday.Renderer(blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			Flags: blackfriday.CompletePage,
		}))
		rndr = &CustomHTMLRenderer{Renderer: rndr}

		rndr.RenderHeader(rw, nil)
		rw.Write([]byte(fmt.Sprintf(scriptHeader, prevSlide, nextSlide)))
		rw.Write([]byte(styleHeader))
		rw.Write([]byte(`<div id="body-inner">`))
		doc.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
			return rndr.RenderNode(rw, node, entering)
		})
		rw.Write([]byte(`</div>`))
		rndr.RenderFooter(rw, nil)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}

	return nil
}
