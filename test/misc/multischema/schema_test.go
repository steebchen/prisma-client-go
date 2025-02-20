//go:generate go run github.com/steebchen/prisma-client-go generate --schema .

package db

import (
	"context"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
	"github.com/steebchen/prisma-client-go/test/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestSchema(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "run",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.Post.FindMany(
				Post.ID.Equals("stuff"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, 0, len(actual))
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
