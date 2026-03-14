package input

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
)

const (
	DefaultTagName = "html"
)

const (
	TypeText     = "text"
	TypeTextarea = "textarea"
	TypePassword = "password"
	TypeEmail    = "email"
	TypeNumber   = "number"
	TypeCheckbox = "checkbox"
)

type Input struct {
	FieldName   string
	InputType   string
	Value       string
	Required    bool
	Pattern     string
	MinLength   int
	MaxLength   int
	Min         int
	Max         int
	ID          string
	Name        string
	Checked     bool
	Placeholder string
	Disabled    bool
	ReadOnly    bool
}

func (i Input) HTML() template.HTML {
	var buf bytes.Buffer

	if i.InputType == TypeTextarea {
		buf.WriteString("<textarea")
		if i.ID != "" {
			buf.WriteString(fmt.Sprintf(` id="%s"`, html.EscapeString(i.ID)))
		}
		if i.Name != "" {
			buf.WriteString(fmt.Sprintf(` name="%s"`, html.EscapeString(i.Name)))
		}
		if i.Required {
			buf.WriteString(" required")
		}
		if i.Pattern != "" {
			buf.WriteString(fmt.Sprintf(` pattern="%s"`, html.EscapeString(i.Pattern)))
		}
		if i.MinLength > 0 {
			buf.WriteString(fmt.Sprintf(` minlength="%d"`, i.MinLength))
		}
		if i.MaxLength > 0 {
			buf.WriteString(fmt.Sprintf(` maxlength="%d"`, i.MaxLength))
		}
		buf.WriteString(">")
		buf.WriteString(html.EscapeString(i.Value))
		buf.WriteString("</textarea>")
	} else if i.InputType == TypeCheckbox {
		buf.WriteString("<input type=\"checkbox\"")
		if i.ID != "" {
			buf.WriteString(fmt.Sprintf(` id="%s"`, html.EscapeString(i.ID)))
		}
		if i.Name != "" {
			buf.WriteString(fmt.Sprintf(` name="%s"`, html.EscapeString(i.Name)))
		}
		if i.Checked {
			buf.WriteString(" checked")
		}

		buf.WriteString(">")
	} else {
		buf.WriteString(fmt.Sprintf("<input type=\"%s\"", html.EscapeString(i.InputType)))

		if i.ID != "" {
			buf.WriteString(fmt.Sprintf(` id="%s"`, html.EscapeString(i.ID)))
		}
		if i.Name != "" {
			buf.WriteString(fmt.Sprintf(` name="%s"`, html.EscapeString(i.Name)))
		}

		if i.Value != "" && i.InputType != TypePassword {
			buf.WriteString(fmt.Sprintf(` value="%s"`, html.EscapeString(i.Value)))
		}

		if i.Required {
			buf.WriteString(" required")
		}
		if i.Pattern != "" {
			buf.WriteString(fmt.Sprintf(` pattern="%s"`, html.EscapeString(i.Pattern)))
		}
		if i.MinLength > 0 {
			buf.WriteString(fmt.Sprintf(` minlength="%d"`, i.MinLength))
		}
		if i.MaxLength > 0 {
			buf.WriteString(fmt.Sprintf(` maxlength="%d"`, i.MaxLength))
		}
		if i.Min != 0 {
			buf.WriteString(fmt.Sprintf(` min="%d"`, i.Min))
		}
		if i.Max != 0 {
			buf.WriteString(fmt.Sprintf(` max="%d"`, i.Max))
		}
		if i.Placeholder != "" {
			buf.WriteString(fmt.Sprintf(` placeholder="%s"`, html.EscapeString(i.Placeholder)))
		}
		if i.Disabled {
			buf.WriteString(" disabled")
		}
		if i.ReadOnly {
			buf.WriteString(" readonly")
		}
		buf.WriteString("/>")
	}

	return template.HTML(buf.String())
}
