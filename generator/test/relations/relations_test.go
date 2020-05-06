package relations

//go:generate go run github.com/prisma/prisma-client-go generate

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/generator/test/hooks"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func str(v string) *string {
	return &v
}

func TestRelations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "marshal json with relation",
		// language=GraphQL
		before: []string{`
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
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			user, err := client.User.FindOne(
				User.Email.Equals("john@example.com"),
			).With(
				User.Posts.Fetch(),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			actual, err := json.Marshal(&user)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := `{"id":"relations","email":"john@example.com","username":"johndoe","name":"John","posts":[{"id":"a","title":"common","content":"a","author":null,"authorID":"relations","comments":null},{"id":"b","title":"common","content":"b","author":null,"authorID":"relations","comments":null}],"comments":null}`
			assert.Equal(t, expected, string(actual))
		},
	}, {
		name: "find by single relation",
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
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
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
					ID:       "a",
					Title:    "common",
					Content:  str("a"),
					AuthorID: "relations",
				},
			}, {
				post{
					ID:       "b",
					Title:    "common",
					Content:  str("b"),
					AuthorID: "relations",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "find by to-many relation",
		// language=GraphQL
		before: []string{`
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
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
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
	}, {
		name: "create and connect",
		// language=GraphQL
		before: []string{`
			mutation {
				createOneUser(data: {
					id: "123",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			title := "What's up?"
			userID := "123"

			created, err := client.Post.CreateOne(
				Post.Title.Set(title),
				Post.Author.Link(
					User.ID.Equals(userID),
				),
				Post.ID.Set("post"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := PostModel{
				post{
					ID:       "post",
					Title:    title,
					AuthorID: "123",
				},
			}

			assert.Equal(t, expected, created)

			posts, err := client.Post.FindMany(
				Post.Title.Equals(title),
				Post.Author.Where(
					User.ID.Equals(userID),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, []PostModel{expected}, posts)
		},
	}, {
		name: "with simple",
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
						}, {
							id: "c",
							title: "common",
							content: "c",
						}, {
							id: "d",
							title: "stuff",
							content: "d",
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
			).With(
				User.Posts.Fetch().Last(2),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				user{
					ID:       "relations",
					Email:    "john@example.com",
					Username: "johndoe",
					Name:     str("John"),
					Posts: []PostModel{{
						post{
							ID:       "c",
							Title:    "common",
							Content:  str("c"),
							AuthorID: "relations",
						},
					}, {
						post{
							ID:       "d",
							Title:    "stuff",
							Content:  str("d"),
							AuthorID: "relations",
						},
					}},
					Comments: nil,
				},
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "with and sub query",
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
						}, {
							id: "c",
							title: "common",
							content: "c",
						}, {
							id: "d",
							title: "unrelated",
							content: "d",
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
			).With(
				User.Posts.Fetch(
					Post.Title.Equals("common"),
				).Last(2),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				user{
					ID:       "relations",
					Email:    "john@example.com",
					Username: "johndoe",
					Name:     str("John"),
					Posts: []PostModel{{
						post{
							ID:       "b",
							Title:    "common",
							Content:  str("b"),
							AuthorID: "relations",
						},
					}, {
						post{
							ID:       "c",
							Title:    "common",
							Content:  str("c"),
							AuthorID: "relations",
						},
					}},
					Comments: nil,
				},
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "with many to many nested",
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
						}, {
							id: "c",
							title: "common",
							content: "c",
						}, {
							id: "d",
							title: "unrelated",
							content: "d",
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
			).With(
				User.Posts.Fetch(
					Post.Title.Equals("common"),
				).Last(2).With(
					Post.Comments.Fetch().First(2),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				user{
					ID:       "relations",
					Email:    "john@example.com",
					Username: "johndoe",
					Name:     str("John"),
					Posts: []PostModel{{
						post{
							ID:       "b",
							Title:    "common",
							Content:  str("b"),
							AuthorID: "relations",
							Comments: []CommentModel{},
						},
					}, {
						post{
							ID:       "c",
							Title:    "common",
							Content:  str("c"),
							AuthorID: "relations",
							Comments: []CommentModel{},
						},
					}},
					Comments: nil,
				},
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "with by accessing methods",
		// language=GraphQL
		before: []string{`
			mutation {
				post: createOnePost(data: {
					id: "post-a",
					title: "common",
					content: "stuff",
					comments: {
						create: [{
							id: "comment-a",
							content: "this is a comment",
							by: {
								connect: {
									id: "john"
								},
							},
						}],
					},
					author: {
						create: {
							id: "john",
							email: "john@example.com",
							username: "johndoe",
							name: "John",
						},
					},
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.Post.FindOne(
				Post.ID.Equals("post-a"),
			).With(
				Post.Comments.Fetch().Last(2),
				Post.Author.Fetch(),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			user := UserModel{
				user{
					ID:       "john",
					Email:    "john@example.com",
					Username: "johndoe",
					Name:     str("John"),
				},
			}

			author, ok := actual.Author()

			assert.Equal(t, user, author)
			assert.Equal(t, true, ok)

			comments := []CommentModel{{
				comment{
					ID:      "comment-a",
					Content: "this is a comment",
					UserID:  "john",
					PostID:  "post-a",
				},
			}}

			assert.Equal(t, comments, actual.Comments())
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
