package composite

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestComposite(t *testing.T) {
	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name:   "",
		before: nil,
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			expected := &UserModel{
				InnerUser: InnerUser{
					FirstName:  "a",
					MiddleName: "b",
					LastName:   "c",
				},
			}

			user, err := client.User.CreateOne(
				User.FirstName.Set("a"),
				User.MiddleName.Set("b"),
				User.LastName.Set("c"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expected, user)

			users, err := client.User.FindMany(
				User.FirstName.Equals("a"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, []UserModel{*expected}, users)

			oneUser, err := client.User.FindUnique(
				User.FirstNameMiddleNameLastName(
					User.FirstName.Equals("a"),
					User.MiddleName.Equals("b"),
					User.LastName.Equals("c"),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expected, oneUser)
		},
	}, {
		name:   "",
		before: nil,
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			expected := &UserModel{
				InnerUser: InnerUser{
					FirstName:  "a",
					MiddleName: "b",
					LastName:   "c",
				},
			}

			user, err := client.User.CreateOne(
				User.FirstName.Set("a"),
				User.MiddleName.Set("b"),
				User.LastName.Set("c"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expected, user)

			users, err := client.User.FindMany(
				User.FirstName.Equals("a"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, []UserModel{*expected}, users)

			oneUser, err := client.User.FindUnique(
				User.FirstNameLastName(
					User.FirstName.Equals("a"),
					User.LastName.Equals("c"),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expected, oneUser)
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
