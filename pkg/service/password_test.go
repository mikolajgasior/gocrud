package service

import (
	"context"
	"testing"

	"github.com/mikolajgasior/gocrud"
)

// PasswordStruct is a minimal struct with a field tagged crud:"pass",
// registered in testService under the "passwordstruct" key.
type PasswordStruct struct {
	ID       uint64
	Name     string `crud:"req len:1,100"`
	Password string `crud:"pass"`
}

// TestList_VerifyPasswordFields verifies that List passes passFieldsToVerify
// through to gocrud's Get, reporting PassOK/PassInvalid for a query that
// matches exactly one record and ignoring keys that aren't password fields.
func TestList_VerifyPasswordFields(t *testing.T) {
	recreatePasswordStructTable()

	obj := &PasswordStruct{Name: "erin", Password: "correct-horse"}
	if err := testCRUD.Save(context.Background(), obj, gocrud.SaveOptions{}); err != nil {
		t.Fatalf("Save failed: %s", err.Error())
	}

	results, passwordFields, err := testService.List(context.Background(), "passwordstruct", 10, 0, "", "",
		map[string]string{"Name": "erin"}, nil, nil, nil,
		map[string]string{
			"Password": "correct-horse",
			"Name":     "erin",
		})
	if err != nil {
		t.Fatalf("List failed: %s", err.Error())
	}
	if len(results) != 1 {
		t.Fatalf("List: expected 1 result, got %d", len(results))
	}
	if got := passwordFields["Password"]; got != gocrud.PassOK {
		t.Errorf("expected PassOK for correct password, got %d", got)
	}
	if _, ok := passwordFields["Name"]; ok {
		t.Errorf("expected no entry for non-password field, got %d", passwordFields["Name"])
	}
	if ps := results[0].(*PasswordStruct); ps.Password != "" {
		t.Errorf("List: expected password field to be empty, got %q", ps.Password)
	}

	_, passwordFields2, err := testService.List(context.Background(), "passwordstruct", 10, 0, "", "",
		map[string]string{"Name": "erin"}, nil, nil, nil,
		map[string]string{"Password": "wrong-password"})
	if err != nil {
		t.Fatalf("List failed: %s", err.Error())
	}
	if got := passwordFields2["Password"]; got != gocrud.PassInvalid {
		t.Errorf("expected PassInvalid for wrong password, got %d", got)
	}
}
