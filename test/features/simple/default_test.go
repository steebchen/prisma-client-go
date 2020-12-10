package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestSimple(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "simple",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			created, err := client.User.CreateOne(
				User.Name.Set("john"),
				User.ID.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			actual, err := client.User.FindOne(User.ID.Equals(created.ID)).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				InternalUser: InternalUser{
					ID:   "123",
					Name: "john",
					Age:  nil,
				},
			}

			assert.Equal(t, expected, actual)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, []test.Database{test.MySQL, test.PostgreSQL}, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
