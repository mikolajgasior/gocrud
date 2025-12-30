package crud

import (
	"context"
	"fmt"
	"testing"
)

// TestDelete tests if Delete removes an object from the database.
func TestDelete(t *testing.T) {
	recreateTestStructTable()

	// Insert an object first
	objSaved := testStructWithData()
	_ = testCRUD.Save(context.Background(), objSaved, SaveOptions{})

	// Delete it
	err := testCRUD.Delete(context.Background(), objSaved, DeleteOptions{})
	if err != nil {
		t.Fatalf("Delete failed to remove: %s", err.(ErrCRUD).Op)
	}

	var rowCount int64
	err = testDB.QueryRow(fmt.Sprintf("SELECT COUNT(*) AS c FROM test_struct WHERE id = %d", objSaved.ID)).Scan(&rowCount)
	if err != nil {
		t.Fatalf("Delete failed to delete struct from the table")
	}
	if rowCount > 0 {
		t.Fatalf("Delete failed to delete struct from the table")
	}
	if objSaved.ID != 0 {
		t.Fatalf("Delete failed to set ID to 0 on the struct")
	}
}
