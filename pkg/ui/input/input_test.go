package input

import (
	"strings"
	"testing"
)

func TestInputHTML_Text(t *testing.T) {
	inp := Input{
		FieldName: "Username",
		InputType: TypeText,
		Name:      "username",
		ID:        "id_username",
		Required:  true,
		Value:     "test_user",
	}

	html := string(inp.HTML())
	if !strings.Contains(html, `<input type="text"`) {
		t.Fatal("Text input type not found")
	}
	if !strings.Contains(html, `name="username"`) {
		t.Fatal("Name attribute missing")
	}
	if !strings.Contains(html, `id="id_username"`) {
		t.Fatal("ID attribute missing")
	}
	if !strings.Contains(html, `required`) {
		t.Fatal("Required attribute missing")
	}
	if !strings.Contains(html, `value="test_user"`) {
		t.Fatal("Value attribute missing")
	}
}

func TestInputHTML_Password(t *testing.T) {
	inp := Input{
		FieldName: "Password",
		InputType: TypePassword,
		Name:      "password",
		ID:        "id_password",
		Required:  true,
		Value:     "secret123", // Should be ignored for password
	}

	html := string(inp.HTML())

	if !strings.Contains(html, `<input type="password"`) {
		t.Fatal("Password input type not found")
	}
	if strings.Contains(html, `value=`) {
		t.Fatal("Password should not have value attribute")
	}
}

func TestInputHTML_Checkbox(t *testing.T) {
	inp := Input{
		FieldName: "Subscribe",
		InputType: TypeCheckbox,
		Name:      "subscribe",
		ID:        "id_subscribe",
		Checked:   true,
	}

	html := string(inp.HTML())

	if !strings.Contains(html, `<input type="checkbox"`) {
		t.Fatal("Checkbox input type not found")
	}
	if !strings.Contains(html, `checked`) {
		t.Fatal("Checked attribute missing")
	}
	if strings.Contains(html, `value=`) {
		t.Fatal("Checkbox should not have value attribute")
	}
}

func TestInputHTML_Textarea(t *testing.T) {
	inp := Input{
		FieldName: "Bio",
		InputType: TypeTextarea,
		Name:      "bio",
		ID:        "id_bio",
		MinLength: 10,
		MaxLength: 500,
		Value:     "Hello world",
	}

	html := string(inp.HTML())

	if !strings.Contains(html, `<textarea`) {
		t.Fatal("Textarea tag not found")
	}
	if !strings.Contains(html, `minlength="10"`) {
		t.Fatal("Minlength attribute missing")
	}
	if !strings.Contains(html, `maxlength="500"`) {
		t.Fatal("Maxlength attribute missing")
	}
	if !strings.Contains(html, `Hello world`) {
		t.Fatal("Textarea value missing")
	}
}

func TestInputHTML_Escaping(t *testing.T) {
	inp := Input{
		FieldName: "Test",
		InputType: TypeText,
		Name:      "test",
		Value:     `"<script>alert('xss')</script>"`,
	}

	html := string(inp.HTML())
	if strings.Contains(html, `<script>`) {
		t.Fatal("Script tag not escaped")
	}
	if !strings.Contains(html, `&lt;script&gt;`) {
		t.Fatal("Script tag not properly escaped")
	}
}
