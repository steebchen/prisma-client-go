package db

import (
	"context"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
	"github.com/steebchen/prisma-client-go/test/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

var Databases = []test.Database{
	test.MySQL,
	test.PostgreSQL,
	test.SQLite,
}

func TestUpsert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "create",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.Post.UpsertOne(
				Post.ID.Equals("upsert"),
			).Create(
				Post.Title.Set("title"),
				Post.Views.Set(0),
				Post.ID.Set("upsert"),
			).Update(
				Post.Title.Set("title"),
				Post.Views.Increment(1),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &PostModel{
				InnerPost: InnerPost{
					ID:    "upsert",
					Title: "title",
					Views: 0,
				},
			}

			massert.Equal(t, expected, actual)
		},
	}, {
		"update",
		// language=GraphQL
		[]string{`
			mutation {
				result: createOnePost(data: {
					id: "upsert",
					title: "title",
					views: 0,
				}) {
					id
				}
			}
		`},
		func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.Post.UpsertOne(
				Post.ID.Equals("upsert"),
			).Create(
				Post.Title.Set("title"),
				Post.Views.Set(0),
				Post.ID.Set("upsert"),
			).Update(
				Post.Title.Set("title"),
				Post.Views.Increment(1),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &PostModel{
				InnerPost: InnerPost{
					ID:    "upsert",
					Title: "title",
					Views: 1,
				},
			}

			massert.Equal(t, expected, actual)
		},
	}, {
		name: "CreateOrUpdate when record does not exist",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.Post.UpsertOne(
				Post.ID.Equals("upsert"),
			).CreateOrUpdate(
				Post.Title.Set("title"),
				Post.Views.Set(0),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &PostModel{
				InnerPost: InnerPost{
					ID:    actual.ID,
					Title: "title",
					Views: 0,
				},
			}

			massert.Equal(t, expected, actual)
		},
	}, {
		name: "CreateOrUpdate when record exists",
		// Create the post first
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			desc := "random desc"
			_, err := client.Post.CreateOne(
				Post.Title.Set("title"),
				Post.Views.Set(0),
				Post.Description.Set(desc),
				Post.ID.Set("upsert"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			// Upsert the post: should only update it since it exists
			actual, err := client.Post.UpsertOne(
				Post.ID.Equals("upsert"),
			).CreateOrUpdate(
				Post.Title.Set("title"),
				Post.Views.Set(2),
				Post.Description.Set("random desc"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &PostModel{
				InnerPost: InnerPost{
					ID:          "upsert",
					Title:       "title",
					Views:       2,
					Description: &desc,
				},
			}

			massert.Equal(t, expected, actual)
		},
	}, {
		name: "transaction",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOnePost(data: {
					id: "upsert",
					title: "title",
					views: 0,
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			query := client.Post.UpsertOne(
				Post.ID.Equals("upsert"),
			).Create(
				Post.Title.Set("title"),
				Post.Views.Set(0),
				Post.ID.Set("upsert"),
			).Update(
				Post.Title.Set("title"),
				Post.Views.Increment(1),
			).Tx()

			if err := client.Prisma.Transaction(query).Exec(ctx); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &PostModel{
				InnerPost: InnerPost{
					ID:    "upsert",
					Title: "title",
					Views: 1,
				},
			}

			massert.Equal(t, expected, query.Result())
		},
	}, {
		name: "transaction CreateOrUpdate",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			// Create the post first
			desc := "random desc"
			_, err := client.Post.CreateOne(
				Post.Title.Set("title"),
				Post.Views.Set(2),
				Post.Description.Set(desc),
				Post.ID.Set("upsert"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			// Upsert: should only update it
			query := client.Post.UpsertOne(
				Post.ID.Equals("upsert"),
			).CreateOrUpdate(
				Post.Title.Set("title"),
				Post.Views.Set(3),
			).Tx()

			if err := client.Prisma.Transaction(query).Exec(ctx); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &PostModel{
				InnerPost: InnerPost{
					ID:          "upsert",
					Title:       "title",
					Views:       3,
					Description: &desc,
				},
			}

			massert.Equal(t, expected, query.Result())
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, Databases, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
