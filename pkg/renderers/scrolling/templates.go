package scrolling

import (
	"html/template"
	"log"
)

func AddStyleOverridesTemplate(root *template.Template) {
	if _, err := root.New("style.overrides").Parse(`
html {
	height: auto;
}

body {
	height: 100%;
    display: flex;
    flex-flow: column;
	justify-content: start;
	position: relative;
}

.slide-wrap {
	position: relative;
	margin-top: 1.5em;
	margin-bottom: 1.5em;
}
`); err != nil {
		log.Fatalf("failed to parse: %s", err)
	}
}

const slideTemplate = `
<!DOCTYPE html>
  <head>
    <title>{{ .Title }}</title>
	<meta charset="utf-8">
  </head>
  <body>
    <style>
      {{ template "style.normalize" . }}
      {{ template "style.markdown" . }}
      {{ template "style.chroma" .}}
      {{ template "style.common" . }}
      {{ template "style.overrides" .}}
    </style>
	{{ range .PreparedSlides }}
    <div class="slide-wrap slide-wrap-halign-{{ .Settings.HAlign }} slide-wrap-valign-{{ .Settings.VAlign }} slide-wrap-talign-{{ .Settings.TAlign }}" style="width: {{ .Settings.XResPX }}px; height: {{ .Settings.YResPX }}px">
      <div class="markdown-body">
        {{ .Content }}
      </div>
      {{ with .Settings.FooterText }}<div class="page-footer">{{ . }}</div>{{ end }}
      <div class="page-number">{{ .PageNum }}/{{ $.PageCount }}</div>
    </div>
    {{ end }}
  </body>
</html>
`
