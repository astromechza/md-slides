<meta valign="center" halign="center" talign="center">
<meta footer="github.com/AstromechZA/md-slides">

# md-slides

AstromechZA 2018

---

<meta footer="">

### Intro

`md-slides` is a tool for presenting html-based presentations based on a Markdown documents.

Online Example (scrolling): [astromechza.github.io/md-slides/](https://astromechza.github.io/md-slides/)

Project URL: [github.com/AstromechZA/md-slides](https://github.com/AstromechZA/md-slides)

Features:

- single static binary with no runtime dependencies
- no external javascript or css
- single & multipage modes
- supports embedding images
- uses `blackfriday` markdown library
- supports custom aspect-ratios and zooms
- prints to pdf well

---

### Prior art

- _many_. Just Google "Markdown slides"
	- reveal.js
	- remark.js
	- GitPitch
	- etc..

- But I believe in building software to meet your own needs, and reinventing the wheel for yourself can be some good fun.

- And this provides a good balance between self-hosted, lock-in-free software, and feature-packed service offerings. *While still making the slide source easily readable in its source or rendered markdown forms*.

---

### Installation

Although you _can_ build and install it from source, we recommend that you pull
the appropriate binary for your system from the project releases page [here](https://github.com/AstromechZA/md-slides/releases).

Or use the installation script to be quick:

```
$ curl https://raw.githubusercontent.com/AstromechZA/md-slides/master/install.sh | OVERRIDE_RELEASE_TAG=~/bin sh
```


```bash
$ wget https://github.com/AstromechZA/md-slides/releases/download/v0.0/md-slides.linux.amd64

# verify integrity
$ wget https://github.com/AstromechZA/md-slides/releases/download/v0.0/SHA256SUMS
$ sha256sum -c SHA256SUMS

$ chmod +x md-slides.linux.amd64
$ mv md-slides.linux.amd64 ~/bin/md-slides
```

#### Dependencies

The following 3rd party libraries (and their dependencies) are statically linked into the binary:

- `github.com/russross/blackfriday` : markdown processing
- `github.com/gorilla/mux` : a better http server router
- `github.com/alecthomas/chroma` : syntax highlighting

---

### subcommands

```
{embedcommand: ["bash", "-c", "./md-slides --help || true"]}
```

### `version` subcommand

```
{embedcommand: ["./md-slides", "version"]}
```

---

### `serve` subcommand

```
{embedcommand: ["bash", "-c", "./md-slides serve --help || true"]}
```

---

### Input format

A single markdown document represents the multiple slides separated by a horizantal rule:

```markdown
# this is slide one

stuff

---

# this is slide two

stuff

---

# this is slide three

stuff
```

---

### Paged vs Scrolling mode

The `-mode` flag can switch between display modes.

- `paged` is the default javascript-powered presentation mode with dynamic scaling. It attempts to scale to fit the presentation window while maintaining the slide aspect ratio and font size. It responds to key controls to advance slides.


- `scrolling` is a **no-javascript** mode with each slide below the other and with no dynamic scaling. This mode prints very well and has page hints for **PDF rendering and export**. This is also the mode used by the `--export-to` flag. It can be useful for distributing slides via email or embedding in other sites via an iframe.

---

### Extension: `{embedcommand: ...}`

Embed command embeds the combined stdout/stdin output of executing a given subcommand. The subcommand should be embeded as a json array like `["ls", "-al", "blah"]`.

For example a slide containing `{embed``command: ["date"]}` could output:

```
{embedcommand: ["date"]}
```

This makes it easy to present live content or to embed dynamic data into the slides based on the output of scripts.

---

### Extension: image dimensions in url fragment

To improve the sizing of images embedded into your slides, you can specify the width and height as url fragment parameters.

```markdown
![alt text](/my/relative/path#height=500px)

![alt text](/my/relative/path#width=500px)
```

---

### Metadata: `halign` and `valign`

These meta tags control where in the slide the content is positioned. They are most effective when used to position content to the bottom or corners of a slide.

```html
<meta valign="bottom" halign="left">
```

Allowed values for `valign` are `top`, `center`, `bottom`

Allowed values for `halign` are `left`, `center`, `right`

Good for:

- title slides
- thank you slides
- slides containing centered images, questions, blocks, etc.

---

<meta valign="bottom" halign="right">

### `Example of <meta valign="bottom" halign="right">`

second line

---

### Metadata: `talign`

The `talign` metadata should be used with `halign` and `valign` to direct the text-alignment of the slide.

```html
<meta talign="left">
```

Allowed values for `talign` are `left`, `center`, `right`

---

<meta talign="right">

### `Example of <meta talign="right">`

second line

---

<meta footer="January 1970 &vert; Some Conference">

### Metadata: `footer`

Add some footer text to the bottom left of each slide using the `footer` metadata key.

The value will persist between pages until overriden by a new value or empty string.

```html
<meta footer="January 1970 &vert; Some Conference">
```

---

<meta footer="">

### Markdown support

Anything that `https://github.com/russross/blackfriday` v2 supports.

- *italic*, **bold**, ~~strike~~
- [ ] todo
- [x] todo done
- [links](http://google.com)
- [footnote links][1]
- images
- `inline code`
- code blocks
    + sublists

> block quotes

[1]: http://google.com

---

### Code highlighting

```go
package main

import "fmt"

func main() {
	for i := 1; i <= 100; i++ {
		result := ""
		if i%3 == 0 { result += "Fizz" }
		if i%5 == 0 { result += "Buzz" }
		if result != "" {
			fmt.Println(result)
			continue
		}
		fmt.Println(i)
	}
}
```

---

### Tables?

<meta valign="center" halign="center">

| Tables   |      Are      |  Cool |
|----------|:-------------:|------:|
| col 1 is |  left-aligned | $1600 |
| col 2 is |    centered   |   $12 |
| col 3 is | right-aligned |    $1 |

---

<meta valign="center" halign="center" talign="center">

### Text above image

![A test image](windmill.jpeg#height=500px)

Text below image

---

### TODO

Development is still in progress:

- [x] v0.0 available on github
- [ ] Rework CSS styles
- [ ] Testing!
- [ ] v1.0?

---

<meta valign="center" halign="center" talign="center">

## Thanks!
