package basic

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

func TestBasic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "Nullability",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "nullability",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
					stuff: null,
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.User.FindUnique(User.Email.Equals("john@example.com")).Exec(ctx)
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
		name: "marshal json",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "marshal",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
					stuff: null,
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			user, err := client.User.FindUnique(User.Email.Equals("john@example.com")).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			actual, err := json.Marshal(&user)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := `{"id":"marshal","email":"john@example.com","username":"johndoe","name":"John","stuff":null}`
			assert.Equal(t, expected, string(actual))
		},
	}, {
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
			actual, err := client.User.FindUnique(User.Email.Equals("jane@findOne.com")).Exec(ctx)
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
			actual, err := client.User.FindMany(User.Username.Equals("john")).Exec(ctx)
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
		name: "FindMany all",
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
			actual, err := client.User.FindMany().Exec(ctx)
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
		name: "FindMany empty",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.User.FindMany().Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, []UserModel{}, actual)
		},
	}, {
		name: "Create",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			created, err := client.User.CreateOne(
				User.Email.Set("email"),
				User.Username.Set("username"),

				// optional values
				User.ID.Set("id"),
				User.Name.Set("name"),
				User.Stuff.Set("stuff"),
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
		name: "Create with optional values",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			name := "name"
			created, err := client.User.CreateOne(
				User.Email.Set("email"),
				User.Username.Set("username"),

				// optional values
				User.ID.Set("id"),
				// set one to a value
				User.Name.SetOptional(&name),
				// set one to nil pointer
				User.Stuff.SetOptional(nil),
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
					Stuff:    nil,
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
			email := "john@example.com"
			updated, err := client.User.FindUnique(
				User.Email.Equals(email),
			).Update(
				// set required value
				User.Username.Set("new-username"),
				// set optional value
				User.Name.Set("New Name"),
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
	}, {
		name: "Update many",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "username",
					name: "1",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "id2",
					email: "email2",
					username: "username",
					name: "2",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			result, err := client.User.FindMany(
				User.Username.Equals("username"),
			).Update(
				User.Name.Set("New Name"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, 2, result.Count)

			actual, err := client.User.FindMany(
				User.Username.Equals("username"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				InnerUser: InnerUser{
					ID:       "id1",
					Email:    "email1",
					Username: "username",
					Name:     str("New Name"),
				},
			}, {
				InnerUser: InnerUser{
					ID:       "id2",
					Email:    "email2",
					Username: "username",
					Name:     str("New Name"),
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "Delete",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "delete",
					email: "john@example.com",
					username: "johndoe",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			email := "john@example.com"
			deleted, err := client.User.FindUnique(
				User.Email.Equals(email),
			).Delete().Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &UserModel{
				InnerUser: InnerUser{
					ID:       "delete",
					Email:    "john@example.com",
					Username: "johndoe",
				},
			}

			assert.Equal(t, expected, deleted)

			actual, err := client.User.FindUnique(User.Email.Equals(email)).Exec(ctx)
			assert.Equal(t, ErrNotFound, err)
			assert.Equal(t, true, actual == nil)
		},
	}, {
		name: "Delete tx",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "delete",
					email: "john@example.com",
					username: "johndoe",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			email := "john@example.com"
			query := client.User.FindUnique(
				User.Email.Equals(email),
			).Delete().Tx()

			if err := client.Prisma.Transaction(query).Exec(ctx); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &UserModel{
				InnerUser: InnerUser{
					ID:       "delete",
					Email:    "john@example.com",
					Username: "johndoe",
				},
			}

			assert.Equal(t, expected, query.Result())

			actual, err := client.User.FindUnique(User.Email.Equals(email)).Exec(ctx)
			assert.Equal(t, ErrNotFound, err)
			assert.Equal(t, true, actual == nil)
		},
	}, {
		name: "Delete many",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "username",
					name: "1",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "id2",
					email: "email2",
					username: "username",
					name: "2",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			result, err := client.User.FindMany(
				User.Username.Equals("username"),
			).Delete().Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, 2, result.Count)

			actual, err := client.User.FindMany(
				User.Username.Equals("username"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, []UserModel{}, actual)
		},
	}, {
		name: "NOT operation",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "username",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "id2",
					email: "email2",
					username: "username",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.User.FindMany(
				User.Not(
					User.Email.Equals("email1"),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				InnerUser: InnerUser{
					ID:       "id2",
					Email:    "email2",
					Username: "username",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "OR operation",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "a",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "id2",
					email: "email2",
					username: "b",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.User.FindMany(
				User.Or(
					User.Email.Equals("email1"),
					User.ID.Equals("id2"),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				InnerUser: InnerUser{
					ID:       "id1",
					Email:    "email1",
					Username: "a",
				},
			}, {
				InnerUser: InnerUser{
					ID:       "id2",
					Email:    "email2",
					Username: "b",
				},
			}}

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
