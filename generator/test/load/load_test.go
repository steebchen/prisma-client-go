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
			type Result struct {
				// relation:
				// User struct{ *UserModel, Posts []Post }

				FindOneUser  UserModel   `json:"findOneUser"`
				FindManyUser []UserModel `json:"findManyUser"`
				FindManyPost []PostModel `json:"findManyPost"`
			}
			var actual Result
			err := client.Load(
				client.User.FindOne(
					User.ID.Equals("relations"),
					// User.Posts.Fetch(),
				),
				client.User.FindMany(User.Name.Equals("John")),
				client.Post.FindMany(Post.Title.Equals("common")),
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
