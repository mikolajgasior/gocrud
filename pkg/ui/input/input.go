package input

const (
	DefaultTagName = "html"
)

// Input types
const (
	TypeText     = "text"
	TypeTextarea = "textarea"
	TypePassword = "password"
	TypeEmail    = "email"
	TypeNumber   = "number"
	TypeCheckbox = "checkbox"
)

// Input represents a single HTML input field with all its attributes as struct fields
type Input struct {
	FieldName   string
	InputType   string // text, password, email, number, checkbox, textarea
	Value       string
	Required    bool
	Pattern     string
	MinLength   int
	MaxLength   int
	Min         int
	Max         int
	ID          string
	Name        string
	Checked     bool // for checkboxes
	Placeholder string
	Disabled    bool
	ReadOnly    bool
	HasValue    bool // whether to include value attribute (false for passwords)
}
