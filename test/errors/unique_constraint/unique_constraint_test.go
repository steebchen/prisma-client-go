package db

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/engine"
	"github.com/prisma/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestIfPresent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "update operations",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.CreateOne(
				User.ID.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			_, err = client.User.CreateOne(
				User.ID.Set("123"),
			).Exec(ctx)
			expect := engine.UniqueConstraintViolationError{
				Message: "Unique constraint failed on the fields: (`id`)",
				Fields:  []string{"id"},
			}
			log.Printf("err: %s", err)
			assert.Error(t, err)
			v, ok := engine.IsUniqueConstraintViolationError(err)
			assert.Equal(t, &expect, v)
			assert.Equal(t, true, ok)
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
