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
- [ ] work out some css style for aspect ratio controlled slides
- [ ] serve static images/videos from relative paths
- [ ] remember to secure against absolute paths or ../ references
