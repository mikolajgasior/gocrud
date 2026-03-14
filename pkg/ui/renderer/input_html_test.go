package renderer

import (
	uibuilder "miko.gs/struct-crud/pkg/ui/builder"
	uiinput "miko.gs/struct-crud/pkg/ui/input"
	"testing"
)

type Test1 struct {
	FirstName     string `html:"req len:5,25"`
	LastName      string `html:"req len:2,50"`
	Age           int    `html:"req val:18,150"`
	Price         int    `html:"req val:0,9999"`
	PostCode      string `html:"req" html_regexp:"^[0-9][0-9]-[0-9][0-9][0-9]$"`
	Email         string `html:"req uiemail"`
	BelowZero     int    `html:"val:-6,-2"`
	DiscountPrice int    `html:"val:0,8000"`
	Country       string `html_regexp:"^[A-Z][A-Z]$"`
	County        string `html:"len:,40 uitextarea"`
	ExcludeMe     string `html:"req"`
}

func TestGenerateHTML(t *testing.T) {
	s := &Test1{}
	inputs, _ := uibuilder.BuildInputs(s, &uibuilder.Options{
		RestrictFields: map[string]struct{}{
			"FirstName": struct{}{},
			"Age":       struct{}{},
			"PostCode":  struct{}{},
			"Email":     struct{}{},
			"Country":   struct{}{},
			"County":    struct{}{},
			"ExcludeMe": struct{}{},
		},
		ExcludeFields: map[string]struct{}{
			"ExcludeMe": struct{}{},
		},
		IDPrefix:   "id_",
		NamePrefix: "name_",
	})

	renderer, err := NewHTMLRenderer()
	if err != nil {
		t.Fatal("Failed to create HTMLRenderer:", err)
	}

	// Build a map for easier lookup
	inputMap := make(map[string]uiinput.Input)
	for _, input := range inputs {
		inputMap[input.FieldName] = input
	}

	// Render each input and compare
	expected := `<input type="text" name="name_FirstName" id="id_FirstName" required minlength="5" maxlength="25"/>`
	got := renderer.Render(inputMap["FirstName"])
	if got != expected {
		t.Fatalf("GenerateHTML failed for 'FirstName'\nExpected: %s\nGot:      %s", expected, got)
	}

	expected = `<input type="number" name="name_Age" id="id_Age" required min="18" max="150"/>`
	got = renderer.Render(inputMap["Age"])
	if got != expected {
		t.Fatalf("GenerateHTML failed for 'Age'\nExpected: %s\nGot:      %s", expected, got)
	}

	expected = `<input type="text" name="name_PostCode" id="id_PostCode" required pattern="^[0-9][0-9]-[0-9][0-9][0-9]$"/>`
	got = renderer.Render(inputMap["PostCode"])
	if got != expected {
		t.Fatalf("GenerateHTML failed for 'PostCode'\nExpected: %s\nGot:      %s", expected, got)
	}

	expected = `<input type="email" name="name_Email" id="id_Email" required/>`
	got = renderer.Render(inputMap["Email"])
	if got != expected {
		t.Fatalf("GenerateHTML failed for 'Email'\nExpected: %s\nGot:      %s", expected, got)
	}

	expected = `<input type="text" name="name_Country" id="id_Country" pattern="^[A-Z][A-Z]$"/>`
	got = renderer.Render(inputMap["Country"])
	if got != expected {
		t.Fatalf("GenerateHTML failed for 'Country'\nExpected: %s\nGot:      %s", expected, got)
	}

	expected = `<textarea name="name_County" id="id_County" maxlength="40"></textarea>`
	got = renderer.Render(inputMap["County"])
	if got != expected {
		t.Fatalf("GenerateHTML failed for 'County'\nExpected: %s\nGot:      %s", expected, got)
	}
}

func TestGenerateHTMLWithValues(t *testing.T) {
	s := &Test1{}
	inputs, _ := uibuilder.BuildInputs(s, &uibuilder.Options{
		RestrictFields: map[string]struct{}{
			"FirstName": struct{}{},
			"Age":       struct{}{},
			"PostCode":  struct{}{},
			"Email":     struct{}{},
			"Country":   struct{}{},
			"County":    struct{}{},
			"ExcludeMe": struct{}{},
		},
		ExcludeFields: map[string]struct{}{
			"ExcludeMe": struct{}{},
		},
		IDPrefix:   "id_",
		NamePrefix: "name_",
		OverwriteValues: map[string]string{
			"FirstName": `Joe "Joe"`,
			"Age":       "40",
			"Email":     "email@example.com",
			"Country":   "XX",
		},
	})

	renderer, err := NewHTMLRenderer()
	if err != nil {
		t.Fatal("Failed to create HTMLRenderer:", err)
	}

	// Build a map for easier lookup
	inputMap := make(map[string]uiinput.Input)
	for _, input := range inputs {
		inputMap[input.FieldName] = input
	}

	expected := `<input type="text" name="name_FirstName" id="id_FirstName" required minlength="5" maxlength="25" value="Joe &#34;Joe&#34;"/>`
	got := renderer.Render(inputMap["FirstName"])
	if got != expected {
		t.Fatalf("GenerateHTML failed for 'FirstName'\nExpected: %s\nGot:      %s", expected, got)
	}

	expected = `<input type="number" name="name_Age" id="id_Age" required min="18" max="150" value="40"/>`
	got = renderer.Render(inputMap["Age"])
	if got != expected {
		t.Fatalf("GenerateHTML failed for 'Age'\nExpected: %s\nGot:      %s", expected, got)
	}

	expected = `<input type="text" name="name_PostCode" id="id_PostCode" required pattern="^[0-9][0-9]-[0-9][0-9][0-9]$"/>`
	got = renderer.Render(inputMap["PostCode"])
	if got != expected {
		t.Fatalf("GenerateHTML failed for 'PostCode'\nExpected: %s\nGot:      %s", expected, got)
	}

	expected = `<input type="email" name="name_Email" id="id_Email" required value="email@example.com"/>`
	got = renderer.Render(inputMap["Email"])
	if got != expected {
		t.Fatalf("GenerateHTML failed for 'Email'\nExpected: %s\nGot:      %s", expected, got)
	}

	expected = `<input type="text" name="name_Country" id="id_Country" pattern="^[A-Z][A-Z]$" value="XX"/>`
	got = renderer.Render(inputMap["Country"])
	if got != expected {
		t.Fatalf("GenerateHTML failed for 'Country'\nExpected: %s\nGot:      %s", expected, got)
	}

	expected = `<textarea name="name_County" id="id_County" maxlength="40"></textarea>`
	got = renderer.Render(inputMap["County"])
	if got != expected {
		t.Fatalf("GenerateHTML failed for 'County'\nExpected: %s\nGot:      %s", expected, got)
	}
}

func TestGenerateHTMLWithFieldValues(t *testing.T) {
	s := &Test1{
		FirstName: "Joe",
		Age:       60,
		Email:     "joe@example.com",
	}
	inputs, fieldsOrder := uibuilder.BuildInputs(s, &uibuilder.Options{
		RestrictFields: map[string]struct{}{
			"FirstName": struct{}{},
			"Age":       struct{}{},
			"Email":     struct{}{},
			"ExcludeMe": struct{}{},
		},
		ExcludeFields: map[string]struct{}{
			"ExcludeMe": struct{}{},
		},
		IDPrefix:   "id_",
		NamePrefix: "name_",
		OverwriteValues: map[string]string{
			"FirstName": `Joe "Joe"`,
		},
		Values: true,
	})

	renderer, err := NewHTMLRenderer()
	if err != nil {
		t.Fatal("Failed to create HTMLRenderer:", err)
	}

	// Build a map for easier lookup
	inputMap := make(map[string]uiinput.Input)
	for _, input := range inputs {
		inputMap[input.FieldName] = input
	}

	expected := `<input type="text" name="name_FirstName" id="id_FirstName" required minlength="5" maxlength="25" value="Joe &#34;Joe&#34;"/>`
	got := renderer.Render(inputMap["FirstName"])
	if got != expected {
		t.Fatalf("GenerateHTML failed for 'FirstName'\nExpected: %s\nGot:      %s", expected, got)
	}

	expected = `<input type="number" name="name_Age" id="id_Age" required min="18" max="150" value="60"/>`
	got = renderer.Render(inputMap["Age"])
	if got != expected {
		t.Fatalf("GenerateHTML failed for 'Age'\nExpected: %s\nGot:      %s", expected, got)
	}

	expected = `<input type="email" name="name_Email" id="id_Email" required value="joe@example.com"/>`
	got = renderer.Render(inputMap["Email"])
	if got != expected {
		t.Fatalf("GenerateHTML failed for 'Email'\nExpected: %s\nGot:      %s", expected, got)
	}

	if fieldsOrder[0] != "FirstName" || fieldsOrder[1] != "Age" || fieldsOrder[2] != "Email" {
		t.Fatalf("GenerateHTML failed to output fields in the correct order\nExpected: [FirstName, Age, Email]\nGot:      %v", fieldsOrder)
	}
}

// Additional test for checkbox handling
func TestGenerateHTMLWithCheckbox(t *testing.T) {
	type TestWithCheckbox struct {
		Subscribe bool `html:""`
		Enabled   bool `html:""`
	}

	s := &TestWithCheckbox{
		Subscribe: true,
		Enabled:   false,
	}

	inputs, _ := uibuilder.BuildInputs(s, &uibuilder.Options{
		Values: true,
	})

	renderer, err := NewHTMLRenderer()
	if err != nil {
		t.Fatal("Failed to create HTMLRenderer:", err)
	}

	expected := `<input type="checkbox" name="Subscribe" checked>`
	got := renderer.Render(inputs[0])
	if got != expected {
		t.Fatalf("Checkbox with true value failed\nExpected: %s\nGot:      %s", expected, got)
	}

	expected = `<input type="checkbox" name="Enabled">`
	got = renderer.Render(inputs[1])
	if got != expected {
		t.Fatalf("Checkbox with false value failed\nExpected: %s\nGot:      %s", expected, got)
	}
}

// Test for password field (should not have value attribute)
func TestGenerateHTMLWithPassword(t *testing.T) {
	type TestWithPassword struct {
		Password string `html:"uipassword req"`
	}

	s := &TestWithPassword{
		Password: "secret123",
	}

	inputs, _ := uibuilder.BuildInputs(s, &uibuilder.Options{
		Values: true,
	})

	renderer, err := NewHTMLRenderer()
	if err != nil {
		t.Fatal("Failed to create HTMLRenderer:", err)
	}

	expected := `<input type="password" name="Password" required/>`
	got := renderer.Render(inputs[0])
	if got != expected {
		t.Fatalf("Password field should not have value attribute\nExpected: %s\nGot:      %s", expected, got)
	}
}

// Test for custom template registration
func TestCustomTemplateRegistration(t *testing.T) {
	renderer, err := NewHTMLRenderer()
	if err != nil {
		t.Fatal("Failed to create HTMLRenderer:", err)
	}

	// Register a custom template for a specific type
	type CustomWidget struct {
		Title string
		Value string
	}

	customTmpl := `<div class="widget"><h2>{{ .Title }}</h2><p>{{ .Value }}</p></div>`
	err = renderer.RegisterTemplate("CustomWidget", customTmpl)
	if err != nil {
		t.Fatal("Failed to register custom template:", err)
	}

	widget := CustomWidget{Title: "Test", Value: "Content"}
	expected := `<div class="widget"><h2>Test</h2><p>Content</p></div>`
	got := renderer.Render(widget)

	if got != expected {
		t.Fatalf("Custom template rendering failed\nExpected: %s\nGot:      %s", expected, got)
	}
}
