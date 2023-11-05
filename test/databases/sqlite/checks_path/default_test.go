package db

import (
	"context"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
	"github.com/stretchr/testify/assert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestSqliteChecks(t *testing.T) {
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
					id: "checks-path-1",
					email: "asdf",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			assert.Equal(t, "file:dev.db", schemaDatasourceURL)

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
			client := NewClient()

			mockDB := test.Start(t, test.SQLite, client.Engine, tt.before)
			defer test.End(t, test.SQLite, client.Engine, mockDB)

			tt.run(t, client, context.Background())
		})
	}
}
