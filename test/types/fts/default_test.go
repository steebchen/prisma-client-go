package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestFullTextSearch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "full text search",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "a",
					name: "john doe",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "b",
					name: "jane doe",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "c",
					name: "unknown dude",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.User.FindMany(
				User.Name.Search("doe"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{
				{
					InnerUser: InnerUser{
						ID:   "a",
						Name: "john doe",
					},
				},
				{
					InnerUser: InnerUser{
						ID:   "b",
						Name: "jane doe",
					},
				},
			}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "relevance",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "123",
					name: "john doe",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			users, err := client.User.FindMany().OrderBy(
				User.Relevance_.Search("hn"),
				User.Relevance_.Fields([]UserOrderByRelevanceFieldEnum{UserOrderByRelevanceFieldEnumName}),
				User.Relevance_.Sort(SortOrderDesc),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			expected := &UserModel{
				InnerUser: InnerUser{
					ID:   "123",
					Name: "john doe",
				},
			}

			assert.Equal(t, []UserModel{*expected}, users)
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
