package db

import (
	"context"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
	"github.com/steebchen/prisma-client-go/test/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestFindFirst(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "find first",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "first",
					email: "john@example.com",
					username: "johndoe",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			first, err := client.User.FindFirst(
				User.Email.Equals("john@example.com"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &UserModel{
				InnerUser: InnerUser{
					ID:       "first",
					Email:    "john@example.com",
					Username: "johndoe",
				},
			}

			massert.Equal(t, expected, first)
		},
	}, {
		name: "return ErrNotFound",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.FindFirst(
				User.Email.Equals("john@example.com"),
			).Exec(ctx)
			massert.Equal(t, ErrNotFound, err)
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
