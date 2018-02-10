package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

window.onresize = function(event) {
	var el = document.getElementById("body-inner");
	var m = 50;
	var wi = el.clientWidth + m;
	var hi = el.clientHeight + m;
	var ws = window.innerWidth / wi;
	var hs = window.innerHeight / hi;
	var ss = Math.min(ws, hs);
	el.style.transform = "scale(" + ss + ")";
};

document.addEventListener("DOMContentLoaded", function(event) {
    window.onresize(null);
});

</script>
`

const styleHeader = `
<style>
html {
	height: 100%;
	font-size: 20px;
}

body {
	height: 100%;
    display: flex;
    flex-flow: column;
	background-color: grey;
	justify-content: center;
}

#body-inner {
	display: flex;
	flex-flow: column;
	align-self: center;
	box-sizing: border-box;
	background-color: #fffff8;
	padding: 1rem;
    border-radius: 0.1rem;
	box-shadow: 0px 0.2rem 0.6rem black;
	padding-left: 3rem;
    padding-right: 3rem;
	position: absolute;
	overflow: hidden;
}

#body-inner.centered {
	justify-content: center;
}
</style>
`

func parseResString(i string) (int, int, error) {
	i = strings.TrimSpace(strings.ToLower(i))
	parts := strings.Split(i, "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("res string '%s' did not contain one 'x'", i)
	}
	xres, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse x value of res string '%s': %s", i, err)
	}
	yres, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse y value of res string '%s': %s", i, err)
	}
	if xres <= 0 {
		return 0, 0, fmt.Errorf("x value of rest string '%s' should be > 0", i)
	}
	if yres <= 0 {
		return 0, 0, fmt.Errorf("y value of rest string '%s' should be > 0", i)
	}
	return int(xres), int(yres), nil
}

func Serve(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	hotFlag := fs.Bool("hot", false, "reload, reparse, and regenerate slides on each refresh")
	resFlag := fs.String("res", "1600x900", "set render aspect ratio and zoom for rendering")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("expected a single source file as argument")
	}
	filename := fs.Arg(0)

	xres, yres, err := parseResString(*resFlag)
	if err != nil {
		return fmt.Errorf("bad res string: %s", err)
	}

	sr := SlideRenderer{Filename: filename, Hot: *hotFlag, XRes: xres, YRes: yres}
	http.HandleFunc("/slides/", func(rw http.ResponseWriter, req *http.Request) {
		snRaw := req.URL.Path[len("/slides/"):]
		if snRaw == "" {
			rw.Header().Set("location", "/slides/0")
			rw.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		sn, err := strconv.ParseInt(snRaw, 10, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(http.StatusText(http.StatusBadRequest)))
			return
		}
		sr.Serve(int(sn), rw, req)
	})

	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		return err
	}

	return nil
}
