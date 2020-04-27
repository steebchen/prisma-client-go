package chains

//go:generate go run github.com/prisma/prisma-client-go generate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/generator/test/hooks"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func str(v string) *string {
	return &v
}

func TestRelationChains(t *testing.T) {
	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "return many",
		// language=GraphQL
		before: []string{`
			mutation {
				unrelated: createOnePost(data: {
					id: "same-relevant-for-test",
					title: "common",
					content: "a",
					author: {
						create: {
							id: "x",
							email: "x",
							username: "x",
							name: "x",
						}
					}
				}) {
					id
				}
			}
		`, `
			mutation {
				user: createOneUser(data: {
					id: "john",
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
						}, {
							id: "c",
							title: "stuff",
							content: "c",
						}],
					},
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.User.FindOne(
				User.Email.Equals("john@example.com"),
			).GetPosts(
				Post.Title.Equals("common"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []PostModel{{
				post{
					ID:       "a",
					Title:    "common",
					Content:  str("a"),
					AuthorID: "john",
				},
			}, {
				post{
					ID:       "b",
					Title:    "common",
					Content:  str("b"),
					AuthorID: "john",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "return one",
		// language=GraphQL
		before: []string{`
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
			}
		`, `
			mutation {
				user: createOneUser(data: {
					id: "john",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
					posts: {
						create: [{
							id: "a",
							title: "common",
							content: "a",
						}],
					},
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.Post.FindOne(
				Post.ID.Equals("a"),
			).GetAuthor().Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				user{
					ID:       "a",
					Email:    "john@example.com",
					Username: "johndoe",
					Name:     str("John"),
				},
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "select one and return many",
		// language=GraphQL
		before: []string{`
			mutation {
				unrelated: createOnePost(data: {
					id: "unrelated",
					title: "unrelated",
					content: "unrelated",
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
			}
		`, `
			mutation {
				user: createOneUser(data: {
					id: "john",
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
						}, {
							id: "c",
							title: "stuff",
							content: "c",
						}],
					},
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.Post.FindOne(
				Post.ID.Equals("a"),
			).GetAuthor().GetPosts().Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []PostModel{{
				post{
					ID:       "a",
					Title:    "common",
					Content:  str("a"),
					AuthorID: "john",
				},
			}, {
				post{
					ID:       "b",
					Title:    "common",
					Content:  str("b"),
					AuthorID: "john",
				},
			}, {
				post{
					ID:       "c",
					Title:    "c",
					Content:  str("c"),
					AuthorID: "c",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "select one and return many with parameters",
		// language=GraphQL
		before: []string{`
			mutation {
				unrelated: createOnePost(data: {
					id: "unrelated",
					title: "unrelated",
					content: "unrelated",
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
			}
		`, `
			mutation {
				user: createOneUser(data: {
					id: "john",
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
						}, {
							id: "c",
							title: "stuff",
							content: "c",
						}],
					},
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.Post.FindOne(
				Post.ID.Equals("c"),
			).GetAuthor().GetPosts(
				Post.Title.Equals("common"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []PostModel{{
				post{
					ID:       "a",
					Title:    "common",
					Content:  str("a"),
					AuthorID: "john",
				},
			}, {
				post{
					ID:       "b",
					Title:    "common",
					Content:  str("b"),
					AuthorID: "john",
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
