package postgresql

import (
	"fmt"
	"testing"

	"github.com/steebchen/prisma-client-go/test/cmd"
	"github.com/steebchen/prisma-client-go/test/setup"
)

const containerName = "go-client-postgres"

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
