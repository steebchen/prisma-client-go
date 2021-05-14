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
		name: "json create",
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

			actual, err := client.User.FindUnique(User.ID.Equals(created.ID)).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "json find by json field",
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
				User.JSON.Equals(data),
				User.JSONOpt.Equals(opt),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "json update field",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var b JSON = []byte(`"hi"`)
			created, err := client.User.CreateOne(
				User.JSON.Set(b),
				User.JSONOpt.Set(b),
				User.ID.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &UserModel{
				InnerUser: InnerUser{
					ID:      "123",
					JSON:    b,
					JSONOpt: &b,
				},
			}

			assert.Equal(t, expected, created)

			updated, err := client.User.FindUnique(
				User.ID.Equals("123"),
			).Update(
				User.JSON.Set(b),
				User.JSONOpt.Set(b),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, expected, updated)

			actual, err := client.User.FindUnique(
				User.ID.Equals("123"),
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
