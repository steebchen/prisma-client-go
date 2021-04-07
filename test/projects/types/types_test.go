package types

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func str(v string) *string {
	return &v
}

func TestTypes(t *testing.T) {
	t.Parallel()

	date, _ := time.Parse(RFC3339Milli, "2000-01-01T00:00:00+00:00")
	date = date.In(time.UTC)

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "complex strings",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			id := `f"hi"'`
			s := "\"'`\n\t}{*.,;:!?1234567890-_–=§±][äö€🤪"
			created, err := client.User.CreateOne(
				User.Str.Set(s),
				User.Int.Set(5),
				User.Float.Set(5.5),
				User.Bool.Set(true),
				User.Date.Set(date),
				User.Type.Set("x"),

				User.ID.Set(id),
				User.CreatedAt.Set(date),
				User.UpdatedAt.Set(date),
				User.StrOpt.Set(s),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &UserModel{
				InnerUser: InnerUser{
					ID:        id,
					CreatedAt: date,
					UpdatedAt: date,
					Str:       s,
					StrOpt:    &s,
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
					Type:      "x",
				},
			}

			assert.Equal(t, expected, created)

			actual, err := client.User.FindUnique(
				User.ID.Equals(id),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, expected, actual)

			actualSlice, err := client.User.FindMany(
				User.StrOpt.Equals(s),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, []UserModel{*expected}, actualSlice)
		},
	}, {
		name: "different field casing",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			created, err := client.User.CreateOne(
				User.Str.Set("str"),
				User.Int.Set(5),
				User.Float.Set(5.5),
				User.Bool.Set(true),
				User.Date.Set(date),
				User.Type.Set("x"),

				User.ID.Set("id"),
				User.CreatedAt.Set(date),
				User.UpdatedAt.Set(date),
				User.UpperCaseTest.Set("test1"),
				User.LowerCaseTest.Set("test2"),
				User.SnakeCaseTest.Set("test3"),
				User.WEiRdLycasEDTest.Set("test4"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &UserModel{
				InnerUser: InnerUser{
					ID:               "id",
					CreatedAt:        date,
					UpdatedAt:        date,
					Str:              "str",
					Int:              5,
					Float:            5.5,
					Bool:             true,
					Date:             date,
					Type:             "x",
					UpperCaseTest:    str("test1"),
					LowerCaseTest:    str("test2"),
					SnakeCaseTest:    str("test3"),
					WEiRdLycasEDTest: str("test4"),
				},
			}

			assert.Equal(t, expected, created)

			actualSlice, err := client.User.FindMany(
				User.UpperCaseTest.Equals("test1"),
				User.LowerCaseTest.Equals("test2"),
				User.SnakeCaseTest.Equals("test3"),
				User.WEiRdLycasEDTest.Equals("test4"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, []UserModel{*expected}, actualSlice)
		},
	}, {
		name: "basic equals",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "",
					strOpt: "str",
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
					type: "x",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			users, err := client.User.FindMany(
				User.ID.Equals("id"),
				User.StrOpt.Equals("str"),
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
				InnerUser: InnerUser{
					ID:        "id",
					CreatedAt: date,
					UpdatedAt: date,
					StrOpt:    str("str"),
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
					Type:      "x",
				},
			}}

			assert.Equal(t, expected, users)
		},
	}, {
		name: "advanced query",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "",
					strOpt: "alongstring",
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
					type: "x",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			before, _ := time.Parse(RFC3339Milli, "1999-01-01T00:00:00Z")

			users, err := client.User.FindMany(
				User.StrOpt.Contains("long"),
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
				InnerUser: InnerUser{
					ID:        "id",
					CreatedAt: date,
					UpdatedAt: date,
					StrOpt:    str("alongstring"),
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
					Type:      "x",
				},
			}}

			assert.Equal(t, expected, users)
		},
	}, {
		name: "failing query for the same field should lead to ErrNotFound",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "",
					strOpt: "alongstring",
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
					type: "x",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			before, _ := time.Parse(RFC3339Milli, "1999-01-01T00:00:00Z")

			_, err := client.User.FindFirst(
				User.StrOpt.Contains("long"),
				User.Bool.Equals(true),
				User.Int.GTE(5),
				User.Int.GT(10), // <- this is the failing part – this ensures all fields are considered in the query
				User.Int.LTE(5),
				User.Int.LT(7),
				User.Float.GTE(5.5),
				User.Float.GT(10), // <- this is the failing part – this ensures all fields are considered in the query
				User.Float.LTE(5.5),
				User.Float.LT(7.3),
				User.Date.Before(time.Now()),
				User.Date.After(before),
				User.CreatedAt.Equals(date),
				User.UpdatedAt.Equals(date),
			).Exec(ctx)
			assert.Equal(t, ErrNotFound, err)
		},
	}, {
		name: "IsNull",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "filled",
					strOpt: "filled",
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
					type: "x",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "id2",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "",
					strOpt: null,
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
					type: "x",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.User.FindMany(
				User.StrOpt.IsNull(),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				InnerUser: InnerUser{
					ID:        "id2",
					CreatedAt: date,
					UpdatedAt: date,
					StrOpt:    nil,
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
					Type:      "x",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "nullable dynamic nil field",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "",
					strOpt: "filled",
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
					type: "x",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "id2",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "",
					strOpt: null,
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
					type: "x",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var s *string = nil
			actual, err := client.User.FindMany(
				User.StrOpt.EqualsOptional(s),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				InnerUser: InnerUser{
					ID:        "id2",
					CreatedAt: date,
					UpdatedAt: date,
					StrOpt:    nil,
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
					Type:      "x",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "nullable dynamic field with value",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "",
					strOpt: null,
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
					type: "x",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "id2",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "",
					strOpt: "filled",
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
					type: "x",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			s := "filled"
			actual, err := client.User.FindMany(
				User.StrOpt.EqualsOptional(&s),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				InnerUser: InnerUser{
					ID:        "id2",
					CreatedAt: date,
					UpdatedAt: date,
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
					StrOpt:    &s,
					Type:      "x",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "IN operation",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "",
					strOpt: "first",
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
					type: "x",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "id2",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "",
					strOpt: "second",
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
					type: "x",
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "id3",
					createdAt: "2000-01-01T00:00:00Z",
					updatedAt: "2000-01-01T00:00:00Z",
					str: "",
					strOpt: "third",
					bool: true,
					date: "2000-01-01T00:00:00Z",
					int: 5,
					float: 5.5,
					type: "x",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			actual, err := client.User.FindMany(
				User.StrOpt.In([]string{"first", "third"}),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				InnerUser: InnerUser{
					ID:        "id1",
					CreatedAt: date,
					UpdatedAt: date,
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
					StrOpt:    str("first"),
					Type:      "x",
				},
			}, {
				InnerUser: InnerUser{
					ID:        "id3",
					CreatedAt: date,
					UpdatedAt: date,
					Int:       5,
					Float:     5.5,
					Bool:      true,
					Date:      date,
					StrOpt:    str("third"),
					Type:      "x",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, []test.Database{test.SQLite, test.MySQL, test.PostgreSQL}, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
