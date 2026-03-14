package input

import (
	"strings"
	"testing"
)

type TestStruct struct {
	Username string `html:"req len:3,20"`
	Email    string `html:"req uiemail"`
	Age      int    `html:"req val:18,120"`
	Password string `html:"uipassword req"`
	Active   bool   `html:""`
	Notes    string `html:"len:0,500 uitextarea"`
}

func TestBuildInputsAndRender(t *testing.T) {
	source := &TestStruct{
		Username: "johndoe",
		Email:    "john@example.com",
		Age:      30,
		Password: "secret123", // This should be ignored in output
		Active:   true,
		Notes:    "Some notes",
	}

	options := &Options{
		Values:     true, // Populate values from the struct
		IDPrefix:   "id_",
		NamePrefix: "name_",
	}

	inputs, order := BuildInputs(source, options)

	if len(inputs) != 6 {
		t.Fatalf("Expected 6 inputs, got %d", len(inputs))
	}
	if len(order) != 6 {
		t.Fatalf("Expected 6 fields in order, got %d", len(order))
	}

	inputMap := make(map[string]Input)
	for _, inp := range inputs {
		inputMap[inp.FieldName] = inp
	}

	usernameInput, ok := inputMap["Username"]
	if !ok {
		t.Fatal("Username input not found")
	}
	usernameHTML := string(usernameInput.HTML())
	if !strings.Contains(usernameHTML, `type="text"`) {
		t.Errorf("Username should be text type, got: %s", usernameHTML)
	}
	if !strings.Contains(usernameHTML, `id="id_Username"`) {
		t.Errorf("Username ID prefix failed, got: %s", usernameHTML)
	}
	if !strings.Contains(usernameHTML, `value="johndoe"`) {
		t.Errorf("Username value not populated, got: %s", usernameHTML)
	}
	if !strings.Contains(usernameHTML, `minlength="3"`) || !strings.Contains(usernameHTML, `maxlength="20"`) {
		t.Errorf("Username validation attributes missing, got: %s", usernameHTML)
	}

	passwordInput, ok := inputMap["Password"]
	if !ok {
		t.Fatal("Password input not found")
	}
	passwordHTML := string(passwordInput.HTML())
	if !strings.Contains(passwordHTML, `type="password"`) {
		t.Errorf("Password should be password type, got: %s", passwordHTML)
	}
	if strings.Contains(passwordHTML, `value=`) {
		t.Errorf("Password should NOT have a value attribute, got: %s", passwordHTML)
	}

	activeInput, ok := inputMap["Active"]
	if !ok {
		t.Fatal("Active input not found")
	}
	activeHTML := string(activeInput.HTML())
	if !strings.Contains(activeHTML, `type="checkbox"`) {
		t.Errorf("Active should be checkbox type, got: %s", activeHTML)
	}
	if !strings.Contains(activeHTML, `checked`) {
		t.Errorf("Active checkbox should be checked, got: %s", activeHTML)
	}

	notesInput, ok := inputMap["Notes"]
	if !ok {
		t.Fatal("Notes input not found")
	}
	notesHTML := string(notesInput.HTML())
	if !strings.Contains(notesHTML, `<textarea`) {
		t.Errorf("Notes should be textarea, got: %s", notesHTML)
	}
	if !strings.Contains(notesHTML, `Some notes`) {
		t.Errorf("Textarea value not populated, got: %s", notesHTML)
	}
}

func TestBuildInputsOrderPreservation(t *testing.T) {
	source := &TestStruct{}
	_, order := BuildInputs(source, &Options{})

	expectedOrder := []string{"Username", "Email", "Age", "Password", "Active", "Notes"}
	if len(order) != len(expectedOrder) {
		t.Fatalf("Order length mismatch: expected %d, got %d", len(expectedOrder), len(order))
	}

	for i, field := range expectedOrder {
		if order[i] != field {
			t.Errorf("Order mismatch at index %d: expected %s, got %s", i, field, order[i])
		}
	}
}

func TestBuildInputsExclusion(t *testing.T) {
	source := &TestStruct{}

	options := &Options{
		ExcludeFields: map[string]struct{}{
			"Age": {},
		},
	}

	inputs, order := BuildInputs(source, options)

	for _, inp := range inputs {
		if inp.FieldName == "Age" {
			t.Fatal("Age field should have been excluded")
		}
	}

	for _, name := range order {
		if name == "Age" {
			t.Fatal("Age should not be in the order list")
		}
	}

	if len(inputs) != 5 {
		t.Fatalf("Expected 5 inputs after exclusion, got %d", len(inputs))
	}
}
