package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

func TestTableCasing(t *testing.T) {
	test.RunParallel(t, []test.Database{test.MySQL, test.PostgreSQL, test.SQLite}, func(t *testing.T, db test.Database, ctx context.Context) {
		client := NewClient()

		// language=GraphQL
		mockDB := test.Start(t, db, client.Engine, []string{`
			mutation {
				result: createOnePost(data: {
					id: "a",
					int: 10,
					float: 10,
					int2: 10,
					float2: 10,
				}) {
					id
				}
			}
		`})
		defer test.End(t, db, client.Engine, mockDB)

		expectedPost := &PostModel{
			InnerPost: InnerPost{
				ID:     "a",
				Int:    13,
				Float:  7.5,
				Int2:   20,
				Float2: 5,
			},
		}

		actualFoundPost, err := client.Post.FindUnique(
			Post.ID.Equals("a"),
		).Update(
			Post.Int.Increment(3),
			Post.Float.Decrement(2.5),
			Post.Int2.Multiply(2),
			Post.Float2.Divide(2),
		).Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expectedPost, actualFoundPost)

		actualUpdatedPost, err := client.Post.FindUnique(
			Post.ID.Equals("a"),
		).Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expectedPost, actualUpdatedPost)
	})
}
