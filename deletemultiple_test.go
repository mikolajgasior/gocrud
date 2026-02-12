package crud

import (
	"context"
	"testing"

	sqlfilters "miko.gs/pgsql-builder/pkg/filters"
)

// TestDeleteMultiple tests if DeleteMultiple removes objects from database based on specified filters
func TestDeleteMultiple(t *testing.T) {
	recreateTestStructTable()

	// Insert some data that should not be removed
	for i := 1; i < 51; i++ {
		ts := testStructWithData()
		ts.ID = 0
		ts.Age = 10 + i
		ts.PrimaryEmail = "another@example.com"
		_ = testCRUD.Save(context.Background(), ts, SaveOptions{})
	}

	// Insert data that should be deleted
	for i := 1; i < 151; i++ {
		ts := testStructWithData()
		ts.ID = 0
		ts.Age = 30
		_ = testCRUD.Save(context.Background(), ts, SaveOptions{})
	}

	// Delete multiple rows from the database
	err := testCRUD.DeleteMultiple(context.Background(), &TestStruct{}, DeleteMultipleOptions{
		Filters: &sqlfilters.Filters{
			"Price": {
				Op:  sqlfilters.OpEqual,
				Val: 444,
			},
			"PrimaryEmail": {
				Op:  sqlfilters.OpEqual,
				Val: "primary@example.com",
			},
		},
	})
	if err != nil {
		t.Fatalf("DeleteMultiple failed to delete objects: %s", err.(*CRUDError).Op)
	}

	cnt, _ := testCRUD.GetCount(context.Background(), &TestStruct{}, GetCountOptions{})
	if cnt != 50 {
		t.Fatalf("DeleteMultiple removed invalid number of rows, there are %d rows left, instead of %d", cnt, 50)
	}
}

// TestDeleteMultipleWithRawQuery tests if DeleteMultiple removes objects from the database based on specified filters
// and condition, that is almost a raw query.
func TestDeleteMultipleWithRawQuery(t *testing.T) {
	recreateTestStructTable()

	// Insert some data that should not be removed
	for i := 1; i < 51; i++ {
		ts := testStructWithData()
		ts.ID = 0
		ts.Age = 10 + i
		ts.PrimaryEmail = "another@example.com"
		_ = testCRUD.Save(context.Background(), ts, SaveOptions{})
	}

	// Insert data that should be deleted
	for i := 1; i < 151; i++ {
		ts := testStructWithData()
		ts.ID = 0
		ts.Age = 30
		_ = testCRUD.Save(context.Background(), ts, SaveOptions{})
	}

	// Delete multiple rows from the database
	err := testCRUD.DeleteMultiple(context.Background(), &TestStruct{}, DeleteMultipleOptions{
		Filters: &sqlfilters.Filters{
			"Price": {
				Op:  sqlfilters.OpEqual,
				Val: 444,
			},
			"PrimaryEmail": {
				Op:  sqlfilters.OpEqual,
				Val: "primary@example.com",
			},
			sqlfilters.Raw: {
				Op: sqlfilters.OpOR,
				Val: []interface{}{
					".Age = ? OR .Age IN (?) OR (.Age = ? AND .PrimaryEmail = ?)",
					31,
					[]int{32, 33, 34},
					35,
					"miko@example.com",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("DeleteMultiple failed to delete objects: %s", err.(*CRUDError).Op)
	}

	cnt, _ := testCRUD.GetCount(context.Background(), &TestStruct{}, GetCountOptions{})
	if cnt != 46 {
		t.Fatalf("DeleteMultiple removed invalid number of rows, there are %d rows left, instead of %d", cnt, 46)
	}
}

// TestDeleteMultipleWithRawQueryOnly tests if DeleteMultiple removes objects from database based on a condition which is somewhat raw query
func TestDeleteMultipleWithRawQueryOnly(t *testing.T) {
	recreateTestStructTable()

	// Insert some data that should not be removed
	for i := 1; i < 51; i++ {
		ts := testStructWithData()
		ts.ID = 0
		ts.Age = 10 + i
		ts.PrimaryEmail = "another@example.com"
		_ = testCRUD.Save(context.Background(), ts, SaveOptions{})
	}

	// Insert data that should be deleted
	for i := 1; i < 151; i++ {
		ts := testStructWithData()
		ts.ID = 0
		ts.Age = 30
		_ = testCRUD.Save(context.Background(), ts, SaveOptions{})
	}

	// Delete multiple rows from the database
	err := testCRUD.DeleteMultiple(context.Background(), &TestStruct{}, DeleteMultipleOptions{
		Filters: &sqlfilters.Filters{
			sqlfilters.Raw: {
				Op: sqlfilters.OpOR,
				Val: []interface{}{
					"(.Price = ? AND .PrimaryEmail = ?) OR (.Age = ? OR .Age IN (?) OR (.Age = ? AND .PrimaryEmail = ?))",
					444,
					"primary@example.com",
					31,
					[]int{32, 33, 34},
					35,
					"miko@example.com",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("DeleteMultiple failed to delete objects: %s", err.(*CRUDError).Op)
	}

	cnt, _ := testCRUD.GetCount(context.Background(), &TestStruct{}, GetCountOptions{})
	if cnt != 46 {
		t.Fatalf("DeleteMultiple removed invalid number of rows, there are %d rows left, instead of %d", cnt, 46)
	}
}
