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

# TODO

- [x] customise html renderer to add todo list functionality
- [x] a bunch of css fixes and styling
- [x] work out some css style for aspect ratio controlled slides
- [ ] **feature**: find a nice way to control aspect ratio/zoom px values via golang
- [ ] **feature**: serve static images/videos from relative paths
- [ ] **security**: remember to secure against absolute paths or ../ references
- [ ] **feature**: `{embed-command: ./my-command --help}`
- [ ] **feature**: allow slides to declare a zoom value
