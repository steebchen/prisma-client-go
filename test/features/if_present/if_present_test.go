package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func str(v string) *string {
	return &v
}

func TestIfPresent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "find and update",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "update",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			maybeQuery := "johndoe"
			newUsername := "new-username"
			_, err := client.User.FindMany(
				User.Email.Equals("john@example.com"),
				// query for this one
				User.Username.EqualsIfPresent(&maybeQuery),
				// ignore this one
				User.Name.EqualsIfPresent(nil),
			).Update(
				// set value
				User.Username.SetIfPresent(&newUsername),
				// don't set because nil
				User.Name.SetIfPresent(nil),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			updated, err := client.User.FindUnique(
				User.ID.Equals("update"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &UserModel{
				InnerUser: InnerUser{
					ID:       "update",
					Email:    "john@example.com",
					Username: "new-username",
					Name:     str("John"),
				},
			}

			assert.Equal(t, expected, updated)

			actual, err := client.User.FindUnique(User.ID.Equals("update")).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "update operations",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "update",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
					age: 1,
					age2: 2,
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			incrementAge := 50
			var newAge2 int
			_, err := client.User.FindUnique(
				User.ID.Equals("update"),
			).Update(
				// set value
				User.Age.IncrementIfPresent(&incrementAge),
				// don't set because nil
				User.Age2.IncrementIfPresent(&newAge2),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			updated, err := client.User.FindUnique(
				User.ID.Equals("update"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			age := 51
			age2 := 2
			expected := &UserModel{
				InnerUser: InnerUser{
					ID:       "update",
					Email:    "john@example.com",
					Username: "johndoe",
					Name:     str("John"),
					Age:      &age,
					Age2:     &age2,
				},
			}

			assert.Equal(t, expected, updated)

			actual, err := client.User.FindUnique(User.ID.Equals("update")).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "with link filled",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "test",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			v := "test"
			post, err := client.Post.CreateOne(
				Post.Title.Set("asdf"),
				Post.Author.Link(
					User.ID.EqualsIfPresent(&v),
				),
				Post.ID.Set("post-1"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &PostModel{
				InnerPost: InnerPost{
					ID:       "post-1",
					Title:    "asdf",
					AuthorID: str("test"),
				},
			}

			assert.Equal(t, expected, post)
		},
	}, {
		name: "with link nil",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "test",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var v *string = nil
			post, err := client.Post.CreateOne(
				Post.Title.Set("asdf"),
				Post.Author.Link(
					User.ID.EqualsIfPresent(v),
				),
				Post.ID.Set("post-1"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &PostModel{
				InnerPost: InnerPost{
					ID:       "post-1",
					Title:    "asdf",
					AuthorID: nil,
				},
			}

			assert.Equal(t, expected, post)
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
