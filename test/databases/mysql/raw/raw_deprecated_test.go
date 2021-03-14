package raw

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

func TestRawDeprecated(t *testing.T) {
	strOpt := "strOpt"
	i := 5
	f := 5.5
	bTrue := 1
	bFalse := 0

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "raw query",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "a",
					str: "str",
					strOpt: "strOpt",
					date: "2020-01-01T00:00:00Z",
					dateOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "id2",
					email: "email2",
					username: "b",
					str: "str",
					strOpt: "strOpt",
					date: "2020-01-01T00:00:00Z",
					dateOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []RawUserModel
			if err := client.QueryRaw(`SELECT * FROM User`).Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []RawUserModel{{
				ID:       "id1",
				Email:    "email1",
				Username: "a",
				Str:      "str",
				StrOpt:   &strOpt,
				Int:      i,
				IntOpt:   &i,
				Float:    f,
				FloatOpt: &f,
				Bool:     bTrue,
				BoolOpt:  &bFalse,
			}, {
				ID:       "id2",
				Email:    "email2",
				Username: "b",
				Str:      "str",
				StrOpt:   &strOpt,
				Int:      i,
				IntOpt:   &i,
				Float:    f,
				FloatOpt: &f,
				Bool:     bTrue,
				BoolOpt:  &bFalse,
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "raw query with parameter",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "a",
					str: "str",
					strOpt: "strOpt",
					date: "2020-01-01T00:00:00Z",
					dateOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "id2",
					email: "email2",
					username: "b",
					str: "str",
					strOpt: "strOpt",
					date: "2020-01-01T00:00:00Z",
					dateOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []RawUserModel
			if err := client.QueryRaw(`SELECT * FROM User WHERE id = ?`, "id2").Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []RawUserModel{{
				ID:       "id2",
				Email:    "email2",
				Username: "b",
				Str:      "str",
				StrOpt:   &strOpt,
				Int:      i,
				IntOpt:   &i,
				Float:    f,
				FloatOpt: &f,
				Bool:     bTrue,
				BoolOpt:  &bFalse,
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "raw query with multiple parameters",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "a",
					str: "str",
					strOpt: "strOpt",
					date: "2020-01-01T00:00:00Z",
					dateOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "id2",
					email: "email2",
					username: "b",
					str: "str",
					strOpt: "strOpt",
					date: "2020-01-01T00:00:00Z",
					dateOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []RawUserModel
			if err := client.QueryRaw(`SELECT * FROM User WHERE id = ? AND email = ?`, "id2", "email2").Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []RawUserModel{{
				ID:       "id2",
				Email:    "email2",
				Username: "b",
				Str:      "str",
				StrOpt:   &strOpt,
				Int:      i,
				IntOpt:   &i,
				Float:    f,
				FloatOpt: &f,
				Bool:     bTrue,
				BoolOpt:  &bFalse,
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "raw query count",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "a",
					str: "str",
					strOpt: "strOpt",
					date: "2020-01-01T00:00:00Z",
					dateOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "id2",
					email: "email2",
					username: "b",
					str: "str",
					strOpt: "strOpt",
					date: "2020-01-01T00:00:00Z",
					dateOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []struct {
				Count int `json:"count"`
			}
			if err := client.QueryRaw(`SELECT COUNT(*) AS count FROM User`).Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, 2, actual[0].Count)
		},
	}, {
		name:   "insert into",
		before: []string{},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			date, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
			if err != nil {
				t.Fatal(err)
			}
			result, err := client.ExecuteRaw("INSERT INTO `User` (`id`, `email`, `username`, `str`, `strOpt`, `date`, `dateOpt`, `int`, `intOpt`, `float`, `floatOpt`, `bool`, `boolOpt`) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)", "a", "a", "a", "a", "a", date, date, 1, 1, 2.0, 2.0, true, false).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, 1, result.Count)
		},
	}, {
		name: "update",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "a",
					str: "str",
					strOpt: "strOpt",
					date: "2020-01-01T00:00:00Z",
					dateOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			result, err := client.ExecuteRaw("UPDATE `User` SET email = 'abc' WHERE id = ?", "id1").Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, 1, result.Count)

			result, err = client.ExecuteRaw("UPDATE `User` SET email = 'abc' WHERE id = ?", "non-existing").Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, 0, result.Count)
		},
	}, {
		name: "raw query with time parameter",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "a",
					str: "str",
					strOpt: "strOpt",
					date: "2010-01-01T00:00:00Z",
					dateOpt: "2010-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
				}) {
					id
				}
			}
		`, `
			mutation {
				result: createOneUser(data: {
					id: "id2",
					email: "email2",
					username: "b",
					str: "str",
					strOpt: "strOpt",
					date: "2020-01-01T00:00:00Z",
					dateOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					bool: true,
					boolOpt: false,
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			date, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
			if err != nil {
				t.Fatal(err)
			}
			var actual []RawUserModel
			if err := client.Prisma.QueryRaw("SELECT * FROM `User` WHERE date = ?", date).Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []RawUserModel{{
				ID:       "id2",
				Email:    "email2",
				Username: "b",
				Str:      "str",
				StrOpt:   &strOpt,
				Int:      i,
				IntOpt:   &i,
				Float:    f,
				FloatOpt: &f,
				Bool:     bTrue,
				BoolOpt:  &bFalse,
			}}

			assert.Equal(t, expected, actual)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()

			mockDB := test.Start(t, test.MySQL, client.Engine, tt.before)
			defer test.End(t, test.MySQL, client.Engine, mockDB)

			tt.run(t, client, context.Background())
		})
	}
}
