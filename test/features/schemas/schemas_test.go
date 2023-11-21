package db

import (
	"context"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
	"github.com/steebchen/prisma-client-go/test/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestSchemas(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "create and find first",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			expected := &FirstModel{
				InnerFirst: InnerFirst{
					ID: "123",
				},
			}

			created, err := client.First.CreateOne(
				First.ID.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, expected, created)

			actual, err := client.First.FindUnique(
				First.ID.Equals("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, expected, actual)
		},
	}, {
		name: "create and find second",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			expected := &SecondModel{
				InnerSecond: InnerSecond{
					ID: "123",
				},
			}

			created, err := client.Second.CreateOne(
				Second.ID.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, expected, created)

			actual, err := client.Second.FindUnique(
				Second.ID.Equals("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, expected, actual)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, []test.Database{test.PostgreSQL}, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
