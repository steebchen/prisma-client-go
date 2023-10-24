package sqlite

import (
	"fmt"
	"testing"

	"github.com/steebchen/prisma-client-go/test/cmd"
	"github.com/steebchen/prisma-client-go/test/setup"
)

var SQLite = &sqLite{}

const dbFileName = "%s.db"
const dbNameSchema = "file:" + dbFileName

type sqLite struct{}

func (*sqLite) Name() string {
	return "sqlite"
}

func (*sqLite) ConnectionString(mockDBName string) string {
	return fmt.Sprintf(dbNameSchema, mockDBName)
}

func (*sqLite) SetupDatabase(*testing.T) string {
	mockDB := setup.RandomString()
	return mockDB
}

func (*sqLite) TeardownDatabase(t *testing.T, name string) {
	if err := cmd.Run("rm", "-rf", fmt.Sprintf(dbFileName, name)); err != nil {
		t.Fatal(err)
	}
}
