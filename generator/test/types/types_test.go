package types

//go:generate go run github.com/prisma/photongo generate

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/photongo/generator/test/hooks"
)

type cx = context.Context
type Func func(t *testing.T, client *Client, ctx cx)

func str(v string) *string {
	return &v
}

func TestTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before string
		run    Func
	}{{
		name: "complex strings",
		run: func(t *testing.T, client *Client, ctx cx) {
			date, _ := time.Parse(RFC3339Milli, "2000-01-01T00:00:00Z")
			id := `f"hi"'`
			s := "\"'`\n\t}{*.,;:!?1234567890-_–=§±][äö€🤪"
			created, err := client.User.CreateOne(
				User.Int.Set(5),
				User.Float.Set(5.5),
				User.Bool.Set(true),
				User.Date.Set(date),

				User.ID.Set(id),
				User.CreatedAt.Set(date),
				User.UpdatedAt.Set(date),
				User.Str.Set(s),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				user{
					ID:        id,
					CreatedAt: date,
					UpdatedAt: date,
					Str:       &s,
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
				User.Str.Equals(s),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, []UserModel{expected}, actualSlice)
		},
	}, {
		name: "enums",
		run: func(t *testing.T, client *Client, ctx cx) {
			date, _ := time.Parse(RFC3339Milli, "2000-01-01T00:00:00Z")

			admin := RoleAdmin
			expected := UserModel{
				user{
					ID:        "123",
					CreatedAt: date,
					UpdatedAt: date,
					Str:       str("a"),
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
					Role:      &admin,
				},
			}

			created, err := client.User.CreateOne(
				User.Int.Set(5),
				User.Float.Set(5.5),
				User.Bool.Set(true),
				User.Date.Set(date),
				User.Role.Set(RoleAdmin),

				User.ID.Set("123"),
				User.Str.Set("a"),
				User.CreatedAt.Set(date),
				User.UpdatedAt.Set(date),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, expected, created)

			actual, err := client.User.FindMany(
				User.Role.Equals(RoleAdmin),
				User.Role.In([]Role{RoleAdmin}),
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
		run: func(t *testing.T, client *Client, ctx cx) {
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
					Str:       str("str"),
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
		run: func(t *testing.T, client *Client, ctx cx) {
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
					Str:       str("alongstring"),
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
				},
			}}

			assert.Equal(t, expected, users)
		},
	}, {
		name: "IsNull",
		// language=GraphQL
		before: `
			mutation {
				a: createOneUser(data: {
					id: "id1",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "filled",
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
				}) {
					id
				}
				b: createOneUser(data: {
					id: "id2",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: null,
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
				}) {
					id
				}
			}
		`,
		run: func(t *testing.T, client *Client, ctx cx) {
			date, _ := time.Parse(RFC3339Milli, "2000-01-01T00:00:00Z")

			actual, err := client.User.FindMany(
				User.Str.IsNull(),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				user{
					ID:        "id2",
					CreatedAt: date,
					UpdatedAt: date,
					Str:       nil,
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "nullable dynamic nil field",
		// language=GraphQL
		before: `
			mutation {
				a: createOneUser(data: {
					id: "id1",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "filled",
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
				}) {
					id
				}
				b: createOneUser(data: {
					id: "id2",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: null,
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
				}) {
					id
				}
			}
		`,
		run: func(t *testing.T, client *Client, ctx cx) {
			date, _ := time.Parse(RFC3339Milli, "2000-01-01T00:00:00Z")

			var s *string = nil
			actual, err := client.User.FindMany(
				User.Str.EqualsOptional(s),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				user{
					ID:        "id2",
					CreatedAt: date,
					UpdatedAt: date,
					Str:       nil,
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "nullable dynamic field with value",
		// language=GraphQL
		before: `
			mutation {
				a: createOneUser(data: {
					id: "id1",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: null,
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
				}) {
					id
				}
				b: createOneUser(data: {
					id: "id2",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "filled",
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
				}) {
					id
				}
			}
		`,
		run: func(t *testing.T, client *Client, ctx cx) {
			date, _ := time.Parse(RFC3339Milli, "2000-01-01T00:00:00Z")

			s := "filled"
			actual, err := client.User.FindMany(
				User.Str.EqualsOptional(&s),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				user{
					ID:        "id2",
					CreatedAt: date,
					UpdatedAt: date,
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
					Str:       &s,
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "IN operation",
		// language=GraphQL
		before: `
			mutation {
				a: createOneUser(data: {
					id: "id1",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "first",
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
				}) {
					id
				}
				b: createOneUser(data: {
					id: "id2",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "second",
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
				}) {
					id
				}
				c: createOneUser(data: {
					id: "id3",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "third",
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
				}) {
					id
				}
			}
		`,
		run: func(t *testing.T, client *Client, ctx cx) {
			date, _ := time.Parse(RFC3339Milli, "2000-01-01T00:00:00Z")

			actual, err := client.User.FindMany(
				User.Str.In([]string{"first", "third"}),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				user{
					ID:        "id1",
					CreatedAt: date,
					UpdatedAt: date,
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
					Str:       str("first"),
				},
			}, {
				user{
					ID:        "id3",
					CreatedAt: date,
					UpdatedAt: date,
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
					Str:       str("third"),
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			hooks.Start(t, client, tt.before, client.do)
			tt.run(t, client, context.Background())
			hooks.End(t, client)
		})
	}
}
