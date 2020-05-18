package relations

//go:generate go run github.com/prisma/prisma-client-go generate

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test/hooks"
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

			expected := `{"id":"relations","email":"john@example.com","username":"johndoe","name":"John","posts":[{"id":"a","title":"common","content":"a","authorID":"relations","categoryID":null,"author":null,"category":null,"comments":null},{"id":"b","title":"common","content":"b","authorID":"relations","categoryID":null,"author":null,"category":null,"comments":null}],"comments":null}`
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
			).OrderBy(Post.ID.Order(ASC)).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []PostModel{{
				RawPost: RawPost{
					ID:       "a",
					Title:    "common",
					Content:  str("a"),
					AuthorID: "relations",
				},
			}, {
				RawPost: RawPost{
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
				RawUser: RawUser{
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
				RawPost: RawPost{
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
				RawUser: RawUser{
					ID:       "relations",
					Email:    "john@example.com",
					Username: "johndoe",
					Name:     str("John"),
				},
				RelationsUser: RelationsUser{
					Posts: []PostModel{{
						RawPost: RawPost{
							ID:       "c",
							Title:    "common",
							Content:  str("c"),
							AuthorID: "relations",
						},
					}, {
						RawPost: RawPost{
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
				RawUser: RawUser{
					ID:       "relations",
					Email:    "john@example.com",
					Username: "johndoe",
					Name:     str("John"),
				},
				RelationsUser: RelationsUser{
					Posts: []PostModel{{
						RawPost: RawPost{
							ID:       "b",
							Title:    "common",
							Content:  str("b"),
							AuthorID: "relations",
						},
					}, {
						RawPost: RawPost{
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
				RawUser: RawUser{
					ID:       "relations",
					Email:    "john@example.com",
					Username: "johndoe",
					Name:     str("John"),
				},
				RelationsUser: RelationsUser{
					Posts: []PostModel{{
						RawPost: RawPost{
							ID:       "b",
							Title:    "common",
							Content:  str("b"),
							AuthorID: "relations",
						},
						RelationsPost: RelationsPost{
							Comments: []CommentModel{},
						},
					}, {
						RawPost: RawPost{
							ID:       "c",
							Title:    "common",
							Content:  str("c"),
							AuthorID: "relations",
						},
						RelationsPost: RelationsPost{
							Comments: []CommentModel{},
						},
					}},
					Comments: nil,
				},
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "with by accessing methods with required relation",
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
				RawUser: RawUser{
					ID:       "john",
					Email:    "john@example.com",
					Username: "johndoe",
					Name:     str("John"),
				},
			}

			author := actual.Author()

			assert.Equal(t, user, author)

			comments := []CommentModel{{
				RawComment: RawComment{
					ID:      "comment-a",
					Content: "this is a comment",
					UserID:  "john",
					PostID:  "post-a",
				},
			}}

			assert.Equal(t, comments, actual.Comments())
		},
	}, {
		name: "create and find with existing optional relation",
		// language=GraphQL
		before: []string{`
			mutation {
				post: createOneUser(data: {
					id: "john",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
				}) {
					id
				}
			}
		`, `
			mutation {
				category: createOneCategory(data: {
					id: "media",
					name: "Media",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			title := "What's up?"
			userID := "john"

			_, err := client.Post.CreateOne(
				Post.Title.Set(title),
				Post.Author.Link(
					User.ID.Equals(userID),
				),
				Post.ID.Set("post"),
				Post.Category.Link(
					Category.ID.Equals("media"),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			actual, err := client.Post.FindOne(
				Post.ID.Equals("post"),
			).With(
				Post.Category.Fetch(),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expectedCategory := CategoryModel{
				RawCategory: RawCategory{
					ID:   "media",
					Name: "Media",
				},
			}

			actualCategory, ok := actual.Category()

			assert.Equal(t, expectedCategory, actualCategory)
			assert.Equal(t, true, ok)
		},
	}, {
		name: "create and find with non-existing optional relations",
		// language=GraphQL
		before: []string{`
			mutation {
				post: createOneUser(data: {
					id: "john",
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
			userID := "john"

			_, err := client.Post.CreateOne(
				Post.Title.Set(title),
				Post.Author.Link(
					User.ID.Equals(userID),
				),
				Post.ID.Set("post"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			actual, err := client.Post.FindOne(
				Post.ID.Equals("post"),
			).With(
				Post.Category.Fetch(),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			actualCategory, ok := actual.Category()

			assert.Equal(t, CategoryModel{}, actualCategory)
			assert.Equal(t, false, ok)
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
