package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func str(v string) *string {
	return &v
}

func TestExportedBuilderFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "FindUnique",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "findOne1",
					email: "john@findOne.com",
					username: "john_doe",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "findOne2",
					email: "jane@findOne.com",
					username: "jane_doe",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var x UserEqualsUniqueWhereParam //nolint:gosimple
			x = User.Email.Equals("jane@findOne.com")
			actual, err := client.User.FindUnique(x).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, "findOne2", actual.ID)
		},
	}, {
		name: "FindMany",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "findMany1",
					email: "1",
					username: "john",
					name: "a",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "findMany2",
					email: "2",
					username: "john",
					name: "b",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var x []UserWhereParam
			x = append(x, User.Username.Equals("john"))
			actual, err := client.User.FindMany(x...).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, []UserModel{{
				InnerUser: InnerUser{
					ID:       "findMany1",
					Email:    "1",
					Username: "john",
					Name:     str("a"),
				},
			}, {
				InnerUser: InnerUser{
					ID:       "findMany2",
					Email:    "2",
					Username: "john",
					Name:     str("b"),
				},
			}}, actual)
		},
	}, {
		name: "Create",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var email UserWithPrismaEmailSetParam //nolint:gosimple
			email = User.Email.Set("email")

			var username UserWithPrismaUsernameSetParam //nolint:gosimple
			username = User.Username.Set("username")

			var optional []UserSetParam
			optional = append(optional, User.ID.Set("id"))
			optional = append(optional, User.Name.Set("name"))
			optional = append(optional, User.Stuff.Set("stuff"))

			created, err := client.User.CreateOne(
				email,
				username,
				optional...,
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &UserModel{
				InnerUser: InnerUser{
					ID:       "id",
					Email:    "email",
					Username: "username",
					Name:     str("name"),
					Stuff:    str("stuff"),
				},
			}

			assert.Equal(t, expected, created)

			actual, err := client.User.FindUnique(User.Email.Equals("email")).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "Update",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "update",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var params []UserSetParam
			params = append(params, User.Username.Set("new-username"))
			params = append(params, User.Name.Set("New Name"))

			email := "john@example.com"
			updated, err := client.User.FindUnique(
				User.Email.Equals(email),
			).Update(
				params...,
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &UserModel{
				InnerUser: InnerUser{
					ID:       "update",
					Email:    email,
					Username: "new-username",
					Name:     str("New Name"),
				},
			}

			assert.Equal(t, expected, updated)

			actual, err := client.User.FindUnique(User.Email.Equals(email)).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, expected, actual)
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
