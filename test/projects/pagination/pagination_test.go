package pagination

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
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
		name: "order by ASC",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOnePost(data: {
					id: "a",
					title: "a",
					content: "a",
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
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.Post.FindMany().OrderBy(
				Post.Title.Order(ASC),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []PostModel{{
				InnerPost: InnerPost{
					ID:      "a",
					Title:   "a",
					Content: "a",
				},
			}, {
				InnerPost: InnerPost{
					ID:      "b",
					Title:   "b",
					Content: "b",
				},
			}, {
				InnerPost: InnerPost{
					ID:      "c",
					Title:   "c",
					Content: "c",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "order by DESC",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOnePost(data: {
					id: "a",
					title: "a",
					content: "a",
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
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.Post.FindMany().OrderBy(
				Post.Title.Order(DESC),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []PostModel{{
				InnerPost: InnerPost{
					ID:      "c",
					Title:   "c",
					Content: "c",
				},
			}, {
				InnerPost: InnerPost{
					ID:      "b",
					Title:   "b",
					Content: "b",
				},
			}, {
				InnerPost: InnerPost{
					ID:      "a",
					Title:   "a",
					Content: "a",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "first 2",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOnePost(data: {
					id: "a",
					title: "a",
					content: "a",
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
					Post.Title.Order(ASC),
				).
				// would return a, b
				Take(2).
				Skip(1).
				// return records after b, which is c
				Cursor(Post.Title.Cursor("b")).
				Exec(ctx)

			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []PostModel{{
				InnerPost: InnerPost{
					ID:      "c",
					Title:   "c",
					Content: "c",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "first 2 skip",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOnePost(data: {
					id: "a",
					title: "a",
					content: "a",
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
					Post.Title.Order(ASC),
				).
				// would return a, b
				Take(2).
				// skip a, return b, c
				Skip(1).
				Exec(ctx)

			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []PostModel{{
				InnerPost: InnerPost{
					ID:      "b",
					Title:   "b",
					Content: "b",
				},
			}, {
				InnerPost: InnerPost{
					ID:      "c",
					Title:   "c",
					Content: "c",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "last 2",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOnePost(data: {
					id: "a",
					title: "a",
					content: "a",
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
					Post.Title.Order(ASC),
				).
				// would return b, c
				Take(-2).
				Skip(1).
				// before c will return b
				Cursor(Post.Title.Cursor("c")).
				Exec(ctx)

			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []PostModel{{
				InnerPost: InnerPost{
					ID:      "a",
					Title:   "a",
					Content: "a",
				},
			}, {
				InnerPost: InnerPost{
					ID:      "b",
					Title:   "b",
					Content: "b",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
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
					Post.IntTest.Order(DESC),
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
			test.RunSerial(t, []test.Database{test.SQLite, test.MySQL, test.PostgreSQL}, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
