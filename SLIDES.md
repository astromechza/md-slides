# md-slides

### Ben Meier 2018
### Oracle Cape Town

---

### 1

---

### 2

```
subcommand := flag.Arg(0)
switch subcommand {
case "serve":
    return Serve(flag.Args()[1:])
default:
    return fmt.Errorf("unknown subcommand '%s'", subcommand)
}
```

---

# HEADINGS
## HEADINGS
### HEADINGS
#### HEADINGS
##### HEADINGS

---

blah

![testimage](http://via.placeholder.com/450x350)

blah

---

# TODO

- [x] customise html renderer to add todo list functionality
- [x] a bunch of css fixes and styling
- [x] work out some css style for aspect ratio controlled slides
- [ ] **feature**: find a nice way to control aspect ratio/zoom px values via golang
- [ ] **feature**: allow slides to declare a zoom value
- [ ] **feature**: allow slides to declare a vertical / horizantal alignment
- [ ] **feature**: serve static images/videos from relative paths
- [ ] **security**: remember to secure against absolute paths or ../ references (may reconsider this)
- [ ] **feature**: `{embed-command: ./my-command --help}`
- [x] **fix**: background color of markdown block does not match background colour of slide
- [ ] **feature**: live reload? can we detect changes to the markdown file and reload it? hot swap it as often as possible after all generation and embedding has taken place?
- [ ] **spike**: explore rendering to html/pdf
- [ ] **spike**: explore theming support (some kind of optional named css dropin)
- [ ] **fix**: todo checkboxes are too small!
- [x] **fix**: overflow on body-inner

---

# A
# A
# A
# A
# A
# A
# A
# A
# A
# A
# A
# A
# A
# A
# A
# A
