package load

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

func TestLoad(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before string
		run    Func
	}{{
		name: "multiple things",
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
			type Result struct {
				FindOneUser  UserModel   `json:"findOneUser"`
				FindManyUser []UserModel `json:"findManyUser"`
				FindManyPost []PostModel `json:"findManyPost"`
			}
			var actual Result
			err := client.Load(
				client.User.FindOne(
					User.ID.Equals("relations"),
				).Load(),
				client.User.FindMany(User.Name.Equals("John")).Load(),
				client.Post.FindMany(Post.Title.Equals("common")).Load(),
			).Exec(ctx, &actual)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := Result{
				FindOneUser: UserModel{
					user{
						ID:       "relations",
						Email:    "john@example.com",
						Username: "johndoe",
						Name:     str("John"),
					},
				},
				FindManyUser: []UserModel{
					{
						user{
							ID:       "relations",
							Email:    "john@example.com",
							Username: "johndoe",
							Name:     str("John"),
						},
					},
				},
				FindManyPost: []PostModel{
					{
						post{
							ID:      "a",
							Title:   "common",
							Content: str("a"),
						},
					},
					{
						post{
							ID:      "b",
							Title:   "common",
							Content: str("b"),
						},
					},
				},
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "fetch a relation",
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
		`,
		run: func(t *testing.T, client Client, ctx cx) {
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
					User.ID.Equals("relations"),
					// query for some users but filter by posts
					User.Posts.Some(
						Post.Title.Equals("2"),
					),

					// for those users which were found, additionally fetch posts
					User.Posts.
						FindMany(
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
							ID:       "relations",
							Email:    "john@example.com",
							Username: "johndoe",
							Name:     str("John"),
						},
					},
					Posts: []PostModel{{
						post{
							ID:      "d",
							Title:   "2",
							Content: str("stuff"),
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
			hooks.Run(t)

			client := NewClient()
			if err := client.Connect(); err != nil {
				t.Fatalf("could not connect %s", err)
				return
			}

			defer func() {
				err := client.Disconnect()
				if err != nil {
					t.Fatalf("could not disconnect: %s", err)
				}
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
