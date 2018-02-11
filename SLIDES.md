<meta valign="bottom" halign="left">

# `md-slides`

> Ben Meier 2018
> \- Cape Town

---

### 1

---

### 2

A small demo of some flag things.

```
subcommand := flag.Arg(0)
switch subcommand {
case "serve":
    return Serve(flag.Args()[1:])
default:
    return fmt.Errorf("unknown subcommand '%s'", subcommand)
}
```

- A
  - B
    - C

> Study hard what interests you the most in the most undisciplined, irreverent and original manner possible.
> <footer>- Richard Feynman</footer>

---

# HEADINGS
## HEADINGS
### HEADINGS
#### HEADINGS
##### HEADINGS

---

Better directory structure

- `/_keynotes` could in future render keynotes in time with the slide transitions (so it can keep them in sync :))
- `/_perf` might drop http perf stats (not that its really necessary :P)
- `/_slides?page=0` defines the slide
- `/<filename>` defines a relative image to load and serve

---

blah

![testimage](http://via.placeholder.com/450x350)

blah

![testimage](./testimage.png)

---

# TODO

- [ ] **feature**: pdf rendering
- [ ] **security**: disable file serving if necessary or by default
- [ ] **security**: blacklist some file patterns from serving
- [ ] **feature**: page numbers
- [ ] **feature**: file path can be a url :D
- [x] **fix**: change slides urls to `/_slides/<int>`
- [x] **feature**: other keys for advancing slides (space, enter)

---

# A
# A
# A
# A
# A
# A

---

## Example of command output

```
{embedcommand: ["bash", "-c", "./md-slides --help || true"]}
```

---

## Rendering to html / pdf

- `wkhtmltopdf` exists but doesn't support css flexbox / grid
- Headless Google Chrome CLI can take very nice `--screenshot` but the `--print-to-pdf` rendering is a bit lacking
- Programatic solution using the Chrome Debug Port https://github.com/mafredri/cdp
- Follow what is done in https://github.com/Szpadel/chrome-headless-render-pdf/blob/master/index.js (NodeJS)

Many advantages:

- First-class page layout and feature support
- Can be embedded in a docker container if necessary
- Cross-platform
- Dependency is small enough to fit in `md-slides` binary
