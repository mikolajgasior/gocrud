package crud

import (
	"fmt"
	"testing"
)

// TestLoad tests if Load properly gets row from the database table and populate object fields with its value
func TestLoad(t *testing.T) {
	recreateTestStructTable()

	// Insert an object first
	objSaved := testStructWithData()
	_ = testCRUD.Save(objSaved, SaveOptions{})

	// Get the object
	objLoaded := &TestStruct{}
	err := testCRUD.Load(objLoaded, fmt.Sprintf("%d", objSaved.Id), LoadOptions{})
	if err != nil {
		t.Fatalf("Load failed to get data: %s", err.(ErrCRUD).Op)
	}

	if !areTestStructObjectsSame(objSaved, objLoaded) {
		t.Fatalf("Load failed to set struct with data: %s", err.(ErrCRUD).Op)
	}
}
