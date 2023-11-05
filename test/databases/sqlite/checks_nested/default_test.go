package raw

import (
	"context"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
	db "github.com/steebchen/prisma-client-go/test/databases/sqlite/checks_nested/prisma"
	"github.com/stretchr/testify/assert"
)

type cx = context.Context
type Func func(t *testing.T, client *db.PrismaClient, ctx cx)

func TestSqliteChecksNested(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "check for connection URL",
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "checks-nested",
					email: "456",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *db.PrismaClient, ctx cx) {
			assert.Equal(t, "file:custom/dev.db", db.SchemaDatasourceURL)

			users, err := client.User.FindMany().Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, 1, len(users))
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			client := db.NewClient()

			mockDB := test.Start(t, test.SQLite, client.Engine, tt.before)
			defer test.End(t, test.SQLite, client.Engine, mockDB)

			tt.run(t, client, context.Background())
		})
	}
}
