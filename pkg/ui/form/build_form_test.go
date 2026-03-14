package form

import (
	"strings"
	"testing"

	uiinput "miko.gs/struct-crud/pkg/ui/input"
)

type TestStruct struct {
	Username string `html:"req len:3,20"`
	Email    string `html:"req uiemail"`
	Age      int    `html:"req val:18,120"`
	Password string `html:"uipassword req"`
	Active   bool   `html:""`
	Notes    string `html:"len:0,500 uitextarea"`
}

func TestBuildFormAndRender(t *testing.T) {
	source := &TestStruct{
		Username: "johndoe",
		Email:    "john@example.com",
		Age:      30,
		Password: "secret123",
		Active:   true,
		Notes:    "Some notes",
	}

	formOpts := &FormOptions{
		Options: uiinput.Options{
			Values:     true,
			IDPrefix:   "id_",
			NamePrefix: "name_",
		},
		Path: "/register-user",
	}

	form, err := BuildForm(source, formOpts)
	if err != nil {
		t.Fatalf("BuildForm failed: %v", err)
	}

	if form.Path != "/register-user" {
		t.Errorf("Expected Path '/register-user', got '%s'", form.Path)
	}
	if len(form.InputsOrder) != 6 {
		t.Fatalf("Expected 6 fields in order, got %d", len(form.InputsOrder))
	}
	if len(form.Inputs) != 6 {
		t.Fatalf("Expected 6 inputs in map, got %d", len(form.Inputs))
	}

	html := string(form.HTML())

	if !strings.Contains(html, `<form`) {
		t.Fatal("Rendered HTML missing <form> tag")
	}
	if !strings.Contains(html, `action="/register-user"`) {
		t.Fatal("Rendered HTML missing correct action path")
	}
	if !strings.Contains(html, `method="POST"`) {
		t.Fatal("Rendered HTML missing method POST")
	}
	if !strings.Contains(html, `<button type="submit">Create</button>`) {
		t.Fatal("Rendered HTML missing submit button")
	}

	expectedLabels := []string{"Username", "Email", "Age", "Password", "Active", "Notes"}
	for _, label := range expectedLabels {
		if !strings.Contains(html, `<label>`+label+`:</label>`) {
			t.Errorf("Missing label for field: %s", label)
		}
	}

	usernameIdx := strings.Index(html, "Username")
	emailIdx := strings.Index(html, "Email")
	ageIdx := strings.Index(html, "Age")
	passwordIdx := strings.Index(html, "Password")
	activeIdx := strings.Index(html, "Active")
	notesIdx := strings.Index(html, "Notes")

	if !(usernameIdx < emailIdx && emailIdx < ageIdx && ageIdx < passwordIdx && passwordIdx < activeIdx && activeIdx < notesIdx) {
		t.Fatalf("Inputs not in correct order. Indices: U=%d, E=%d, A=%d, P=%d, Ac=%d, N=%d",
			usernameIdx, emailIdx, ageIdx, passwordIdx, activeIdx, notesIdx)
	}

	if !strings.Contains(html, `id="id_Username"`) {
		t.Error("Username ID prefix missing in form")
	}
	if !strings.Contains(html, `value="johndoe"`) {
		t.Error("Username value missing in form")
	}

	if strings.Contains(html, `type="password"`) {
		pwdStart := strings.Index(html, `type="password"`)
		if pwdStart == -1 {
			t.Fatal("Password input type not found")
		}

		pwdBlock := html[pwdStart:]
		if len(pwdBlock) > 200 {
			pwdBlock = pwdBlock[:200]
		}
		if strings.Contains(pwdBlock, `value=`) {
			t.Error("Password field should not have a value attribute in the form")
		}
	}

	if !strings.Contains(html, `type="checkbox"`) {
		t.Error("Checkbox input type missing")
	}

	if !strings.Contains(html, `checked`) {
		t.Error("Active checkbox should be checked")
	}

	if !strings.Contains(html, `<textarea`) {
		t.Error("Textarea input missing")
	}
	if !strings.Contains(html, `Some notes`) {
		t.Error("Textarea value missing")
	}
}

func TestBuildFormEmptyPath(t *testing.T) {
	source := &TestStruct{}

	form, err := BuildForm(source, &FormOptions{
		Options: uiinput.Options{},
		Path:    "", // Explicitly empty
	})
	if err != nil {
		t.Fatalf("BuildForm failed: %v", err)
	}

	html := string(form.HTML())

	if !strings.Contains(html, `action=""`) {
		t.Error("Form should have empty action attribute if path is empty")
	}
}

func TestBuildFormExclusion(t *testing.T) {
	source := &TestStruct{}

	form, err := BuildForm(source, &FormOptions{
		Options: uiinput.Options{
			ExcludeFields: map[string]struct{}{
				"Age": {},
			},
		},
		Path: "/test",
	})
	if err != nil {
		t.Fatalf("BuildForm failed: %v", err)
	}

	html := string(form.HTML())

	if strings.Contains(html, "Age") {
		t.Error("Age field should be excluded from the form")
	}

	if len(form.InputsOrder) != 5 {
		t.Fatalf("Expected 5 fields after exclusion, got %d", len(form.InputsOrder))
	}
}
