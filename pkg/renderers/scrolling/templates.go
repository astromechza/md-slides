package scrolling

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
	  body {
          display: flex;
          flex-flow: column;
          justify-content: start;
		  align-items: center;
      }
      .page-wrap {
          width: {{ $.PageXResPX }}px;
          height: {{ $.AdjustedPageYResPX }}px;
      }
      .slide-wrap {
          position: relative;
      }
      @page {
	      size: {{ $.PageXResPX }}px {{ $.PageYResPX }}px;
      }
    </style>
	{{ range .PreparedSlides }}
    <div class="page-wrap">
    <div class="slide-wrap slide-wrap-halign-{{ .Settings.HAlign }} slide-wrap-valign-{{ .Settings.VAlign }} slide-wrap-talign-{{ .Settings.TAlign }}" style="width: {{ .Settings.XResPX }}px; left: {{ .PageLeft }}px; top: {{ .PageTop }}px; height: {{ .Settings.YResPX }}px">
      <div class="markdown-body">
        {{ .Content }}
      </div>
      {{ with .Settings.FooterText }}<div class="page-footer">{{ . }}</div>{{ end }}
      <div class="page-number">{{ .PageNum }}/{{ $.PageCount }}</div>
    </div>
    </div>
    {{ end }}
  </body>
</html>
`
