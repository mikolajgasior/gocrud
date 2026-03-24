package service

import (
	"context"
	"testing"
	"time"

	"miko.gs/gocrud/pkg/test"
)

func TestSaveFromForm(t *testing.T) {
	recreateTestStructTable()

	objSaved := &test.TestStruct{}
	urlValues := test.TestStructURLValues()
	prefix := "prefix"
	now := time.Now().UTC().Unix()
	err := testService.SaveFromForm(context.Background(), objSaved, urlValues, prefix, now, 311)
	if err != nil {
		t.Fatalf("Save failed to insert struct to the table: %s", err.Error())
	}

	objLoaded := test.FatalIfNoTestStructInTheDatabase(t, testDB)
	if objLoaded.ID == 0 || !test.AreTestStructObjectsSame(objSaved, objLoaded) {
		t.Fatalf("Save failed to insert struct to the table")
	}

	// if creation works then no need for testing update
}
