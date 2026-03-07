package api

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	structcrud "miko.gs/struct-crud"
	"miko.gs/struct-crud/pkg/http/cors"
	"miko.gs/struct-crud/pkg/service"
	"miko.gs/struct-crud/pkg/test"
)

var testUser = "testuser"
var testPassword = "testpass"
var testName = "testdb"
var testDB *sql.DB

var dockerPool *dockertest.Pool
var dockerResource *dockertest.Resource

var testCRUD *structcrud.CRUD
var testService *service.CRUD
var testHandler *Handler

func TestMain(m *testing.M) {
	dockerPool, dockerResource, testDB = test.CreateDocker(testUser, testPassword, testName)

	var code int

	defer func() {
		test.RemoveDocker(dockerPool, dockerResource)
		os.Exit(code)
	}()

	createCRUD()
	createService()
	createHandler()

	code = m.Run()
}

func createCRUD() {
	testCRUD = structcrud.New(testDB, structcrud.Options{})
}

func createService() {
	testService = service.New(map[string]func() interface{}{
		"teststruct": func() interface{} {
			return &test.TestStruct{}
		},
	}, testDB)
}

func createHandler() {
	testHandler = &Handler{
		cors: &cors.CORS{},
		paths: map[string]func() interface{}{
			"teststruct": func() interface{} {
				return &test.TestStruct{}
			},
		},
		svc: testService,
	}
}

func recreateTestStructTable() {
	_ = testCRUD.DropTable(context.Background(), &test.TestStruct{})
	_ = testCRUD.CreateTable(context.Background(), &test.TestStruct{})
}
