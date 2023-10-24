package raw

import (
	"context"
	checks_nested_db "github.com/steebchen/prisma-client-go/test/databases/sqlite/checks_nested/prisma"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *checks_nested_db.PrismaClient, ctx cx)

func TestSqliteChecksNested(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "check for connection URL",
		run: func(t *testing.T, client *checks_nested_db.PrismaClient, ctx cx) {
			assert.Equal(t, "sqlite:prisma/dev.db", checks_nested_db.SchemaConnectionURL)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			client := checks_nested_db.NewClient()

			mockDB := test.Start(t, test.SQLite, client.Engine, tt.before)
			defer test.End(t, test.SQLite, client.Engine, mockDB)

			tt.run(t, client, context.Background())
		})
	}
}
