package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/prisma/prisma-client-go/cli"
	"github.com/prisma/prisma-client-go/engine"
	"github.com/prisma/prisma-client-go/test/cmd"
	"github.com/prisma/prisma-client-go/test/setup/mysql"
	"github.com/prisma/prisma-client-go/test/setup/postgresql"
	"github.com/prisma/prisma-client-go/test/setup/sqlite"
)

type Database interface {
	Name() string
	ConnectionString(mockDBName string) string
	Setup()
	Teardown()
	SetupDatabase(t *testing.T) string
	TeardownDatabase(t *testing.T, mockDBName string)
}

type Engine interface {
	Connect() error
	Disconnect() error
	Do(context.Context, string, interface{}) error
}

var MySQL = mysql.MySQL
var PostgreSQL = postgresql.PostgreSQL
var SQLite = sqlite.SQLite

var Databases = []Database{
	mysql.MySQL,
	postgresql.PostgreSQL,
	sqlite.SQLite,
}

const schemaTemplate = "schema.temp.%s.prisma"

func replaceSchema(t *testing.T, db Database, e *engine.Engine, schemaPath string, mockDB string) {
	e.ReplaceSchema(func(schema string) string {
		for _, fromDB := range Databases {
			schema = strings.ReplaceAll(schema, fmt.Sprintf(`"%s"`, fromDB.Name()), fmt.Sprintf(`"%s"`, db.Name()))
		}
		return schema
	})
	e.ReplaceSchema(func(schema string) string {
		return strings.ReplaceAll(schema, `env("__REPLACE__")`, fmt.Sprintf(`"%s"`, db.ConnectionString(mockDB)))
	})
	if err := ioutil.WriteFile(schemaPath, []byte(e.Schema), 0644); err != nil {
		t.Fatal(err)
	}
}

func Start(t *testing.T, db Database, e *engine.Engine, queries []string) string {
	mockDB := db.SetupDatabase(t)

	schemaPath := fmt.Sprintf(schemaTemplate, db.Name())
	replaceSchema(t, db, e, schemaPath, mockDB)

	migrate(t, schemaPath)

	if err := e.Connect(); err != nil {
		t.Fatalf("could not connect: %s", err)
		return ""
	}

	for _, b := range queries {
		var response engine.GQLResponse
		if err := e.Do(context.Background(), b, &response); err != nil {
			End(t, db, e, mockDB)
			t.Fatalf("could not send mock query %s", err)
		}
		if response.Errors != nil {
			End(t, db, e, mockDB)
			t.Fatalf("mock query has errors %+v", response)
		}
	}

	log.Printf("")
	log.Printf("---")
	log.Printf("")

	return mockDB
}

func End(t *testing.T, db Database, e Engine, mockDBName string) {
	defer teardown(t, db, mockDBName)

	if err := e.Disconnect(); err != nil {
		t.Fatalf("could not disconnect: %s", err)
	}
}

func teardown(t *testing.T, db Database, mockDBName string) {
	if err := cmd.Run("rm", "-rf", fmt.Sprintf(schemaTemplate, db.Name())); err != nil {
		t.Fatal(err)
	}

	db.TeardownDatabase(t, mockDBName)

	cleanup(t)
}

func RunSerial(t *testing.T, dbs []Database, invoke func(t *testing.T, db Database, ctx context.Context)) {
	run(t, dbs, invoke)
}

func RunParallel(t *testing.T, dbs []Database, invoke func(t *testing.T, db Database, ctx context.Context)) {
	t.Parallel()

	run(t, dbs, invoke)
}

func run(t *testing.T, dbs []Database, invoke func(t *testing.T, db Database, ctx context.Context)) {
	for _, db := range dbs {
		db := db
		t.Run(db.Name(), func(t *testing.T) {
			invoke(t, db, context.Background())
		})
	}
}

func migrate(t *testing.T, schemaPath string) {
	cleanup(t)

	verbose := os.Getenv("PRISMA_CLIENT_GO_TEST_MIGRATE_LOGS") != ""
	if err := cli.Run([]string{"migrate", "save", "--experimental", "--schema=./" + schemaPath, "--create-db", "--name", "init"}, verbose); err != nil {
		t.Fatalf("could not run migrate save --experimental %s", err)
	}

	if err := cli.Run([]string{"migrate", "up", "--experimental", "--schema=./" + schemaPath}, verbose); err != nil {
		t.Fatalf("could not run migrate save --experimental %s", err)
	}
}

func cleanup(t *testing.T) {
	if err := cmd.Run("rm", "-rf", "migrations"); err != nil {
		t.Fatal(err)
	}
}
