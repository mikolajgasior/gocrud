package service

import (
	"context"
	"database/sql"
	"os"
	"testing"

	structcrud "codeberg.org/mikolajgasior/gocrud"
	"codeberg.org/mikolajgasior/gocrud/pkg/test"
	"github.com/ory/dockertest/v3"
)

var testUser = "testuser"
var testPassword = "testpass"
var testName = "testdb"
var testDB *sql.DB

var dockerPool *dockertest.Pool
var dockerResource *dockertest.Resource

var testCRUD *structcrud.CRUD
var testService *CRUD

func TestMain(m *testing.M) {
	dockerPool, dockerResource, testDB = test.CreateDocker(testUser, testPassword, testName)

	var code int

	defer func() {
		test.RemoveDocker(dockerPool, dockerResource)
		os.Exit(code)
	}()

	createCRUD()
	createService()

	code = m.Run()
}

func createCRUD() {
	testCRUD = structcrud.New(testDB, structcrud.Options{Dialect: structcrud.DialectPostgres})
}

func createService() {
	testService = New(map[string]func() interface{}{
		"teststruct": func() interface{} {
			return &test.TestStruct{}
		},
	}, testDB, structcrud.DialectPostgres)
}

func recreateTestStructTable() {
	_ = testCRUD.DropTable(context.Background(), &test.TestStruct{})
	_ = testCRUD.CreateTable(context.Background(), &test.TestStruct{})
}
