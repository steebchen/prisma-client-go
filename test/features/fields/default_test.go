package db

import (
	"context"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
	"github.com/steebchen/prisma-client-go/test/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func str(v string) *string {
	return &v
}

func i(v int) *int {
	return &v
}

func TestFields(t *testing.T) {
	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name:   "omit",
		before: nil,
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			user, err := client.User.CreateOne(
				User.Name.Set("a"),
				User.Keep.Set("keep"),
				User.ID.Set("123"),
				User.Password.Set("password"),
				User.Age.Set(20),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, &UserModel{
				InnerUser: InnerUser{
					ID:       "123",
					Name:     "a",
					Keep:     "keep",
					Password: str("password"),
					Age:      i(20),
				},
			}, user)

			expected := &UserModel{
				InnerUser: InnerUser{
					Keep: "keep",
				},
			}

			users, err := client.User.FindMany(
				User.Name.Equals("a"),
			).Omit(
				User.ID.Field(),
				User.Password.Field(),
				User.Age.Field(),
				User.Name.Field(),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, []UserModel{*expected}, users)

			oneUser, err := client.User.FindUnique(
				User.ID.Equals("123"),
			).Omit(
				User.ID.Field(),
				User.Password.Field(),
				User.Age.Field(),
				User.Name.Field(),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, expected, oneUser)
		},
	}, {
		name:   "select",
		before: nil,
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			user, err := client.User.CreateOne(
				User.Name.Set("a"),
				User.Keep.Set("keep"),
				User.ID.Set("123"),
				User.Password.Set("password"),
				User.Age.Set(20),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, &UserModel{
				InnerUser: InnerUser{
					ID:       "123",
					Name:     "a",
					Keep:     "keep",
					Password: str("password"),
					Age:      i(20),
				},
			}, user)

			expected := &UserModel{
				InnerUser: InnerUser{
					Keep: "keep",
				},
			}

			users, err := client.User.FindMany(
				User.Name.Equals("a"),
			).Select(
				User.Keep.Field(),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, []UserModel{*expected}, users)

			oneUser, err := client.User.FindUnique(
				User.ID.Equals("123"),
			).Select(
				User.Keep.Field(),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, expected, oneUser)
		},
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, test.Databases, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
