package gocrud

import (
	"context"
	"fmt"
	"testing"
)

// PasswordStruct is a minimal struct with a field tagged crud:"pass".
// It maps to the table "password_struct".
type PasswordStruct struct {
	ID       uint64
	Name     string `crud:"req len:1,100"`
	Password string `crud:"pass"`
}

func recreatePasswordStructTable() {
	_ = testCRUD.DropTable(context.Background(), &PasswordStruct{})
	_ = testCRUD.CreateTable(context.Background(), &PasswordStruct{})
}

// TestLoad_PasswordFieldIsZeroed verifies that Load returns an object with the
// password field cleared even though a bcrypt hash exists in the database.
func TestLoad_PasswordFieldIsZeroed(t *testing.T) {
	recreatePasswordStructTable()

	obj := &PasswordStruct{Name: "alice", Password: "secret123"}
	if err := testCRUD.Save(context.Background(), obj, SaveOptions{}); err != nil {
		t.Fatalf("Save failed: %s", err.Error())
	}

	loaded := &PasswordStruct{}
	if output := testCRUD.Load(context.Background(), loaded, fmt.Sprintf("%d", obj.ID), LoadOptions{}); output.Error != nil {
		t.Fatalf("Load failed: %s", output.Error.Error())
	}

	if loaded.Password != "" {
		t.Errorf("Load: expected password field to be empty, got %q", loaded.Password)
	}
	if loaded.Name != "alice" {
		t.Errorf("Load: expected Name to be %q, got %q", "alice", loaded.Name)
	}
}

// TestGet_PasswordFieldIsZeroed verifies that Get returns objects with the
// password field cleared.
func TestGet_PasswordFieldIsZeroed(t *testing.T) {
	recreatePasswordStructTable()

	for _, name := range []string{"bob", "carol"} {
		obj := &PasswordStruct{Name: name, Password: "hunter2"}
		if err := testCRUD.Save(context.Background(), obj, SaveOptions{}); err != nil {
			t.Fatalf("Save failed: %s", err.Error())
		}
	}

	results, err := testCRUD.Get(context.Background(), func() interface{} {
		return &PasswordStruct{}
	}, GetOptions{})
	if err != nil {
		t.Fatalf("Get failed: %s", err.Error())
	}

	if len(results) != 2 {
		t.Fatalf("Get: expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		ps := r.(*PasswordStruct)
		if ps.Password != "" {
			t.Errorf("Get: expected password field to be empty for %q, got %q", ps.Name, ps.Password)
		}
	}
}

// TestLoad_VerifyPasswordFields verifies that Load reports PassOK and
// PassInvalid for the password fields listed in
// LoadOptions.VerifyPasswordFields, and ignores keys that do not name an
// actual password field.
func TestLoad_VerifyPasswordFields(t *testing.T) {
	recreatePasswordStructTable()

	obj := &PasswordStruct{Name: "dave", Password: "correct-horse"}
	if err := testCRUD.Save(context.Background(), obj, SaveOptions{}); err != nil {
		t.Fatalf("Save failed: %s", err.Error())
	}

	loaded := &PasswordStruct{}
	output := testCRUD.Load(context.Background(), loaded, fmt.Sprintf("%d", obj.ID), LoadOptions{
		VerifyPasswordFields: map[string]string{
			"Password": "correct-horse",
			"Name":     "dave",
		},
	})
	if output.Error != nil {
		t.Fatalf("Load failed: %s", output.Error.Error())
	}

	if got := output.PasswordFields["Password"]; got != PassOK {
		t.Errorf("expected PassOK for correct password, got %d", got)
	}
	if _, ok := output.PasswordFields["Name"]; ok {
		t.Errorf("expected no entry for non-password field, got %d", output.PasswordFields["Name"])
	}
	if len(output.PasswordFields) != 1 {
		t.Errorf("expected exactly 1 entry in PasswordFields, got %d", len(output.PasswordFields))
	}

	loaded2 := &PasswordStruct{}
	output2 := testCRUD.Load(context.Background(), loaded2, fmt.Sprintf("%d", obj.ID), LoadOptions{
		VerifyPasswordFields: map[string]string{
			"Password": "wrong-password",
		},
	})
	if output2.Error != nil {
		t.Fatalf("Load failed: %s", output2.Error.Error())
	}
	if got := output2.PasswordFields["Password"]; got != PassInvalid {
		t.Errorf("expected PassInvalid for wrong password, got %d", got)
	}

	if loaded2.Password != "" {
		t.Errorf("Load: expected password field to be empty, got %q", loaded2.Password)
	}
}
