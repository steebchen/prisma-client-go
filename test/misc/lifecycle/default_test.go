package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/steebchen/prisma-client-go/test"
	"github.com/steebchen/prisma-client-go/test/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, db test.Database, client *PrismaClient, ctx cx)

func TestLifecycle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "connect success",
		run: func(t *testing.T, db test.Database, client *PrismaClient, ctx cx) {
			// manually setup testing
			mockDBName := db.SetupDatabase(t)
			test.Migrate(t, db, client.Engine, mockDBName)

			defer test.Teardown(t, db, mockDBName)

			if err := client.Connect(); err != nil {
				t.Fatalf("fail %s", err)
			}

			_, err := client.User.CreateOne(
				User.ID.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			if err := client.Disconnect(); err != nil {
				t.Fatalf("fail %s", err)
			}
		},
	}, {
		name: "connect pending",
		run: func(t *testing.T, db test.Database, client *PrismaClient, ctx cx) {

			// manually setup testing
			mockDBName := db.SetupDatabase(t)
			test.Migrate(t, db, client.Engine, mockDBName)

			defer test.Teardown(t, db, mockDBName)

			_, err := client.User.CreateOne(
				User.ID.Set("123"),
			).Exec(ctx)

			assert.NotEqual(t, err, nil)
			massert.Equal(t, "request failed: client is not connected yet", err.Error())
		},
	}, {
		name: "already disconnected",
		run: func(t *testing.T, db test.Database, client *PrismaClient, ctx cx) {
			// manually setup testing
			mockDBName := db.SetupDatabase(t)
			test.Migrate(t, db, client.Engine, mockDBName)

			defer test.Teardown(t, db, mockDBName)

			if err := client.Connect(); err != nil {
				t.Fatalf("fail %s", err)
			}

			_, err := client.User.CreateOne(
				User.ID.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			if err := client.Disconnect(); err != nil {
				t.Fatalf("fail %s", err)
			}

			_, err = client.User.CreateOne(
				User.ID.Set("456"),
			).Exec(ctx)
			assert.NotEqual(t, err, nil)
			massert.Equal(t, "request failed: client is already disconnected", err.Error())
		},
	}, {
		name: "connect err on async query engine error",
		run: func(t *testing.T, db test.Database, client *PrismaClient, ctx cx) {
			// use a non-working db URL to force a connection error

			err := client.Connect()

			assert.Regexp(t, "Environment variable not found: __REPLACE__", err.Error())
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, []test.Database{test.MySQL, test.PostgreSQL, test.MongoDB}, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()

				tt.run(t, db, client, context.Background())
			})
		})
	}
}
