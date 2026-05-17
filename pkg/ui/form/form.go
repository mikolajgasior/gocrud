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

	buf.WriteString(`<form action="#" method="POST" action="`)
	buf.WriteString(html.EscapeString(f.Path))
	buf.WriteString(`">`)

	for _, name := range f.InputsOrder {
		if inp, ok := f.Inputs[name]; ok {
			buf.WriteString(`<div class="form-group"><label class="form-label" for="`)
			buf.WriteString(html.EscapeString(inp.Name))
			buf.WriteString(`">`)
			buf.WriteString(html.EscapeString(name))
			buf.WriteString(`</label>`)
			buf.WriteString(string(inp.HTML()))
			buf.WriteString(`</div>`)
		}
	}

	buf.WriteString(`<div class="form-buttons"><button type="submit" class="btn-submit">Submit</button><button type="reset" class="btn-reset">Reset</button></div></form>`)

	return template.HTML(buf.String())
}
