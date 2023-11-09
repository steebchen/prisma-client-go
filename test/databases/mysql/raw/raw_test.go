package raw

import (
	"context"
	"testing"
	"time"

	"github.com/steebchen/prisma-client-go/test"
	"github.com/steebchen/prisma-client-go/test/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

type CustomRawUserModel struct {
	ID       string   `json:"id"`
	Email    string   `json:"email"`
	Username string   `json:"username"`
	Name     *string  `json:"name"`
	Stuff    *string  `json:"stuff"`
	Str      string   `json:"str"`
	StrOpt   *string  `json:"strOpt"`
	Int      int      `json:"int"`
	IntOpt   *int     `json:"intOpt"`
	Float    float64  `json:"float"`
	FloatOpt *float64 `json:"floatOpt"`
	Bool     int      `json:"bool"`
	BoolOpt  *int     `json:"boolOpt"`
}

func TestRaw(t *testing.T) {
	t.Parallel()

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
			var actual []CustomRawUserModel
			if err := client.Prisma.QueryRaw("select * from `User`").Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []CustomRawUserModel{{
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

			massert.Equal(t, expected, actual)
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
			var actual []CustomRawUserModel
			if err := client.Prisma.QueryRaw("select * from `User` where id = ?", "id2").Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []CustomRawUserModel{{
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

			massert.Equal(t, expected, actual)
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
			var actual []CustomRawUserModel
			if err := client.Prisma.QueryRaw("select * from `User` where id = ? and email = ?", "id2", "email2").Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []CustomRawUserModel{{
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

			massert.Equal(t, expected, actual)
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
				Count string `json:"count"`
			}
			if err := client.Prisma.QueryRaw("select count(*) as count from `User`").Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, "2", actual[0].Count)
		},
	}, {
		name:   "insert into",
		before: []string{},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			date, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
			if err != nil {
				t.Fatal(err)
			}
			result, err := client.Prisma.ExecuteRaw("insert into `User` (`id`, `email`, `username`, `str`, `strOpt`, `date`, `dateOpt`, `int`, `intOpt`, `float`, `floatOpt`, `bool`, `boolOpt`) values(?,?,?,?,?,?,?,?,?,?,?,?,?)", "a", "a", "a", "a", "a", date, date, 1, 1, 2.0, 2.0, true, false).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, 1, result.Count)
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
			result, err := client.Prisma.ExecuteRaw("update `User` set email = 'abc' where id = ?", "id1").Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, 1, result.Count)

			result, err = client.Prisma.ExecuteRaw("update `User` set email = 'abc' where id = ?", "non-existing").Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, 0, result.Count)
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
			var actual []CustomRawUserModel
			if err := client.Prisma.QueryRaw("select * from `User` where date = ?", date).Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []CustomRawUserModel{{
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

			massert.Equal(t, expected, actual)
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
