package css

import (
	"html/template"
	"log"
)

func AddMarkdownStyleTemplate(root *template.Template) {
	if _, err := root.New("style.markdown").Parse(`
@charset "UTF-8";

/* Import ET Book styles
   adapted from https://github.com/edwardtufte/et-book/blob/gh-pages/et-book.css */

.markdown-body {
  font-family: Palatino, "Palatino Linotype", "Palatino LT STD", "Book Antiqua", Georgia, serif;
  color: #111;
  counter-reset: sidenote-counter;
}

h1, h2, h3, h4, h5 {
  margin-top: 1rem;
  line-height: 1;
  font-weight: 400;
}

h1 {
  margin-bottom: 1.5rem;
  font-size: 3.0rem;
}

h2 {
  margin-bottom: 0;
  font-size: 2.5rem;
}

h3 {
  font-size: 2.0rem;
  margin-bottom: 0;
}

h4 {
  font-size: 1.5rem;
  margin-bottom: 0;
}

hr {
  display: block;
  height: 1px;
  width: 55%;
  border: 0;
  border-top: 1px solid #ccc;
  margin: 1em 0;
  padding: 0;
}

p.subtitle {
  font-style: italic;
  margin-top: 1rem;
  margin-bottom: 1rem;
  font-size: 1.8rem;
  display: block;
  line-height: 1;
}

.numeral {
  font-family: et-book-roman-old-style;
}

.danger {
  color: red;
}

article {
  position: relative;
  padding: 5rem 0rem;
}

section {
  padding-top: 1rem;
  padding-bottom: 1rem;
}

p, ol, ul {
  font-size: 1.4rem;
  line-height: 2rem;
}

p {
  margin-top: 1.4rem;
  margin-bottom: 1.4rem;
  padding-right: 0;
  vertical-align: baseline;
}
/* Chapter Epigraphs */
div.epigraph {
  margin: 5em 0;
}

div.epigraph > blockquote {
  margin-top: 3em;
  margin-bottom: 3em;
}

div.epigraph > blockquote, div.epigraph > blockquote > p {
  font-style: italic;
}

div.epigraph > blockquote > footer {
  font-style: normal;
}

div.epigraph > blockquote > footer > cite {
  font-style: italic;
}
/* end chapter epigraphs styles */
blockquote {
  margin: 0;
  margin-right: 1em;
  font-size: 1.4rem;
  font-style: italic;
  padding-left: 1em;
  border-left: 0.25em solid #111;
}

blockquote p {
  margin-right: 40px;
}

blockquote footer {
  font-size: 1.1rem;
  text-align: right;
}

section > p, section > footer, section > table {
  width: 55%;
}

/* 50 + 5 == 55, to be the same width as paragraph */
section > ol, section > ul {
  width: 50%;
  -webkit-padding-start: 5%;
}

li:not(:first-child) {
  margin-top: 0.25rem;
}

figure {
  padding: 0;
  border: 0;
  font-size: 100%;
  font: inherit;
  vertical-align: baseline;
  max-width: 55%;
  -webkit-margin-start: 0;
  -webkit-margin-end: 0;
  margin: 0 0 3em 0;
}

figcaption {
  float: right;
  clear: right;
  margin-top: 0;
  margin-bottom: 0;
  font-size: 1.1rem;
  line-height: 1.6;
  vertical-align: baseline;
  position: relative;
  max-width: 40%;
}

figure.fullwidth figcaption {
  margin-right: 24%;
}

/* Links: replicate underline that clears descenders */
a:link, a:visited {
  color: inherit;
}

a:link {
  text-decoration: underline;
}

/* Sidenotes, margin notes, figures, captions */
img {
  max-width: 100%;
}

.sidenote, .marginnote {
  float: right;
  clear: right;
  margin-right: -60%;
  width: 50%;
  margin-top: 0;
  margin-bottom: 0;
  font-size: 1.1rem;
  line-height: 1.3;
  vertical-align: baseline;
  position: relative;
}

.sidenote-number {
  counter-increment: sidenote-counter;
}

.sidenote-number:after, .sidenote:before {
  font-family: et-book-roman-old-style;
  position: relative;
  vertical-align: baseline;
}

.sidenote-number:after {
  content: counter(sidenote-counter);
  font-size: 1rem;
  top: -0.5rem;
  left: 0.1rem;
}

.sidenote:before {
  content: counter(sidenote-counter) " ";
  top: -0.5rem;
}

blockquote .sidenote, blockquote .marginnote {
  margin-right: -82%;
  min-width: 59%;
  text-align: left;
}

div.fullwidth, table.fullwidth {
  width: 100%;
}

div.table-wrapper {
  overflow-x: auto;
  font-family: "Trebuchet MS", "Gill Sans", "Gill Sans MT", sans-serif;
}

.sans {
  font-family: "Gill Sans", "Gill Sans MT", Calibri, sans-serif;
  letter-spacing: .03em;
}

code {
  font-family: Consolas, "Liberation Mono", Menlo, Courier, monospace;
  font-size: 1.0rem;
  line-height: 1.42;
}

.sans > code {
  font-size: 1.2rem;
}

h1 > code, h2 > code, h3 > code {
  font-size: 0.80em;
}

.marginnote > code, .sidenote > code {
  font-size: 1rem;
}

pre.code {
  font-size: 0.9em;
  width: 52.5%;
  margin-left: 2.5%;
  overflow-x: auto;
}

pre.code.fullwidth {
  width: 90%;
}

.fullwidth {
  max-width: 90%;
  clear: both;
}

span.newthought {
  font-variant: small-caps;
  font-size: 1.2em;
}

input.margin-toggle {
  display: none;
}

label.sidenote-number {
  display: inline;
}

label.margin-toggle:not(.sidenote-number) {
  display: none;
}

.iframe-wrapper {
  position: relative;
  padding-bottom: 56.25%;
    /* 16:9 */
  padding-top: 25px;
  height: 0;
}

.iframe-wrapper iframe {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

table {
  font-size: 1.4rem;
  padding: 0;
}

table tr {
  border-top: 1px solid #cccccc;
  background-color: white;
  margin: 0;
  padding: 0;
}

table tr:nth-child(2n) {
  background-color: #f8f8f8;
}

table tr th {
  background-color: #f8f8f8;
  font-weight: bold;
  border: 1px solid #cccccc;
  text-align: left;
  margin: 0;
  padding: 0.5rem 1rem;
}

table tr td {
  border: 1px solid #cccccc;
  text-align: left;
  margin: 0;
  padding: 0.5rem 1rem;
}

table tr th :first-child, table tr td :first-child {
  margin-top: 0;
}

table tr th :last-child, table tr td :last-child {
  margin-bottom: 0;
}
`); err != nil {
		log.Fatalf("failed to parse: %s", err)
	}
}
