package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestArrays(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "query for one",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					items: {
						set: ["a", "b", "c"],
					},
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			user, err := client.User.FindUnique(
				User.ID.Equals("id1"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				InternalUser: InternalUser{
					ID:    "id1",
					Items: []string{"a", "b", "c"},
				},
			}

			assert.Equal(t, expected, user)
		},
	}, {
		name: "create one",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			user, err := client.User.CreateOne(
				User.ID.Set("id"),
				User.Items.Set([]string{"a", "b", "c"}),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				InternalUser: InternalUser{
					ID:    "id",
					Items: []string{"a", "b", "c"},
				},
			}

			assert.Equal(t, expected, user)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()

			mockDB := test.Start(t, test.PostgreSQL, client.Engine, tt.before)
			defer test.End(t, test.PostgreSQL, client.Engine, mockDB)

			tt.run(t, client, context.Background())
		})
	}
}
