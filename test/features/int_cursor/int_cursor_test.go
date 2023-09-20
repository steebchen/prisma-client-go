package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/steebchen/prisma-client-go/test"
	"github.com/steebchen/prisma-client-go/test/setup/mysql"
	"github.com/steebchen/prisma-client-go/test/setup/postgresql"
	"github.com/steebchen/prisma-client-go/test/setup/sqlite"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestPagination(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "int cursor",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOnePost(data: {
					id: "a",
					title: "a",
					content: "a",
					intTest: 3,
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOnePost(data: {
					id: "c",
					title: "c",
					content: "c",
					intTest: 1,
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOnePost(data: {
					id: "b",
					title: "b",
					content: "b",
					intTest: 2,
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.
				Post.
				FindMany().
				OrderBy(
					Post.IntTest.Order(SortOrderDesc),
				).
				Cursor(Post.IntTest.Cursor(2)).
				Exec(ctx)

			if err != nil {
				t.Fatalf("fail %s", err)
			}

			b := 2
			c := 1
			expected := []PostModel{{
				InnerPost: InnerPost{
					ID:      "b",
					Title:   "b",
					Content: "b",
					IntTest: &b,
				},
			}, {
				InnerPost: InnerPost{
					ID:      "c",
					Title:   "c",
					Content: "c",
					IntTest: &c,
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, []test.Database{mysql.MySQL, postgresql.PostgreSQL, sqlite.SQLite}, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
