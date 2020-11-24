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

func TestIfPresent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
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
			maybeQuery := "johndoe"
			newUsername := "new-username"
			email := "john@example.com"
			_, err := client.User.FindMany(
				User.Email.Equals(email),
				// query for this one
				User.Username.EqualsIfPresent(&maybeQuery),
				// ignore this one
				User.Name.EqualsIfPresent(nil),
			).Update(
				// set value
				User.Username.SetIfPresent(&newUsername),
				// don't set because nil
				User.Name.SetIfPresent(nil),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			updated, err := client.User.FindOne(
				User.ID.Equals("update"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				InternalUser: InternalUser{
					ID:       "update",
					Email:    email,
					Username: "new-username",
					Name:     str("John"),
				},
			}

			assert.Equal(t, expected, updated)

			actual, err := client.User.FindOne(User.Email.Equals(email)).Exec(ctx)
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
