package postgresql

import (
	"fmt"
	"log"
	"testing"

	"github.com/prisma/prisma-client-go/test/cmd"
	"github.com/prisma/prisma-client-go/test/setup"
)

const containerName = "go-client-postgres"
const image = "postgres:12.5-alpine"

var PostgreSQL = &postgreSQL{}

type postgreSQL struct{}

func (*postgreSQL) Name() string {
	return "postgresql"
}

func (*postgreSQL) ConnectionString(mockDBName string) string {
	return fmt.Sprintf("postgresql://postgres:pw@localhost:5433/%s", mockDBName)
}

func (*postgreSQL) SetupDatabase(t *testing.T) string {
	mockDB := setup.RandomString()

	exec(t, fmt.Sprintf("CREATE DATABASE %s", mockDB))

	return mockDB
}

func (*postgreSQL) TeardownDatabase(t *testing.T, mockDB string) {
	exec(t, fmt.Sprintf("DROP DATABASE %s", mockDB))
}

func exec(t *testing.T, query string) {
	if err := cmd.Run("docker", "exec", "-t", containerName, "psql", "-U", "postgres", "-c", query); err != nil {
		t.Fatal(err)
	}
}

func (db *postgreSQL) Setup() {
	if err := cmd.Run("docker", "run", "--name", containerName, "-p", "5433:5432", "-e", "POSTGRES_PASSWORD=pw", "-d", image); err != nil {
		panic(err)
	}
}

func (db *postgreSQL) Teardown() {
	if err := cmd.Run("docker", "stop", containerName); err != nil {
		log.Println(err)
	}

	if err := cmd.Run("docker", "rm", containerName, "--force"); err != nil {
		log.Println(err)
	}
}
