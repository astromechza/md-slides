package slide

import "github.com/russross/blackfriday"

type Collection struct {
	Title            string
	WorkingDirectory string
	Slides           []*Slide
}

func ParseSlidesFromDom(node *blackfriday.Node, anchorNode *Slide) ([]*Slide, error) {
	if anchorNode == nil {
		anchorNode = new(Slide)
	}
	var slides []*Slide

	var currentDoc *blackfriday.Node
	var previousSlide = anchorNode

	currentNode := node.FirstChild
	for currentNode != nil {
		nextNode := currentNode.Next

		if currentNode.Type == blackfriday.HorizontalRule {
			if currentDoc != nil {
				currentSlide := NewSlide(*currentDoc, previousSlide)
				slides = append(slides, currentSlide)
				currentDoc = nil
				previousSlide = currentSlide
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
		currentSlide := NewSlide(*currentDoc, previousSlide)
		slides = append(slides, currentSlide)
		previousSlide = currentSlide
	}
	return slides, nil
}
