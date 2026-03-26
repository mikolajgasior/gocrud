package form

import (
	"bytes"
	"html"
	"html/template"

	uiinput "codeberg.org/mikolajgasior/gocrud/pkg/ui/input"
)

type Form struct {
	Path        string
	InputsOrder []string
	Inputs      map[string]uiinput.Input
}

func (f Form) HTML() template.HTML {
	var buf bytes.Buffer

	buf.WriteString(`<form class="login" style="margin: auto; width: var(--login-form-width);" method="POST" action="`)
	buf.WriteString(html.EscapeString(f.Path))
	buf.WriteString(`">`)

	for _, name := range f.InputsOrder {
		if inp, ok := f.Inputs[name]; ok {
			buf.WriteString(`<p class="input">`)
			buf.WriteString(`<label>`)
			buf.WriteString(html.EscapeString(name))
			buf.WriteString(`:</label>`)
			buf.WriteString(string(inp.HTML()))
			buf.WriteString(`</p>`)
		}
	}

	buf.WriteString(`<p class="buttons">`)
	buf.WriteString(`<button type="submit">Create</button>`)
	buf.WriteString(`</p>`)
	buf.WriteString(`</form>`)

	return template.HTML(buf.String())
}
