package load

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

func TestLoad(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "multiple things",
		// language=GraphQL
		before: []string{`
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
						}],
					},
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			type Result struct {
				FindOneUser  UserModel   `json:"findOneUser"`
				FindManyUser []UserModel `json:"findManyUser"`
				FindManyPost []PostModel `json:"findManyPost"`
			}
			var actual Result
			err := client.Load(
				client.User.FindOne(
					User.ID.Equals("john"),
				).Load(),
				// TODO currently only one query is supported
				// client.User.FindMany(User.Name.Equals("John")).Load(),
				// client.Post.FindMany(Post.Title.Equals("common")).Load(),
			).Exec(ctx, &actual)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := Result{
				FindOneUser: UserModel{
					user{
						ID:       "john",
						Email:    "john@example.com",
						Username: "johndoe",
						Name:     str("John"),
					},
				},
				// FindManyUser: []UserModel{
				// 	{
				// 		user{
				// 			ID:       "john",
				// 			Email:    "john@example.com",
				// 			Username: "johndoe",
				// 			Name:     str("John"),
				// 		},
				// 	},
				// },
				// FindManyPost: []PostModel{
				// 	{
				// 		post{
				// 			ID:       "a",
				// 			Title:    "common",
				// 			Content:  str("a"),
				// 			AuthorID: "john",
				// 		},
				// 	},
				// 	{
				// 		post{
				// 			ID:       "b",
				// 			Title:    "common",
				// 			Content:  str("b"),
				// 			AuthorID: "john",
				// 		},
				// 	},
				// },
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "fetch a `many` relation",
		// language=GraphQL
		before: []string{`
			mutation {
				user: createOneUser(data: {
					id: "john",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
					posts: {
						create: [{
							id: "a",
							title: "1",
							content: "x",
						}, {
							id: "b",
							title: "1",
							content: "x",
						}, {
							id: "c",
							title: "2",
							content: "stuff",
						}, {
							id: "d",
							title: "2",
							content: "stuff",
						}, {
							id: "e",
							title: "2",
							content: "non-stuff",
						}],
					},
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			type UserResponse struct {
				UserModel
				Posts []PostModel `json:"posts"`
			}
			type Result struct {
				FindManyUser []UserResponse `json:"findManyUser"`
			}
			var actual Result
			err := client.Load(
				client.User.FindMany(
					User.ID.Equals("john"),
					// query for some users but filter by posts
					User.Posts.Some(
						Post.Title.Equals("2"),
					),

					// for those users which were found, additionally fetch posts
					User.Posts.
						Fetch(
							Post.Content.Equals("stuff"),
						).
						// returns last of matching content "stuff", meaning id d will be returned instead of c
						Last(1).
						Load(),
				).Load(),
			).Exec(ctx, &actual)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := Result{
				FindManyUser: []UserResponse{{
					UserModel: UserModel{
						user{
							ID:       "john",
							Email:    "john@example.com",
							Username: "johndoe",
							Name:     str("John"),
						},
					},
					Posts: []PostModel{{
						post{
							ID:       "d",
							Title:    "2",
							Content:  str("stuff"),
							AuthorID: "john",
						},
					}},
				}},
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "fetch a `one` relation",
		// language=GraphQL
		before: []string{`
			mutation {
				user: createOneUser(data: {
					id: "john",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
					posts: {
						create: [{
							id: "a",
							title: "1",
							content: "x",
						}, {
							id: "b",
							title: "2",
							content: "x",
						}],
					},
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			type PostResponse struct {
				PostModel
				Author UserModel `json:"author"`
			}
			type Result struct {
				FindManyPost []PostResponse `json:"findManyPost"`
			}
			var actual Result
			err := client.Load(
				client.Post.FindMany(
					Post.Title.Equals("1"),
					Post.Author.Fetch().Load(),
				).Load(),
			).Exec(ctx, &actual)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := Result{
				FindManyPost: []PostResponse{{
					PostModel: PostModel{
						post{
							ID:       "a",
							Title:    "1",
							Content:  str("x"),
							AuthorID: "john",
						},
					},
					Author: UserModel{
						user{
							ID:       "john",
							Email:    "john@example.com",
							Username: "johndoe",
							Name:     str("John"),
						},
					},
				}},
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "deeply nested relation",
		// language=GraphQL
		before: []string{`
			mutation {
				user: createOneUser(data: {
					id: "john",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
					posts: {
						create: [{
							id: "post-a",
							title: "post a",
							content: "post a",
							comments: {
								create: [{
									id: "comment-a-1",
									content: "a 1",
									by: {
										connect: {
											id: "john",
										},
									},
								}, {
									id: "comment-a-2",
									content: "a 2",
									by: {
										connect: {
											id: "john",
										},
									},
								}],
							},
						}, {
							id: "post-b",
							title: "post b",
							content: "post b",
						}],
					},
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			type PostResponse struct {
				PostModel
				Comments []CommentModel `json:"comments"`
			}
			type UserResponse struct {
				UserModel
				Posts []PostResponse `json:"posts"`
			}
			type Result struct {
				FindManyUser []UserResponse `json:"findManyUser"`
			}
			var actual Result
			err := client.Load(
				client.User.FindMany(
					User.ID.Equals("john"),

					// for those users which were found, additionally fetch posts
					User.Posts.
						Fetch(
							Post.Content.Equals("post a"),

							// for each post, also load all comments
							Post.Comments.Fetch(
								Comment.Content.Equals("a 2"),
							).Load(),
						).
						// returns last of matching content "stuff", meaning id d will be returned instead of c
						Last(1).
						Load(),
				).Load(),
			).Exec(ctx, &actual)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := Result{
				FindManyUser: []UserResponse{{
					UserModel: UserModel{
						user{
							ID:       "john",
							Email:    "john@example.com",
							Username: "johndoe",
							Name:     str("John"),
						},
					},
					Posts: []PostResponse{{
						PostModel: PostModel{
							post{
								ID:       "post-a",
								Title:    "post a",
								Content:  str("post a"),
								AuthorID: "john",
							},
						},
						Comments: []CommentModel{
							{
								comment{
									ID:      "comment-a-2",
									Content: "a 2",
									UserID:  "john",
									PostID:  "post-a",
								},
							},
						},
					}},
				}},
			}

			assert.Equal(t, expected, actual)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			hooks.Start(t, client.Engine, tt.before)
			tt.run(t, client, context.Background())
			hooks.End(t, client.Engine)
		})
	}
}
