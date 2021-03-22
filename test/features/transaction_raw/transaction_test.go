package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/runtime/types"
	"github.com/prisma/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestTransaction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "query raw",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "123",
					email: "john@example.com",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			e := client.Prisma.QueryRaw(`SELECT * FROM "User"`).Tx()

			if err := client.Prisma.Transaction(e).Exec(ctx); err != nil {
				t.Fatal(err)
			}

			var v []UserModel
			if err := e.Into(&v); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, []UserModel{{
				InnerUser: InnerUser{
					ID:    "123",
					Email: "john@example.com",
				},
				RelationsUser: RelationsUser{},
			}}, v)
		},
	}, {
		name: "execute raw",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "123",
					email: "john@example.com",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			e := client.Prisma.ExecuteRaw(`UPDATE "User" SET email = $1 WHERE id = $2`, "new-email", "123").Tx()

			if err := client.Prisma.Transaction(e).Exec(ctx); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, &types.BatchResult{
				Count: 1,
			}, e.Result())

			// --

			actual, err := client.User.FindMany().Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			expected := []UserModel{{
				InnerUser: InnerUser{
					ID:    "123",
					Email: "new-email",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, []test.Database{test.PostgreSQL}, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
