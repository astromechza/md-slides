package paged

import (
	"html/template"
	"log"
)

func AddScriptTemplate(root *template.Template) {
	if _, err := root.New("script.all").Parse(`
document.onkeydown = function(evt) {
	evt = evt || window.event
	if ([13, 32, 33, 39, 40].indexOf(evt.keyCode) >= 0) {
		window.location = "{{ .URLPath }}?page={{ .PageNext }}";
	}
	if ([8, 34, 37, 38].indexOf(evt.keyCode) >= 0 ) {
		window.location = "{{ .URLPath }}?page={{ .PagePrev }}";
	}
}

window.onresize = function(event) {
	var el = document.getElementsByClassName("slide-wrap")[0];
	var m = 40;
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
`); err != nil {
		log.Fatalf("failed to parse: %s", err)
	}
}

func AddStyleOverridesTemplate(root *template.Template) {
	if _, err := root.New("style.overrides").Parse(`
html {
	height: 100%;
}

body {
	height: 100%;
    display: flex;
    flex-flow: column;
	justify-content: center;
}
`); err != nil {
		log.Fatalf("failed to parse: %s", err)
	}
}

const slideTemplate = `
<!DOCTYPE html>
  <head>
    <title>{{ .Title }} ({{ .PageNum }}/{{ .PageCount }})</title>
	<meta charset="utf-8">
  </head>
  <body>
    <script>
    {{ template "script.all" . }}
    </script>
    <style>
      {{ template "style.normalize" . }}
      {{ template "style.markdown" . }}
      {{ template "style.chroma" .}}
      {{ template "style.common" . }}
      {{ template "style.overrides" .}}
    </style>
    <div class="slide-wrap slide-wrap-halign-{{ .Settings.HAlign }} slide-wrap-valign-{{ .Settings.VAlign }} slide-wrap-talign-{{ .Settings.TAlign }}" style="width: {{ .Settings.XResPX }}px; height: {{ .Settings.YResPX }}px; color: {{ .Settings.FontColor }}">
      <div class="markdown-body">
        {{ .SlideContent }}
      </div>
      {{ with .Settings.FooterText }}<div class="page-footer">{{ . }}</div>{{ end }}
      <div class="page-number">{{ .PageNum }}/{{ .PageCount }}</div>
    </div>
  </body>
</html>
`
