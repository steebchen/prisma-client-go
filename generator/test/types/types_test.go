package types

//go:generate prisma2 generate

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/photongo/generator/test/hooks"
)

type cx = context.Context
type Func func(t *testing.T, client Client, ctx cx)

func TestTypes(t *testing.T) {
	t.Parallel()

	t.Skip("blocked by ignored default values: https://github.com/prisma/prisma2/issues/964")

	tests := []struct {
		name   string
		before string
		run    Func
	}{{
		name: "complex strings",
		run: func(t *testing.T, client Client, ctx cx) {
			date, _ := time.Parse(RFC3339Milli, "2000-01-01T00:00:00Z")
			id := `f"hi"'`
			str := "\"'`\n\t}{*.,;:!?1234567890-_–=§±][äö€🤪"
			created, err := client.User.CreateOne(
				User.Str.Set(str),
				User.Int.Set(5),
				User.Float.Set(5.5),
				User.Bool.Set(true),
				User.Date.Set(date),

				User.ID.Set(id),
				User.CreatedAt.Set(date),
				User.UpdatedAt.Set(date),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				user{
					ID:        id,
					CreatedAt: date,
					UpdatedAt: date,
					Str:       str,
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
				},
			}

			assert.Equal(t, expected, created)

			actual, err := client.User.FindOne(
				User.ID.Equals(id),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, expected, actual)

			actualSlice, err := client.User.FindMany(
				User.Str.Equals(str),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, []UserModel{expected}, actualSlice)
		},
	}, {
		name: "enums",
		run: func(t *testing.T, client Client, ctx cx) {
			date, _ := time.Parse(RFC3339Milli, "2000-01-01T00:00:00Z")

			admin := RoleAdmin
			expected := UserModel{
				user{
					ID:    "123",
					Str:   "a",
					Int:   5,
					Float: 5.5,
					Bool:  true,
					Date:  date,
					Role:  &admin,
				},
			}

			created, err := client.User.CreateOne(
				User.Str.Set("a"),
				User.Int.Set(5),
				User.Float.Set(5.5),
				User.Bool.Set(true),
				User.Date.Set(date),
				User.Role.Set(RoleAdmin),
				User.ID.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, expected, created)

			actual, err := client.User.FindMany(
				User.Role.Equals(RoleAdmin),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, []UserModel{expected}, actual)
		},
	}, {
		name: "basic equals",
		// language=GraphQL
		before: `
			mutation {
				a: createOneUser(data: {
					id: "id",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "str",
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
				}) {
					id
				}
			}
		`,
		run: func(t *testing.T, client Client, ctx cx) {
			date, _ := time.Parse(RFC3339Milli, "2000-01-01T00:00:00Z")

			users, err := client.User.FindMany(
				User.ID.Equals("id"),
				User.Str.Equals("str"),
				User.Bool.Equals(true),
				User.Date.Equals(date),
				User.Float.Equals(5.5),
				User.Int.Equals(5),
				User.CreatedAt.Equals(date),
				User.UpdatedAt.Equals(date),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				user{
					ID:        "id",
					CreatedAt: date,
					UpdatedAt: date,
					Str:       "str",
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
				},
			}}

			assert.Equal(t, expected, users)
		},
	}, {
		name: "advanced query",
		// language=GraphQL
		before: `
			mutation {
				a: createOneUser(data: {
					id: "id",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "alongstring",
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
				}) {
					id
				}
			}
		`,
		run: func(t *testing.T, client Client, ctx cx) {
			date, _ := time.Parse(RFC3339Milli, "2000-01-01T00:00:00Z")
			before, _ := time.Parse(RFC3339Milli, "1999-01-01T00:00:00Z")

			users, err := client.User.FindMany(
				User.Str.Contains("long"),
				User.Bool.Equals(true),
				User.Int.GTE(5),
				User.Int.GT(3),
				User.Int.LTE(5),
				User.Int.LT(7),
				User.Float.GTE(5.5),
				User.Float.GT(2.7),
				User.Float.LTE(5.5),
				User.Float.LT(7.3),
				User.Date.Before(time.Now()),
				User.Date.After(before),
				User.CreatedAt.Equals(date),
				User.UpdatedAt.Equals(date),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				user{
					ID:        "id",
					CreatedAt: date,
					UpdatedAt: date,
					Str:       "alongstring",
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
				},
			}}

			assert.Equal(t, expected, users)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			hooks.Run(t)

			client := NewClient()
			if err := client.Connect(); err != nil {
				t.Fatalf("could not connect %s", err)
				return
			}

			defer func() {
				_ = client.Disconnect()
				// TODO blocked by prisma-engine panicking on disconnect
				// if err != nil {
				// 	t.Fatalf("could not disconnect %s", err)
				// }
			}()

			ctx := context.Background()

			if tt.before != "" {
				var response gqlResponse
				err := client.do(ctx, tt.before, &response)
				if err != nil {
					t.Fatalf("could not send mock query %s", err)
				}
				if response.Errors != nil {
					t.Fatalf("mock query has errors %+v", response)
				}
			}

			tt.run(t, client, ctx)
		})
	}
}
