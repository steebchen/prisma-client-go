package relations

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
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

			expected := `{"id":"relations","email":"john@example.com","username":"johndoe","name":"John","roleID":null,"role":null,"posts":[{"id":"a","title":"common","content":"a","authorID":"relations","categoryID":null,"author":null,"Category":null,"comments":null},{"id":"b","title":"common","content":"b","authorID":"relations","categoryID":null,"author":null,"Category":null,"comments":null}],"comments":null}`
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
				InternalPost: InternalPost{
					ID:       "a",
					Title:    "common",
					Content:  str("a"),
					AuthorID: "relations",
				},
			}, {
				InternalPost: InternalPost{
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
				InternalUser: InternalUser{
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
				InternalPost: InternalPost{
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
				User.Posts.Fetch().Take(-2),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				InternalUser: InternalUser{
					ID:       "relations",
					Email:    "john@example.com",
					Username: "johndoe",
					Name:     str("John"),
				},
				RelationsUser: RelationsUser{
					Posts: []PostModel{{
						InternalPost: InternalPost{
							ID:       "c",
							Title:    "common",
							Content:  str("c"),
							AuthorID: "relations",
						},
					}, {
						InternalPost: InternalPost{
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
		name: "CreateOne with relation",
		// language=GraphQL
		before: []string{`
			mutation {
				user: createOneUser(data: {
					id: "relations",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.Post.CreateOne(
				Post.Title.Set("hi"),
				Post.Author.Link(
					User.ID.Equals("relations"),
				),
				Post.ID.Set("post1"),
			).With(
				Post.Author.Fetch(),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := PostModel{
				InternalPost: InternalPost{
					ID:       "post1",
					Title:    "hi",
					AuthorID: "relations",
				},
				RelationsPost: RelationsPost{
					Author: &UserModel{
						InternalUser: InternalUser{
							ID:       "relations",
							Email:    "john@example.com",
							Username: "johndoe",
							Name:     str("John"),
						},
					},
				},
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "unlink",
		// language=GraphQL
		before: []string{`
			mutation {
				user: createOneUser(data: {
					id: "relations",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
					role: {
						create: {
							id: "admin",
							name: "Admin",
						},
					},
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.User.FindOne(
				User.ID.Equals("relations"),
			).With(
				User.Role.Fetch(),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				InternalUser: InternalUser{
					ID:       "relations",
					Email:    "john@example.com",
					Username: "johndoe",
					Name:     str("John"),
					RoleID:   str("admin"),
				},
				RelationsUser: RelationsUser{
					Role: &RoleModel{
						InternalRole: InternalRole{
							ID:   "admin",
							Name: "Admin",
						},
					},
				},
			}

			assert.Equal(t, expected, actual)

			actual, err = client.User.FindOne(
				User.ID.Equals("relations"),
			).With(
				User.Role.Fetch(),
			).Update(
				User.Role.Unlink(),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expectedEmpty := UserModel{
				InternalUser: InternalUser{
					ID:       "relations",
					Email:    "john@example.com",
					Username: "johndoe",
					Name:     str("John"),
					RoleID:   nil,
				},
			}

			assert.Equal(t, expectedEmpty, actual)

			actual, err = client.User.FindOne(
				User.ID.Equals("relations"),
			).With(
				User.Role.Fetch(),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, expectedEmpty, actual)

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
				).Take(-2),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				InternalUser: InternalUser{
					ID:       "relations",
					Email:    "john@example.com",
					Username: "johndoe",
					Name:     str("John"),
				},
				RelationsUser: RelationsUser{
					Posts: []PostModel{{
						InternalPost: InternalPost{
							ID:       "b",
							Title:    "common",
							Content:  str("b"),
							AuthorID: "relations",
						},
					}, {
						InternalPost: InternalPost{
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
				).Take(-2).With(
					Post.Comments.Fetch().Take(2),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				InternalUser: InternalUser{
					ID:       "relations",
					Email:    "john@example.com",
					Username: "johndoe",
					Name:     str("John"),
				},
				RelationsUser: RelationsUser{
					Posts: []PostModel{{
						InternalPost: InternalPost{
							ID:       "b",
							Title:    "common",
							Content:  str("b"),
							AuthorID: "relations",
						},
						RelationsPost: RelationsPost{
							Comments: []CommentModel{},
						},
					}, {
						InternalPost: InternalPost{
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
				Post.Comments.Fetch().Take(-2),
				Post.Author.Fetch(),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			user := UserModel{
				InternalUser: InternalUser{
					ID:       "john",
					Email:    "john@example.com",
					Username: "johndoe",
					Name:     str("John"),
				},
			}

			author := actual.Author()

			assert.Equal(t, user, author)

			comments := []CommentModel{{
				InternalComment: InternalComment{
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
				InternalCategory: InternalCategory{
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
			test.RunSerial(t, []test.Database{test.SQLite, test.MySQL, test.PostgreSQL}, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
