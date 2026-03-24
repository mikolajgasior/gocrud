package service

import (
	"context"
	"testing"

	structcrud "miko.gs/gocrud"
	"miko.gs/gocrud/pkg/test"
)

func TestDelete(t *testing.T) {
	recreateTestStructTable()

	// Insert an object first
	objSaved := test.TestStructWithData()
	_ = testCRUD.Save(context.Background(), objSaved, structcrud.SaveOptions{})

	// Delete it
	err := testService.Delete(context.Background(), "teststruct", objSaved.ID)
	if err != nil {
		t.Fatalf("Delete failed to remove: %s", err.Error())
	}

	test.FatalIfTestStructNotDeletedInTheDatabase(t, testDB, objSaved.ID)
}
