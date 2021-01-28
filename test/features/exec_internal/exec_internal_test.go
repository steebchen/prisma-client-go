package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestExecInner(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "find unique",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "123",
					username: "johndoe",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.User.FindUnique(
				User.ID.Equals("123"),
			).ExecInner(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &InnerUser{
				ID:       "123",
				Username: "johndoe",
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "find many",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "123",
					username: "johndoe",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.User.FindMany().ExecInner(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []InnerUser{{
				ID:       "123",
				Username: "johndoe",
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "find first",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "123",
					username: "johndoe",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.User.FindFirst().ExecInner(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &InnerUser{
				ID:       "123",
				Username: "johndoe",
			}

			assert.Equal(t, expected, actual)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, []test.Database{test.SQLite, test.MySQL, test.PostgreSQL}, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
