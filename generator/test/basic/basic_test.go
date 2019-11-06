package basic

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

func TestBasic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before string
		run    Func
	}{{
		name: "Nullability",
		// language=GraphQL
		before: `
			mutation {
				createOneUser(data: {
					id: "nullability",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
					stuff: null,
				}) {
					id
				}
			}
		`,
		run: func(t *testing.T, client Client, ctx cx) {
			actual, err := client.User.FindOne(User.Email.Equals("john@example.com")).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			name, ok := actual.Name()
			assert.Equal(t, true, ok)
			assert.Equal(t, "John", name)

			stuff, ok := actual.Stuff()
			assert.Equal(t, false, ok)
			assert.Equal(t, "", stuff)
		},
	}, {
		name: "FindOne",
		// language=GraphQL
		before: `
			mutation {
				a: createOneUser(data: {
					id: "findOne1",
					email: "john@findOne.com",
					username: "john_doe",
				}) {
					id
				}
				b: createOneUser(data: {
					id: "findOne2",
					email: "jane@findOne.com",
					username: "jane_doe",
				}) {
					id
				}
			}
		`,
		run: func(t *testing.T, client Client, ctx cx) {
			actual, err := client.User.FindOne(User.Email.Equals("jane@findOne.com")).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, "findOne2", actual.ID)
		},
	}, {
		name: "FindMany",
		// language=GraphQL
		before: `
				mutation {
					a: createOneUser(data: {
						id: "findMany1",
						email: "1",
						username: "john",
						name: "a",
					}) {
						id
					}
					b: createOneUser(data: {
						id: "findMany2",
						email: "2",
						username: "john",
						name: "b",
					}) {
						id
					}
				}
			`,
		run: func(t *testing.T, client Client, ctx cx) {
			actual, err := client.User.FindMany(User.Username.Equals("john")).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, []UserModel{{
				user{
					ID:       "findMany1",
					Email:    "1",
					Username: "john",
					Name:     str("a"),
				},
			}, {
				user{
					ID:       "findMany2",
					Email:    "2",
					Username: "john",
					Name:     str("b"),
				},
			}}, actual)
		},
	}, {
		name: "Create",
		run: func(t *testing.T, client Client, ctx cx) {
			created, err := client.User.CreateOne(
				User.ID.Set("id"),
				User.Email.Set("email"),
				User.Username.Set("username"),

				// optional values
				User.Name.Set("name"),
				User.Stuff.Set("stuff"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				user{
					ID:       "id",
					Email:    "email",
					Username: "username",
					Name:     str("name"),
					Stuff:    str("stuff"),
				},
			}

			assert.Equal(t, expected, created)

			actual, err := client.User.FindOne(User.Email.Equals("email")).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "Update",
		// language=GraphQL
		before: `
			mutation {
				createOneUser(data: {
					id: "update",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
				}) {
					id
				}
			}
		`,
		run: func(t *testing.T, client Client, ctx cx) {
			email := "john@example.com"
			updated, err := client.User.UpdateOne(
				User.Email.Equals(email),
			).Data(
				// set required value
				User.Username.Set("new-username"),
				// set optional value
				User.Name.Set("New Name"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				user{
					ID:       "update",
					Email:    email,
					Username: "new-username",
					Name:     str("New Name"),
				},
			}

			assert.Equal(t, expected, updated)

			actual, err := client.User.FindOne(User.Email.Equals(email)).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "Delete",
		// language=GraphQL
		before: `
			mutation {
				createOneUser(data: {
					id: "delete",
					email: "john@example.com",
					username: "johndoe",
				}) {
					id
				}
			}
		`,
		run: func(t *testing.T, client Client, ctx cx) {
			email := "john@example.com"
			deleted, err := client.User.DeleteOne(
				User.Email.Equals(email),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				user{
					ID:       "delete",
					Email:    "john@example.com",
					Username: "johndoe",
				},
			}

			assert.Equal(t, expected, deleted)

			actual, err := client.User.FindOne(User.Email.Equals(email)).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, UserModel{}, actual)
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
