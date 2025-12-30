package crud

import (
	"context"
	"testing"

	sqlbuilder "github.com/keenbytes/pgsql-builder"
)

// TestUpdateMultiple tests if UpdateMultiple update objects from database based on specified filters
func TestUpdateMultiple(t *testing.T) {
	recreateTestStructTable()

	// Insert some data that should not be removed
	for i := 1; i < 51; i++ {
		ts := testStructWithData()
		ts.ID = 0
		ts.Age = 10 + i
		ts.PrimaryEmail = "another@example.com"
		_ = testCRUD.Save(context.Background(), ts, SaveOptions{})
	}

	// Insert data that should be updated
	for i := 1; i < 151; i++ {
		ts := testStructWithData()
		ts.ID = 0
		ts.Age = 30
		ts.PrimaryEmail = "changeme@example.com"
		_ = testCRUD.Save(context.Background(), ts, SaveOptions{})
	}

	// Update multiple rows from the database
	err := testCRUD.UpdateMultiple(context.Background(), &TestStruct{}, map[string]interface{}{
		"PrimaryEmail": "newemail@example.com",
		"Age":          98,
	},
		UpdateMultipleOptions{
			Filters: &sqlbuilder.Filters{
				"Price": {
					Op:  sqlbuilder.OpEqual,
					Val: 444,
				},
				"PrimaryEmail": {
					Op:  sqlbuilder.OpEqual,
					Val: "changeme@example.com",
				},
			},
		})
	if err != nil {
		t.Fatalf("UpdateMultiple failed to update objects: %s %s", err.(ErrCRUD).Op, err.(ErrCRUD).Err.Error())
	}

	cnt, _ := testCRUD.GetCount(context.Background(), &TestStruct{}, GetCountOptions{
		Filters: &sqlbuilder.Filters{
			"PrimaryEmail": {
				Op:  sqlbuilder.OpEqual,
				Val: "newemail@example.com",
			},
			"Age": {
				Op:  sqlbuilder.OpEqual,
				Val: 98,
			},
		},
	})
	if cnt != 150 {
		t.Fatalf("UpdateMultiple updated invalid number of rows, there are %d rows left, instead of %d", cnt, 150)
	}
}
