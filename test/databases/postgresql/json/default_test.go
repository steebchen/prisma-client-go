// Package db tests for json filters just for MySQL, as the generated client syntax differs
package db

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "json filter",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			x := struct {
				Attr string `json:"attr"`
			}{
				Attr: "stuff",
			}
			data, err := json.Marshal(x)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			created, err := client.User.CreateOne(
				User.JSON.Set(data),
				User.JSONOpt.Set([]byte(`"hi"`)),
				User.ID.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			var opt JSON = []byte(`"hi"`)
			expected := &UserModel{
				InnerUser: InnerUser{
					ID:      "123",
					JSON:    data,
					JSONOpt: &opt,
				},
			}

			assert.Equal(t, expected, created)

			actual, err := client.User.FindFirst(
				User.JSON.Path([]string{"attr"}),
				User.JSON.StringContains("stuff"),
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
			test.RunSerial(t, []test.Database{test.PostgreSQL}, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
