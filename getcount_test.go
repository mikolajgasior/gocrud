package crud

import (
	"context"
	"testing"

	sqlbuilder "github.com/keenbytes/pgsql-builder"
)

func TestGetCount(t *testing.T) {
	recreateTestStructTable()

	// Insert some data that should be ignored by GetCount later on
	for i := 1; i < 51; i++ {
		ts := testStructWithData()
		ts.ID = 0
		ts.Age = 10 + i
		ts.Price = 444
		ts.PrimaryEmail = "another@example.com"
		_ = testCRUD.Save(context.Background(), ts, SaveOptions{})
	}

	// Insert data that should be selected by filters
	for i := 1; i < 151; i++ {
		ts := testStructWithData()
		ts.ID = 0
		ts.Age = 30
		_ = testCRUD.Save(context.Background(), ts, SaveOptions{})
	}

	// Get the data from the database
	cnt, err := testCRUD.GetCount(context.Background(), &TestStruct{}, GetCountOptions{
		Filters: &sqlbuilder.Filters{
			"Price":        {Op: sqlbuilder.OpEqual, Val: 444},
			"PrimaryEmail": {Op: sqlbuilder.OpEqual, Val: "primary@example.com"},
		},
	})
	if err != nil {
		t.Fatalf("Get failed to return list of objects: %s", err.(ErrCRUD).Op)
	}
	if cnt != 150 {
		t.Fatalf("Get failed to return list of objects, want %v, got %v", 150, cnt)
	}
}

func TestGetCountWithRawQuery(t *testing.T) {
	recreateTestStructTable()

	// Insert some data that should be ignored by GetCount later on
	for i := 1; i < 51; i++ {
		ts := testStructWithData()
		ts.ID = 0
		ts.Age = 10 + i
		ts.Price = 444
		ts.PrimaryEmail = "another@example.com"
		_ = testCRUD.Save(context.Background(), ts, SaveOptions{})
	}

	// Insert data that should be selected by filters
	for i := 1; i < 151; i++ {
		ts := testStructWithData()
		ts.ID = 0
		ts.Age = 30
		_ = testCRUD.Save(context.Background(), ts, SaveOptions{})
	}

	// Get the data from the database
	cnt, err := testCRUD.GetCount(context.Background(), &TestStruct{}, GetCountOptions{
		Filters: &sqlbuilder.Filters{
			"Price":        {Op: sqlbuilder.OpEqual, Val: 444},
			"PrimaryEmail": {Op: sqlbuilder.OpEqual, Val: "primary@example.com"},
			sqlbuilder.Raw: {
				Op: sqlbuilder.OpOR,
				Val: []interface{}{
					".PrimaryEmail = ? AND .Age IN (?)",
					"another@example.com",
					[]int{32, 33, 34},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Get failed to return list of objects: %s", err.(ErrCRUD).Op)
	}
	if cnt != 153 {
		t.Fatalf("Get failed to return list of objects, want %v, got %v", 153, cnt)
	}
}
