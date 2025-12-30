package crud

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

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

// Test struct for all the tests
type TestStruct struct {
	ID             int64
	Flags          int64
	PrimaryEmail   string `crud:"req"`
	EmailSecondary string `crud:"req email"`
	FirstName      string `crud:"req lenmin:2 lenmax:30"`
	LastName       string `crud:"req lenmin:0 lenmax:255"`
	Age            int    `crud:"valmin:1 valmax:120"`
	Price          int    `crud:"valmin:0 valmax:999"`
	PostCode       string `crud:"req lenmin:6 regexp:^[0-9]{2}\\-[0-9]{3}$"`
	PostCode2      string `crud:"lenmin:6" 2db_regexp:"^[0-9]{2}\\-[0-9]{3}$"`
	Password       string `json:"password"`
	Key            string `crud:"req uniq lenmin:30 lenmax:255"`
	CreatedAt      int64
	CreatedBy      int64
	ModifiedAt     int64
	ModifiedBy     int64
}

func TestMain(m *testing.M) {
	createDocker()

	var code int

	defer func() {
		removeDocker()
		os.Exit(code)
	}()

	createOrm()
	code = m.Run()
}

func createDocker() {
	var err error
	dockerPool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	dockerResource, err = dockerPool.Run("postgres", "16-alpine", []string{"POSTGRES_PASSWORD=" + testPassword, "POSTGRES_USER=" + testUser, "POSTGRES_DB=" + testName})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	if err = dockerPool.Retry(func() error {
		var err error
		testDB, err = sql.Open("postgres", fmt.Sprintf("host=localhost user=%s password=%s port=%s dbname=%s sslmode=disable", testUser, testPassword, dockerResource.GetPort("5432/tcp"), testName))
		if err != nil {
			return err
		}
		return testDB.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err.Error())
	}
}

func createOrm() {
	testCRUD = New(testDB, Options{})
}

func removeDocker() {
	err := dockerPool.Purge(dockerResource)
	if err != nil {
		log.Fatalf("Error removing docker container: %s", err.Error())
	}
}

func testStructWithData() *TestStruct {
	ts := &TestStruct{}
	ts.Flags = 4
	ts.PrimaryEmail = "primary@example.com"
	ts.EmailSecondary = "secondary@example.com"
	ts.FirstName = "John"
	ts.LastName = "Smith"
	ts.Age = 37
	ts.Price = 444
	ts.PostCode = "00-000"
	ts.PostCode2 = "11-111"
	ts.Password = "yyy"
	ts.Key = fmt.Sprintf("12345679012345678901234567890%d", time.Now().UnixNano())
	return ts
}

func recreateTestStructTable() {
	_ = testCRUD.DropTable(context.Background(), &TestStruct{})
	_ = testCRUD.CreateTable(context.Background(), &TestStruct{})
}

func areTestStructObjectsSame(ts1 *TestStruct, ts2 *TestStruct) bool {
	if ts1.Flags != ts2.Flags {
		return false
	}
	if ts1.PrimaryEmail != ts2.PrimaryEmail {
		return false
	}
	if ts1.EmailSecondary != ts2.EmailSecondary {
		return false
	}
	if ts1.FirstName != ts2.FirstName {
		return false
	}
	if ts1.LastName != ts2.LastName {
		return false
	}
	if ts1.Age != ts2.Age {
		return false
	}
	if ts1.Price != ts2.Price {
		return false
	}
	if ts1.PostCode != ts2.PostCode {
		return false
	}
	if ts1.PostCode2 != ts2.PostCode2 {
		return false
	}
	if ts1.Password != ts2.Password {
		return false
	}
	if ts1.CreatedBy != ts2.CreatedBy {
		return false
	}
	if ts1.Key != ts2.Key {
		return false
	}
	return true
}
