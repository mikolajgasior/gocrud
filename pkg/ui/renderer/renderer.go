package renderer

import "html/template"

// NewHTMLRenderer creates a new HTMLRenderer with default templates
func NewHTMLRenderer() (*HTMLRenderer, error) {
	renderer := &HTMLRenderer{
		templates: make(map[string]*template.Template),
	}

	// Register default Input template (compact, no extra whitespace)
	inputTmplStr := `{{- if eq .InputType "textarea" 
-}}<textarea{{ 
if .Name }} name="{{ .Name }}"{{ end }}{{ 
if .ID }} id="{{ .ID }}"{{ end }}{{ 
if .Required }} required{{ end }}{{ 
if .Pattern }} pattern="{{ .Pattern }}"{{ end }}{{ 
if .MinLength }} minlength="{{ .MinLength }}"{{ end }}{{ 
if .MaxLength }} maxlength="{{ .MaxLength }}"{{ end }}>{{ .Value }}</textarea>{{- 
else if eq .InputType "checkbox" -}}<input type="checkbox"{{ 
if .Name }} name="{{ .Name }}"{{ end }}{{ 
if .ID }} id="{{ .ID }}"{{ end }}{{ 
if .Checked }} checked{{ end }}>{{- 
else -}}<input type="{{ .InputType }}"{{ 
if .Name }} name="{{ .Name }}"{{ end }}{{ 
if .ID }} id="{{ .ID }}"{{ end }}{{ 
if .Required }} required{{ end }}{{ 
if .Pattern }} pattern="{{ .Pattern }}"{{ end }}{{ 
if .MinLength }} minlength="{{ .MinLength }}"{{ end }}{{ 
if .MaxLength }} maxlength="{{ .MaxLength }}"{{ end }}{{ 
if .Min }} min="{{ .Min }}"{{ end }}{{ 
if .Max }} max="{{ .Max }}"{{ end }}{{ 
if .Value }} value="{{ .Value }}"{{ end }}{{ 
if .Placeholder }} placeholder="{{ .Placeholder }}"{{ end }}{{ 
if .Disabled }} disabled{{ end }}{{ 
if .ReadOnly }} readonly{{ end }}/>{{- end -}}`

	inputTmpl, err := template.New("input").Parse(inputTmplStr)
	if err != nil {
		return nil, err
	}
	renderer.templates["Input"] = inputTmpl
	renderer.template = inputTmpl

	return renderer, nil
}
