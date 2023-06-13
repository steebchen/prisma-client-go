package raw

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/steebchen/prisma-client-go/test"
	"github.com/steebchen/prisma-client-go/test/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestRawTypesInternal(t *testing.T) {
	t.Parallel()

	var strOpt RawString = "strOpt"
	var i RawInt = 5
	var f RawFloat = 5.5
	var b RawBoolean = false
	var d = RawDecimal{Decimal: decimal.NewFromFloat(5.5)}
	var jsn = RawJSON{RawMessage: []byte(`{"field":"value"}`)}
	var jsonOpt = &jsn
	var bytes RawBytes = []byte(`{"field":"value"}`)
	var bytesOpt = &bytes

	dateOrig, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}

	date := RawDateTime{Time: dateOrig}

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
					bytes: "eyJmaWVsZCI6InZhbHVlIn0=",
					bytesOpt: "eyJmaWVsZCI6InZhbHVlIn0=",
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
					bytes: "eyJmaWVsZCI6InZhbHVlIn0=",
					bytesOpt: "eyJmaWVsZCI6InZhbHVlIn0=",
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
				JSON:       jsn,
				JSONOpt:    jsonOpt,
				Bytes:      bytes,
				BytesOpt:   bytesOpt,
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
				JSON:       jsn,
				JSONOpt:    jsonOpt,
				Bytes:      bytes,
				BytesOpt:   bytesOpt,
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
					bytes: "eyJmaWVsZCI6InZhbHVlIn0=",
					bytesOpt: "eyJmaWVsZCI6InZhbHVlIn0=",
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
					bytes: "eyJmaWVsZCI6InZhbHVlIn0=",
					bytesOpt: "eyJmaWVsZCI6InZhbHVlIn0=",
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
				JSON:       jsn,
				JSONOpt:    jsonOpt,
				Bytes:      bytes,
				BytesOpt:   bytesOpt,
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
					bytes: "eyJmaWVsZCI6InZhbHVlIn0=",
					bytesOpt: "eyJmaWVsZCI6InZhbHVlIn0=",
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
					bytes: "eyJmaWVsZCI6InZhbHVlIn0=",
					bytesOpt: "eyJmaWVsZCI6InZhbHVlIn0=",
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
				JSON:       jsn,
				JSONOpt:    jsonOpt,
				Bytes:      bytes,
				BytesOpt:   bytesOpt,
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
					bytes: "eyJmaWVsZCI6InZhbHVlIn0=",
					bytesOpt: "eyJmaWVsZCI6InZhbHVlIn0=",
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
					bytes: "eyJmaWVsZCI6InZhbHVlIn0=",
					bytesOpt: "eyJmaWVsZCI6InZhbHVlIn0=",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []struct {
				Count BigInt `json:"count"`
			}
			if err := client.Prisma.QueryRaw(`select count(*) as count from "User"`).Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, 2, actual[0].Count)
		},
	}, {
		name:   "insert into",
		before: []string{},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			result, err := client.Prisma.ExecuteRaw(`insert into "User" ("id", "email", "username", "str", "strOpt", "time", "timeOpt", "int", "intOpt", "float", "floatOpt", "decimal", "decimalOpt", "bool", "boolOpt", "json", "jsonOpt", "bytes", "bytesOpt") values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)`, "a", "a", "a", "a", "a", date, &date, 1, 1, 2.0, 2.0, 2.0, 2.0, true, false, jsn, jsonOpt, bytes, bytesOpt).Exec(ctx)
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
					bytes: "eyJmaWVsZCI6InZhbHVlIn0=",
					bytesOpt: "eyJmaWVsZCI6InZhbHVlIn0=",
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
					bytes: "eyJmaWVsZCI6InZhbHVlIn0=",
					bytesOpt: "eyJmaWVsZCI6InZhbHVlIn0=",
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
					bytes: "eyJmaWVsZCI6InZhbHVlIn0=",
					bytesOpt: "eyJmaWVsZCI6InZhbHVlIn0=",
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
				JSON:       jsn,
				JSONOpt:    jsonOpt,
				Bytes:      bytes,
				BytesOpt:   bytesOpt,
			}}

			massert.Equal(t, expected, actual)
		},
	}, {
		name: "raw query with json parameter",
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
					json: "{\"field\":\"a\"}",
					jsonOpt: "{\"field\":\"a\"}",
					bytes: "eyJmaWVsZCI6InZhbHVlIn0=",
					bytesOpt: "eyJmaWVsZCI6InZhbHVlIn0=",
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
					json: "{\"field\":\"b\"}",
					jsonOpt: "{\"field\":\"b\"}",
					bytes: "eyJmaWVsZCI6InZhbHVlIn0=",
					bytesOpt: "eyJmaWVsZCI6InZhbHVlIn0=",
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			j := RawJSON{RawMessage: []byte(`{"field":"b"}`)}

			var actual []RawUserModel
			if err := client.Prisma.QueryRaw(`select * from "User" where "json"->>'field' = 'b'`, j).Exec(ctx, &actual); err != nil {
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
				JSON:       j,
				JSONOpt:    &j,
				Bytes:      bytes,
				BytesOpt:   bytesOpt,
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
