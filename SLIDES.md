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

- `/_slides?N` serves the slide html
- `/<filename>` defines a relative image to load and serve

---

blah

![testimage](http://via.placeholder.com/450x350)

blah

![testimage](./testimage.png)

[../../../../../etc/hosts](../../../../../etc/hosts)

---

# TODO

- [ ] **fix** multislide doesn't work with custom resolution
- [ ] **feature** code hightlighting
- [ ] **fix** major improvements using templating

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
