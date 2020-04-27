package pagination

//go:generate go run github.com/prisma/prisma-client-go generate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/generator/test/hooks"
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
				a: createOnePost(data: {
					id: "a",
					title: "a",
					content: "a",
				}) {
					id
				}
			}
		`, `
			mutation {
				c: createOnePost(data: {
					id: "c",
					title: "c",
					content: "c",
				}) {
					id
				}
			}
		`, `
			mutation {
				b: createOnePost(data: {
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
				post{
					ID:      "a",
					Title:   "a",
					Content: "a",
				},
			}, {
				post{
					ID:      "b",
					Title:   "b",
					Content: "b",
				},
			}, {
				post{
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
				a: createOnePost(data: {
					id: "a",
					title: "a",
					content: "a",
				}) {
					id
				}
			}
		`, `
			mutation {
				c: createOnePost(data: {
					id: "c",
					title: "c",
					content: "c",
				}) {
					id
				}
			}
		`, `
			mutation {
				b: createOnePost(data: {
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
				post{
					ID:      "c",
					Title:   "c",
					Content: "c",
				},
			}, {
				post{
					ID:      "b",
					Title:   "b",
					Content: "b",
				},
			}, {
				post{
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
				a: createOnePost(data: {
					id: "a",
					title: "a",
					content: "a",
				}) {
					id
				}
			}
		`, `
			mutation {
				c: createOnePost(data: {
					id: "c",
					title: "c",
					content: "c",
				}) {
					id
				}
			}
		`, `
			mutation {
				b: createOnePost(data: {
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
				First(2).
				// return records after b, which is c
				After(Post.Title.Cursor("b")).
				Exec(ctx)

			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []PostModel{{
				post{
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
				a: createOnePost(data: {
					id: "a",
					title: "a",
					content: "a",
				}) {
					id
				}
			}
		`, `
			mutation {
				c: createOnePost(data: {
					id: "c",
					title: "c",
					content: "c",
				}) {
					id
				}
			}
		`, `
			mutation {
				b: createOnePost(data: {
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
				First(2).
				// skip a, return b, c
				Skip(1).
				Exec(ctx)

			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []PostModel{{
				post{
					ID:      "b",
					Title:   "b",
					Content: "b",
				},
			}, {
				post{
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
				a: createOnePost(data: {
					id: "a",
					title: "a",
					content: "a",
				}) {
					id
				}
			}
		`, `
			mutation {
				c: createOnePost(data: {
					id: "c",
					title: "c",
					content: "c",
				}) {
					id
				}
			}
		`, `
			mutation {
				b: createOnePost(data: {
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
				Last(2).
				// before c will return b
				Before(Post.Title.Cursor("c")).
				Exec(ctx)

			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []PostModel{{
				post{
					ID:      "a",
					Title:   "a",
					Content: "a",
				},
			}, {
				post{
					ID:      "b",
					Title:   "b",
					Content: "b",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			hooks.Start(t, client.Engine, tt.before)
			defer hooks.End(t, client.Engine)
			tt.run(t, client, context.Background())
		})
	}
}
