package form

import (
	"strings"
	"testing"

	uiinput "miko.gs/struct-crud/pkg/ui/input"
)

func TestFormRendering(t *testing.T) {
	inputs := []uiinput.Input{
		{FieldName: "Username", InputType: "text", Name: "username"},
		{FieldName: "Email", InputType: "email", Name: "email"},
		{FieldName: "Age", InputType: "number", Name: "age"},
	}

	inputMap := make(map[string]uiinput.Input)
	for _, inp := range inputs {
		inputMap[inp.FieldName] = inp
	}

	form := &Form{
		Path:        "/create-user",
		InputsOrder: []string{"Username", "Email", "Age"},
		Inputs:      inputMap,
	}

	html := form.HTML()

	if !strings.Contains(string(html), "<form") {
		t.Fatal("Form rendering failed: no <form> tag found")
	}

	if !strings.Contains(string(html), `action="/create-user"`) {
		t.Fatal("Form rendering failed: incorrect action path")
	}

	inputCount := strings.Count(string(html), "<input")
	if inputCount < 2 {
		t.Fatalf("Form rendering failed: expected at least 2 inputs, got %d", inputCount)
	}

	usernamePos := strings.Index(string(html), "Username")
	emailPos := strings.Index(string(html), "Email")
	agePos := strings.Index(string(html), "Age")

	if usernamePos == -1 || emailPos == -1 || agePos == -1 {
		t.Fatal("Form rendering failed: not all inputs found in output")
	}

	if !(usernamePos < emailPos && emailPos < agePos) {
		t.Fatalf("Form rendering failed: inputs not in correct order\nUsername pos: %d, Email pos: %d, Age pos: %d",
			usernamePos, emailPos, agePos)
	}

	if !strings.Contains(string(html), "type=\"submit\"") {
		t.Fatal("Form rendering failed: no submit button found")
	}
}

func TestFormRenderingEmptyInputs(t *testing.T) {
	form := &Form{
		Path:        "/empty-form",
		InputsOrder: []string{},
		Inputs:      make(map[string]uiinput.Input),
	}

	html := form.HTML()

	if !strings.Contains(string(html), "<form") {
		t.Fatal("Empty form rendering failed: no <form> tag found")
	}

	if !strings.Contains(string(html), `action="/empty-form"`) {
		t.Fatal("Empty form rendering failed: incorrect action path")
	}

	if !strings.Contains(string(html), "type=\"submit\"") {
		t.Fatal("Empty form rendering failed: no submit button found")
	}
}

func TestFormRenderingWithSingleInput(t *testing.T) {
	inputs := []uiinput.Input{
		{FieldName: "OnlyField", InputType: "text", Name: "only_field"},
	}

	inputMap := make(map[string]uiinput.Input)
	for _, inp := range inputs {
		inputMap[inp.FieldName] = inp
	}

	form := &Form{
		Path:        "/single-input",
		InputsOrder: []string{"OnlyField"},
		Inputs:      inputMap,
	}

	html := form.HTML()

	if !strings.Contains(string(html), "<form") {
		t.Fatal("Single input form rendering failed: no <form> tag found")
	}

	if !strings.Contains(string(html), `action="/single-input"`) {
		t.Fatal("Single input form rendering failed: incorrect action path")
	}

	inputCount := strings.Count(string(html), "<input")
	if inputCount != 1 {
		t.Fatalf("Single input form rendering failed: expected 1 input, got %d", inputCount)
	}

	if !strings.Contains(string(html), "OnlyField") {
		t.Fatal("Single input form rendering failed: field name not found")
	}
}
