package raw

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/prisma/prisma-client-go/runtime/types/raw"
	"github.com/prisma/prisma-client-go/test"
	"github.com/prisma/prisma-client-go/test/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

type RawUserModel struct {
	ID         raw.String    `json:"id"`
	Email      raw.String    `json:"email"`
	Username   raw.String    `json:"username"`
	Name       *raw.String   `json:"name"`
	Stuff      *raw.String   `json:"stuff"`
	Str        raw.String    `json:"str"`
	StrOpt     *raw.String   `json:"strOpt"`
	Int        raw.Int       `json:"int"`
	IntOpt     *raw.Int      `json:"intOpt"`
	Float      raw.Float     `json:"float"`
	FloatOpt   *raw.Float    `json:"floatOpt"`
	Bool       raw.Boolean   `json:"bool"`
	BoolOpt    *raw.Boolean  `json:"boolOpt"`
	Time       raw.DateTime  `json:"time"`
	TimeOpt    *raw.DateTime `json:"timeOpt"`
	Decimal    raw.Decimal   `json:"decimal"`
	DecimalOpt *raw.Decimal  `json:"decimalOpt"`
	JSON       raw.JSON      `json:"json"`
	JSONOpt    *raw.JSON     `json:"jsonOpt"`
}

func TestRaw(t *testing.T) {
	t.Parallel()

	var strOpt raw.String = "strOpt"
	var i raw.Int = 5
	var f raw.Float = 5.5
	var b raw.Boolean = false
	d := raw.Decimal{Decimal: decimal.NewFromFloat(5.5)}
	json := raw.JSON(`{"field":"value"}`)
	jsonOpt := &json

	dateOrig, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}

	date := raw.DateTime{Time: dateOrig}

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
					time: "2020-01-01T00:00:00Z",
					timeOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					decimal: 5.5,
					decimalOpt: 5.5,
					bool: true,
					boolOpt: false,
					json: "{\"field\":\"value\"}",
					jsonOpt: "{\"field\":\"value\"}",
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
					time: "2020-01-01T00:00:00Z",
					timeOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					decimal: 5.5,
					decimalOpt: 5.5,
					bool: true,
					boolOpt: false,
					json: "{\"field\":\"value\"}",
					jsonOpt: "{\"field\":\"value\"}",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []RawUserModel
			if err := client.Prisma.QueryRaw(`select * from "User"`).Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []RawUserModel{{
				ID:         "id1",
				Email:      "email1",
				Username:   "a",
				Str:        "str",
				StrOpt:     &strOpt,
				Int:        i,
				IntOpt:     &i,
				Float:      f,
				FloatOpt:   &f,
				Decimal:    d,
				DecimalOpt: &d,
				Bool:       true,
				BoolOpt:    &b,
				Time:       date,
				TimeOpt:    &date,
				JSON:       json,
				JSONOpt:    jsonOpt,
			}, {
				ID:         "id2",
				Email:      "email2",
				Username:   "b",
				Str:        "str",
				StrOpt:     &strOpt,
				Int:        i,
				IntOpt:     &i,
				Float:      f,
				FloatOpt:   &f,
				Decimal:    d,
				DecimalOpt: &d,
				Bool:       true,
				BoolOpt:    &b,
				Time:       date,
				TimeOpt:    &date,
				JSON:       json,
				JSONOpt:    jsonOpt,
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
					time: "2020-01-01T00:00:00Z",
					timeOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					decimal: 5.5,
					decimalOpt: 5.5,
					bool: true,
					boolOpt: false,
					json: "{\"field\":\"value\"}",
					jsonOpt: "{\"field\":\"value\"}",
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
					time: "2020-01-01T00:00:00Z",
					timeOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					decimal: 5.5,
					decimalOpt: 5.5,
					bool: true,
					boolOpt: false,
					json: "{\"field\":\"value\"}",
					jsonOpt: "{\"field\":\"value\"}",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []RawUserModel
			if err := client.Prisma.QueryRaw(`select * from "User" where id = $1`, "id2").Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []RawUserModel{{
				ID:         "id2",
				Email:      "email2",
				Username:   "b",
				Str:        "str",
				StrOpt:     &strOpt,
				Int:        i,
				IntOpt:     &i,
				Float:      f,
				FloatOpt:   &f,
				Decimal:    d,
				DecimalOpt: &d,
				Bool:       true,
				BoolOpt:    &b,
				Time:       date,
				TimeOpt:    &date,
				JSON:       json,
				JSONOpt:    jsonOpt,
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
					time: "2020-01-01T00:00:00Z",
					timeOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					decimal: 5.5,
					decimalOpt: 5.5,
					bool: true,
					boolOpt: false,
					json: "{\"field\":\"value\"}",
					jsonOpt: "{\"field\":\"value\"}",
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
					time: "2020-01-01T00:00:00Z",
					timeOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					decimal: 5.5,
					decimalOpt: 5.5,
					bool: true,
					boolOpt: false,
					json: "{\"field\":\"value\"}",
					jsonOpt: "{\"field\":\"value\"}",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []RawUserModel
			if err := client.Prisma.QueryRaw(`select * from "User" where id = $1 and email = $2`, "id2", "email2").Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []RawUserModel{{
				ID:         "id2",
				Email:      "email2",
				Username:   "b",
				Str:        "str",
				StrOpt:     &strOpt,
				Int:        i,
				IntOpt:     &i,
				Float:      f,
				FloatOpt:   &f,
				Decimal:    d,
				DecimalOpt: &d,
				Bool:       true,
				BoolOpt:    &b,
				Time:       date,
				TimeOpt:    &date,
				JSON:       json,
				JSONOpt:    jsonOpt,
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
					time: "2020-01-01T00:00:00Z",
					timeOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					decimal: 5.5,
					decimalOpt: 5.5,
					bool: true,
					boolOpt: false,
					json: "{\"field\":\"value\"}",
					jsonOpt: "{\"field\":\"value\"}",
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
					time: "2020-01-01T00:00:00Z",
					timeOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					decimal: 5.5,
					decimalOpt: 5.5,
					bool: true,
					boolOpt: false,
					json: "{\"field\":\"value\"}",
					jsonOpt: "{\"field\":\"value\"}",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []struct {
				Count struct {
					Value string `json:"prisma__value"`
				} `json:"count"`
			}
			if err := client.Prisma.QueryRaw(`select count(*) as count from "User"`).Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, "2", actual[0].Count.Value)
		},
	}, {
		name:   "insert into",
		before: []string{},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			result, err := client.Prisma.ExecuteRaw(`insert into "User" ("id", "email", "username", "str", "strOpt", "time", "timeOpt", "int", "intOpt", "float", "floatOpt", "decimal", "decimalOpt", "bool", "boolOpt", "json", "jsonOpt") values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`, "a", "a", "a", "a", "a", date, date, 1, 1, 2.0, 2.0, 2.0, 2.0, true, false, json, json).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, 1, result.Count)
		},
	}, {
		name: "update",
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "a",
					str: "str",
					strOpt: "strOpt",
					time: "2020-01-01T00:00:00Z",
					timeOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					decimal: 5.5,
					decimalOpt: 5.5,
					bool: true,
					boolOpt: false,
					json: "{\"field\":\"value\"}",
					jsonOpt: "{\"field\":\"value\"}",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			result, err := client.Prisma.ExecuteRaw(`update "User" set email = 'abc' where id = $1`, "id1").Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, 1, result.Count)

			result, err = client.Prisma.ExecuteRaw(`update "User" set email = 'abc' where id = $1`, "non-existing").Exec(ctx)
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
					time: "2010-01-01T00:00:00Z",
					timeOpt: "2010-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					decimal: 5.5,
					decimalOpt: 5.5,
					bool: true,
					boolOpt: false,
					json: "{\"field\":\"value\"}",
					jsonOpt: "{\"field\":\"value\"}",
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
					time: "2020-01-01T00:00:00Z",
					timeOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					decimal: 5.5,
					decimalOpt: 5.5,
					bool: true,
					boolOpt: false,
					json: "{\"field\":\"value\"}",
					jsonOpt: "{\"field\":\"value\"}",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []RawUserModel
			if err := client.Prisma.QueryRaw(`select * from "User" where "time" = $1`, date).Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []RawUserModel{{
				ID:         "id2",
				Email:      "email2",
				Username:   "b",
				Str:        "str",
				StrOpt:     &strOpt,
				Int:        i,
				IntOpt:     &i,
				Float:      f,
				FloatOpt:   &f,
				Decimal:    d,
				DecimalOpt: &d,
				Bool:       true,
				BoolOpt:    &b,
				Time:       date,
				TimeOpt:    &date,
				JSON:       json,
				JSONOpt:    jsonOpt,
			}}

			massert.Equal(t, expected, actual)
		},
	}, {
		name: "raw query with json parameter",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "a",
					str: "str",
					strOpt: "strOpt",
					time: "2010-01-01T00:00:00Z",
					timeOpt: "2010-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					decimal: 5.5,
					decimalOpt: 5.5,
					bool: true,
					boolOpt: false,
					json: "{\"field\":\"value\"}",
					jsonOpt: "{\"field\":\"value\"}",
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
					time: "2020-01-01T00:00:00Z",
					timeOpt: "2020-01-01T00:00:00Z",
					int: 5,
					intOpt: 5,
					float: 5.5,
					floatOpt: 5.5,
					decimal: 5.5,
					decimalOpt: 5.5,
					bool: true,
					boolOpt: false,
					json: "{\"field\":\"value\"}",
					jsonOpt: "{\"field\":\"value\"}",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []RawUserModel
			if err := client.Prisma.QueryRaw(`select * from "User" where "json" = $1`, json).Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []RawUserModel{{
				ID:         "id2",
				Email:      "email2",
				Username:   "b",
				Str:        "str",
				StrOpt:     &strOpt,
				Int:        i,
				IntOpt:     &i,
				Float:      f,
				FloatOpt:   &f,
				Decimal:    d,
				DecimalOpt: &d,
				Bool:       true,
				BoolOpt:    &b,
				Time:       date,
				TimeOpt:    &date,
				JSON:       json,
				JSONOpt:    jsonOpt,
			}}

			massert.Equal(t, expected, actual)
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
