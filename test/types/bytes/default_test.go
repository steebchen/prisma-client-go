package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "bytes create",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			a := []byte("a")
			b := []byte("b")
			created, err := client.User.CreateOne(
				User.Bytes.Set(a),
				User.BytesOpt.Set(b),
				User.ID.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &UserModel{
				InnerUser: InnerUser{
					ID:       "123",
					Bytes:    a,
					BytesOpt: &b,
				},
			}

			assert.Equal(t, expected, created)

			actual, err := client.User.FindUnique(User.ID.Equals(created.ID)).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "bytes find by bytes field",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			a := []byte("a")
			b := []byte("b")
			created, err := client.User.CreateOne(
				User.Bytes.Set(a),
				User.BytesOpt.Set(b),
				User.ID.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &UserModel{
				InnerUser: InnerUser{
					ID:       "123",
					Bytes:    a,
					BytesOpt: &b,
				},
			}

			assert.Equal(t, expected, created)

			actual, err := client.User.FindFirst(
				User.Bytes.Equals(a),
				User.BytesOpt.Equals(b),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
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
