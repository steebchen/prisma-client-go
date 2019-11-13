package relations

//go:generate prisma2 generate

import (
	"context"
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
		name: "find by single relation",
		// language=GraphQL
		before: `
			mutation {
				unrelated: createOnePost(data: {
					id: "nope",
					title: "nope",
					content: "nope",
					author: {
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
				Post.Author.Where(
					User.Email.Equals("john@example.com"),
				),
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
		name: "find by to-many relation",
		// language=GraphQL
		before: `
			mutation {
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
							comments: {
								create: [{
								id: "comment1",
								content: "comment 1",
								by: {
									connect: {
										id: "relations"
									}
								}
							}]
							}
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
				User.Posts.Some(
					Post.Title.Equals("common"),
					Post.Comments.Every(
						Comment.Content.Contains("comment"),
					),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				user{
					ID:       "relations",
					Email:    "john@example.com",
					Username: "johndoe",
					Name:     str("John"),
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
				var response gqlResponse
				err := client.do(ctx, tt.before, &response)
				if err != nil {
					t.Fatalf("could not send mock query %s", err)
				}
				if response.Errors != nil {
					t.Fatalf("mock query has errors %+v", response)
				}
			}

			tt.run(t, client, ctx)
		})
	}
}
