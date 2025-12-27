package crud

import "testing"

type TestCompany struct {
	ID            int64
	Name          string
	TestEmployees []TestEmployee `crud:"on_del:del del_field:TestCompanyID"`
}

type TestEmployee struct {
	ID              int64
	Name            string
	TestCompanyID   int64
	TestCreditCards []TestCreditCard `crud:"on_del:del del_field:TestEmployeeID"`
	TestComments    []TestComment    `crud:"on_del:upd del_field:TestEmployeeID del_upd_field:TestEmployeeID del_upd_val:0"`
}

type TestComment struct {
	ID             int64
	Comment        string
	TestEmployeeID int64
}

type TestCreditCard struct {
	ID             int64
	Number         string
	TestEmployeeID int64
}

func createTestCascadeDeleteObjects() interface{} {
	recreateTestDeleteCascadeTables()

	// create two companies
	i := 1

	company := &TestCompany{
		Name: "Company",
		ID:   int64(i),
	}
	_ = testCRUD.Save(company, SaveOptions{})

	// create two employees in each company
	for j := 1; j < 3; j++ {
		employee := &TestEmployee{
			Name:          "Employee",
			ID:            int64(i*10 + j),
			TestCompanyID: int64(i),
		}
		_ = testCRUD.Save(employee, SaveOptions{})

		// create two credit cards and two comments for each employee
		for k := 1; k < 3; k++ {
			creditCard := &TestCreditCard{
				ID:             int64(i*100 + j*10 + k),
				Number:         "number",
				TestEmployeeID: int64(i*10 + j),
			}
			_ = testCRUD.Save(creditCard, SaveOptions{})

			comment := &TestComment{
				ID:             int64(i*100 + j*10 + k),
				Comment:        "comment",
				TestEmployeeID: int64(i*10 + j),
			}
			_ = testCRUD.Save(comment, SaveOptions{})
		}
	}

	return company
}

func recreateTestDeleteCascadeTables() {
	_ = testCRUD.DropTable(&TestCompany{})
	_ = testCRUD.DropTable(&TestEmployee{})
	_ = testCRUD.DropTable(&TestCreditCard{})
	_ = testCRUD.DropTable(&TestComment{})
	_ = testCRUD.CreateTable(&TestCompany{})
	_ = testCRUD.CreateTable(&TestEmployee{})
	_ = testCRUD.CreateTable(&TestCreditCard{})
	_ = testCRUD.CreateTable(&TestComment{})
}

func TestDeleteCascade(t *testing.T) {
	// Create a test parent (with children) and get its ID.
	p := createTestCascadeDeleteObjects()

	// Check if test objects are added to the database.
	var cnt int
	err2 := testDB.QueryRow("SELECT COUNT(*) FROM test_company").Scan(&cnt)
	if err2 != nil {
		t.Fatalf("Failed to select count: %s", err2.Error())
	}
	if cnt != 1 {
		t.Fatalf("Number of parent objects in the database before running the test is invalid")
	}
	err2 = testDB.QueryRow("SELECT COUNT(*) FROM test_employee WHERE id IN (11, 12) AND test_company_id != 0").Scan(&cnt)
	if err2 != nil {
		t.Fatalf("Failed to select count: %s", err2.Error())
	}
	if cnt != 2 {
		t.Fatalf("Number of children objects in the database before running the test is invalid")
	}
	err2 = testDB.QueryRow("SELECT COUNT(*) FROM test_credit_card WHERE id IN (111, 112, 121, 122) AND test_employee_id != 0").Scan(&cnt)
	if err2 != nil {
		t.Fatalf("Failed to select count: %s", err2.Error())
	}
	if cnt != 4 {
		t.Fatalf("Number of grand children objects in the database before running the test is invalid")
	}
	err2 = testDB.QueryRow("SELECT COUNT(*) FROM test_comment WHERE id IN (111, 112, 121, 122) AND test_employee_id != 0").Scan(&cnt)
	if err2 != nil {
		t.Fatalf("Failed to select count: %s", err2.Error())
	}
	if cnt != 4 {
		t.Fatalf("Number of grand children objects in the database before running the test is invalid")
	}

	// Delete the parent (company) object.
	err1 := testCRUD.Delete(p, DeleteOptions{})
	if err1 != nil {
		t.Fatalf("Failed to run Delete successfully: %v", err1.(ErrCRUD).Op)
	}

	// Parent should be removed.
	err2 = testDB.QueryRow("SELECT COUNT(*) FROM test_company").Scan(&cnt)
	if err2 != nil {
		t.Fatalf("Failed to select count: %s", err2.Error())
	}
	if cnt > 0 {
		t.Fatalf("Delete failed to remove parent object")
	}

	// Children should be removed (employees).
	err2 = testDB.QueryRow("SELECT COUNT(*) FROM test_employee WHERE id IN (11, 12)").Scan(&cnt)
	if err2 != nil {
		t.Fatalf("Failed to select count: %s", err2.Error())
	}
	if cnt != 0 {
		t.Fatalf("Delete failed to remove children objects")
	}

	// Grand children, tagged to be deleted, should be removed (credit cards).
	err2 = testDB.QueryRow("SELECT COUNT(*) FROM test_credit_card WHERE id IN (111, 112, 121, 122)").Scan(&cnt)
	if err2 != nil {
		t.Fatalf("Failed to select count: %s", err2.Error())
	}
	if cnt > 0 {
		t.Fatalf("Delete failed to remove grand children")
	}

	// Grand children, tagged to be updated, should be updated (comments).
	err2 = testDB.QueryRow("SELECT COUNT(*) FROM test_comment WHERE id IN (111, 112, 121, 122) AND test_employee_id=0").Scan(&cnt)
	if err2 != nil {
		t.Fatalf("Failed to select count: %s", err2.Error())
	}
	if cnt != 4 {
		t.Fatalf("Delete failed to not update grand children")
	}
}
