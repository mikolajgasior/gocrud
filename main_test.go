package crud

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"codeberg.org/mikolajgasior/gocrud/pkg/test"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
)

var testUser = "testuser"
var testPassword = "testpass"
var testName = "testdb"
var testDB *sql.DB

var dockerPool *dockertest.Pool
var dockerResource *dockertest.Resource

var testCRUD *CRUD

func TestMain(m *testing.M) {
	dockerPool, dockerResource, testDB = test.CreateDocker(testUser, testPassword, testName)

	var code int

	defer func() {
		test.RemoveDocker(dockerPool, dockerResource)
		os.Exit(code)
	}()

	createCRUD()
	code = m.Run()
}

func createCRUD() {
	testCRUD = New(testDB, Options{})
}

func recreateTestStructTable() {
	_ = testCRUD.DropTable(context.Background(), &test.TestStruct{})
	_ = testCRUD.CreateTable(context.Background(), &test.TestStruct{})
}
