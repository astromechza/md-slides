package sliderenderer

import (
	"fmt"
	"html/template"
)

func LoadTemplates() (*template.Template, error) {
	root := template.New("root")

	root = root.Funcs(template.FuncMap{
		"add": func(i, j int) int {
			return i + j
		},
	})

	if _, err := root.New("scripts").Parse(scriptHeader); err != nil {
		return nil, fmt.Errorf("failed to load scripts: %s", err)
	}
	if _, err := root.New("normalize-css").Parse(normalizeCSS); err != nil {
		return nil, fmt.Errorf("failed to load normalize-css: %s", err)
	}
	if _, err := root.New("markdown-css").Parse(markdownCSS); err != nil {
		return nil, fmt.Errorf("failed to load markdown-css: %s", err)
	}
	if _, err := root.New("other-css").Parse(styleHeader); err != nil {
		return nil, fmt.Errorf("failed to load other-css: %s", err)
	}
	if _, err := root.New("multipage-css").Parse(styleMultiHeader); err != nil {
		return nil, fmt.Errorf("failed to load multipage-css: %s", err)
	}

	if _, err := root.New("slide-prefix").Parse(`
	<div class="{{ .SlideClasses }}" style="width: {{ .XRes }}px; height: {{ .YRes }}px; background: {{ .BGCSS }}">
		<div class="markdown-body">
	`); err != nil {
		return nil, fmt.Errorf("failed to load slide-prefix: %s", err)
	}

	if _, err := root.New("slide-suffix").Parse(`
			</div>
			<div class="page-number">{{ .PageNum }}/{{ len .CachedSlides }}</div>
		</div>
		`); err != nil {
		return nil, fmt.Errorf("failed to load slide-suffix: %s", err)
	}

	return root, nil
}