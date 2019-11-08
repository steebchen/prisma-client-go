package relations

//go:generate prisma2 generate

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/photongo/generator/test/hooks"
)

type cx = context.Context
type Func func(t *testing.T, client Client, ctx cx)

func str(v string) *string {
	return &v
}

func TestRelations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before string
		run    Func
	}{{
		name: "find many posts by user",
		// language=GraphQL
		before: `
			mutation {
				unrelated: createOnePost(data: {
					id: "nope",
					title: "nope",
					content: "nope",
					user: {
						create: {
							id: "unrelated",
							email: "unrelated",
							username: "unrelated",
							name: "unrelated",
						}
					}
				}) {
					id
				}

				user: createOneUser(data: {
					id: "relations",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
					posts: {
						create: [{
							id: "a",
							title: "common",
							content: "a",
						}, {
							id: "b",
							title: "common",
							content: "b",
						}],
					},
				}) {
					id
				}
			}
		`,
		run: func(t *testing.T, client Client, ctx cx) {
			actual, err := client.Post.FindMany(
				Post.Title.Equals("common"),
				Post.User().Email.Equals("john@example.com"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []PostModel{{
				post{
					ID:      "a",
					Title:   "common",
					Content: str("a"),
				},
			}, {
				post{
					ID:      "b",
					Title:   "common",
					Content: str("b"),
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "find many posts by user",
		// language=GraphQL
		before: `
			mutation {
				unrelated: createOnePost(data: {
					id: "nope",
					title: "nope",
					content: "nope",
					user: {
						create: {
							id: "unrelated",
							email: "unrelated",
							username: "unrelated",
							name: "unrelated",
						}
					}
				}) {
					id
				}

				user: createOneUser(data: {
					id: "relations",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
					posts: {
						create: [{
							id: "a",
							title: "common",
							content: "a",
						}, {
							id: "b",
							title: "common",
							content: "b",
						}],
					},
				}) {
					id
				}
			}
		`,
		run: func(t *testing.T, client Client, ctx cx) {
			actual, err := client.User.FindMany(
				User.Email.Equals("john@example.com"),
				// Post.Content.Equals("f"), // should fail
				User.Posts().Title.Equals("common"),
				User.Posts().User().Email.Equals("common"),
				User.Posts().User().Posts().Content.Equals("common"),
				User.Posts().User().Posts().User().Email.Equals("common"),
				User.Comments().User().Posts().Title.Equals("common"),
				User.Comments().User().Posts().Title.Equals("common"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			p, err := client.Post.FindMany(
				// User.Email.Equals("john@example.com"), // should fail
				Post.Title.Equals("common"),
				Post.User().Email.Equals("common"),
				Post.User().Posts().Content.Equals("common"),
				Post.User().Posts().User().Email.Equals("common"),
				Post.Comments().Post().Title.Equals("common"),
				Post.Comments().Post().User().Email.Equals("common"),
				Post.Comments().Post().Comments().Content.Equals("common"),
				Post.Comments().User().Posts().Title.Equals("common"),
				Post.Comments().User().Posts().Title.Equals("common"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			Post.User()
			Post.Comments()

			User.Posts()
			User.Comments()

			Comment.User()
			Comment.Post()

			Post.User().Email.Equals("f")
			Post.Comments().Content.Equals("f")
			Post.User(). /*Comments()*/ Email.Equals("f")
			Post.User().Posts().Title.Equals("f")
			Post.User().Posts().User().Email.Equals("f")
			Post.User().Posts().User().Posts().User().Email.Equals("f")
			Post.User().Posts().User().Posts().Content.Equals("f")

			Post.User().Posts().Comments()
			Post.Comments().Post()

			Post.Comments().Content.Equals("f")
			Post.Comments().Post().Content.Equals("f")
			Post.Comments().Post().Comments().Post().Comments().Post().Comments().Content.Equals("f")

			expected := []PostModel{{
				post{
					ID:      "a",
					Title:   "common",
					Content: str("a"),
				},
			}, {
				post{
					ID:      "b",
					Title:   "common",
					Content: str("b"),
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			hooks.Run(t)

			client := NewClient()
			if err := client.Connect(); err != nil {
				t.Fatalf("could not connect %s", err)
				return
			}

			defer func() {
				_ = client.Disconnect()
				// TODO blocked by prisma-engine panicking on disconnect
				// if err != nil {
				// 	t.Fatalf("could not disconnect %s", err)
				// }
			}()

			ctx := context.Background()

			if tt.before != "" {
				response, err := client.gql.Raw(ctx, tt.before, map[string]interface{}{})
				log.Printf("mock response query %+v", response)
				if err != nil {
					t.Fatalf("could not send mock query %s %+v", err, response)
				}
			}

			tt.run(t, client, ctx)
		})
	}
}
