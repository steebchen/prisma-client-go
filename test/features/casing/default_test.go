package db

import (
	"context"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
	"github.com/steebchen/prisma-client-go/test/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestCaseSensitivity(t *testing.T) {
	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name:   "case sensitivity",
		before: nil,
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			user, err := client.User.CreateOne(
				User.Name.Set("THIS is me"),
				User.ID.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			expected := &UserModel{
				InnerUser: InnerUser{
					ID:   "123",
					Name: "THIS is me",
				},
			}

			massert.Equal(t, expected, user)

			user, err = client.User.FindFirst(
				User.And(
					User.Name.Contains("THIS"),
					User.Name.Mode(QueryModeInsensitive),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, expected, user)

			user, err = client.User.FindFirst(
				User.Name.Contains("this"),
				User.Name.Mode(QueryModeInsensitive),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, expected, user)
		},
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, []test.Database{test.MongoDB, test.PostgreSQL}, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
