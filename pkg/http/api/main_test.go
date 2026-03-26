package api

import (
	"context"
	"database/sql"
	"os"
	"testing"

	structcrud "codeberg.org/mikolajgasior/gocrud"
	"codeberg.org/mikolajgasior/gocrud/pkg/http/cors"
	"codeberg.org/mikolajgasior/gocrud/pkg/service"
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
