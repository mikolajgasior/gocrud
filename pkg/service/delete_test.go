package service

import (
	"context"
	"testing"

	"codeberg.org/mikolajgasior/gocrud"
	"codeberg.org/mikolajgasior/gocrud/pkg/test"
)

func TestDelete(t *testing.T) {
	recreateTestStructTable()

	// Insert an object first
	objSaved := test.TestStructWithData()
	_ = testCRUD.Save(context.Background(), objSaved, crud.SaveOptions{})

	// Delete it
	err := testService.Delete(context.Background(), "teststruct", objSaved.ID)
	if err != nil {
		t.Fatalf("Delete failed to remove: %s", err.Error())
	}

	test.FatalIfTestStructNotDeletedInTheDatabase(t, testDB, objSaved.ID)
}
