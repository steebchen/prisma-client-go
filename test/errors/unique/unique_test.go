package db

import (
	"context"
	"github.com/steebchen/prisma-client-go/runtime/types"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/steebchen/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestNotFound(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "unique constraint violation on email",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.CreateOne(
				User.Email.Set("john@example.com"),
				User.Username.Set("username"),
			).Exec(ctx)
			assert.Equal(t, nil, err)

			_, err = client.User.CreateOne(
				User.Email.Set("john@example.com"),
				User.Username.Set("username"),
			).Exec(ctx)

			violation, ok := types.CheckUniqueConstraintError(err) // IsUniqueConstraint
			assert.Equal(t, types.RealErrUniqueConstraint{
				Field: types.UserModelNameField, // User.Name.Field()
			}, violation)

			assert.Equal(t, true, ok)
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
