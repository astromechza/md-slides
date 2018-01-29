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
	node := blackfriday.New().Parse(content)
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

		doc := documents[sn]
		log.Printf("Rendering doc %d: %v", sn, doc)
		rndr := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			Flags: blackfriday.CompletePage,
		})
		rndr.RenderHeader(rw, nil)
		doc.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
			return rndr.RenderNode(rw, node, entering)
		})
		rndr.RenderFooter(rw, nil)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}

	return nil
}
