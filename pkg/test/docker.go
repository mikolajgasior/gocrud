package test

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ory/dockertest/v3"
)

func CreateDocker(user, pass, name string) (dockerPool *dockertest.Pool, dockerResource *dockertest.Resource, db *sql.DB) {
	var err error
	dockerPool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	dockerResource, err = dockerPool.Run("postgres", "16-alpine", []string{"POSTGRES_PASSWORD=" + pass, "POSTGRES_USER=" + user, "POSTGRES_DB=" + name})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	if err = dockerPool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", fmt.Sprintf("host=localhost user=%s password=%s port=%s dbname=%s sslmode=disable", user, pass, dockerResource.GetPort("5432/tcp"), name))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err.Error())
	}

	return
}

func RemoveDocker(dockerPool *dockertest.Pool, dockerResource *dockertest.Resource) {
	err := dockerPool.Purge(dockerResource)
	if err != nil {
		log.Fatalf("Error removing docker container: %s", err.Error())
	}
}
