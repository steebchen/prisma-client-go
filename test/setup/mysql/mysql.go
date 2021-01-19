package mysql

import (
	"fmt"
	"log"
	"testing"

	"github.com/prisma/prisma-client-go/test/cmd"
	"github.com/prisma/prisma-client-go/test/setup"
)

const containerName = "go-client-mysql"
const image = "mysql:5.7"

var MySQL = &mySQL{}

type mySQL struct{}

func (*mySQL) Name() string {
	return "mysql"
}

func (*mySQL) ConnectionString(mockDBName string) string {
	return fmt.Sprintf("mysql://root:pw@localhost:3307/%s", mockDBName)
}

func (*mySQL) SetupDatabase(t *testing.T) string {
	mockDB := setup.RandomString()

	exec(t, fmt.Sprintf("CREATE DATABASE %s", mockDB))

	return mockDB
}

func (*mySQL) TeardownDatabase(t *testing.T, mockDB string) {
	exec(t, fmt.Sprintf("DROP DATABASE %s", mockDB))
}

func exec(t *testing.T, query string) {
	if err := cmd.Run("docker", "exec", "-t", containerName, "mysql", "--user=root", "--password=pw", fmt.Sprintf("--execute=%s", query)); err != nil {
		t.Fatal(err)
	}
}

func (*mySQL) Setup() {
	if err := cmd.Run("docker", "run", "--name", containerName, "-p", "3307:3306", "-e", "MYSQL_DATABASE=testing", "-e", "MYSQL_ROOT_PASSWORD=pw", "-d", image); err != nil {
		panic(err)
	}
}

func (*mySQL) Teardown() {
	if err := cmd.Run("docker", "stop", containerName); err != nil {
		log.Println(err)
	}

	if err := cmd.Run("docker", "rm", containerName, "--force"); err != nil {
		log.Println(err)
	}
}
