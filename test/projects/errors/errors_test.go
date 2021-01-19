package errors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestBasic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "FindUnique not found",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.FindUnique(User.Email.Equals("404")).Exec(ctx)

			assert.Equal(t, ErrNotFound, err)
		},
	}, {
		name: "Update not found",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.FindUnique(
				User.Email.Equals("404"),
			).Update(
				User.Name.Set("x"),
			).Exec(ctx)

			assert.Equal(t, ErrNotFound, err)
		},
	}, {
		name: "Delete not found",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.FindUnique(
				User.Email.Equals("404"),
			).Delete().Exec(ctx)

			assert.Equal(t, ErrNotFound, err)
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
