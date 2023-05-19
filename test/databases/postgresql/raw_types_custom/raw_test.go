package raw

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/prisma/prisma-client-go/test"
	"github.com/prisma/prisma-client-go/test/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

type RawUserModel1 struct {
	ID           RawString    `json:"id"`
	Email        RawString    `json:"email"`
	Username     RawString    `json:"username"`
	Str          RawString    `json:"str"`
	StrOpt       *RawString   `json:"strOpt,omitempty"`
	StrEmpty     *RawString   `json:"strEmpty,omitempty"`
	Time         RawDateTime  `json:"time"`
	TimeOpt      *RawDateTime `json:"timeOpt,omitempty"`
	TimeEmpty    *RawDateTime `json:"timeEmpty,omitempty"`
	Int          RawInt       `json:"int"`
	IntOpt       *RawInt      `json:"intOpt,omitempty"`
	IntEmpty     *RawInt      `json:"intEmpty,omitempty"`
	Float        RawFloat     `json:"float"`
	FloatOpt     *RawFloat    `json:"floatOpt,omitempty"`
	FloatEmpty   *RawFloat    `json:"floatEmpty,omitempty"`
	Bool         RawBoolean   `json:"bool"`
	BoolOpt      *RawBoolean  `json:"boolOpt,omitempty"`
	BoolEmpty    *RawBoolean  `json:"boolEmpty,omitempty"`
	Decimal      RawDecimal   `json:"decimal"`
	DecimalOpt   *RawDecimal  `json:"decimalOpt,omitempty"`
	DecimalEmpty *RawDecimal  `json:"decimalEmpty,omitempty"`
	JSON         RawJSON      `json:"json"`
	JSONOpt      *RawJSON     `json:"jsonOpt,omitempty"`
	JSONEmpty    *RawJSON     `json:"jsonEmpty,omitempty"`
	Bytes        RawBytes     `json:"bytes"`
	BytesOpt     *RawBytes    `json:"bytesOpt,omitempty"`
	BytesEmpty   *RawBytes    `json:"bytesEmpty,omitempty"`
}

type CustomRawUserModel struct {
	ID           string           `json:"id"`
	Email        string           `json:"email"`
	Username     string           `json:"username"`
	Str          string           `json:"str"`
	StrOpt       *string          `json:"strOpt,omitempty"`
	StrEmpty     *string          `json:"strEmpty,omitempty"`
	Time         time.Time        `json:"time"`
	TimeOpt      *time.Time       `json:"timeOpt,omitempty"`
	TimeEmpty    *time.Time       `json:"timeEmpty,omitempty"`
	Int          int              `json:"int"`
	IntOpt       *int             `json:"intOpt,omitempty"`
	IntEmpty     *int             `json:"intEmpty,omitempty"`
	Float        float64          `json:"float"`
	FloatOpt     *float64         `json:"floatOpt,omitempty"`
	FloatEmpty   *float64         `json:"floatEmpty,omitempty"`
	Bool         bool             `json:"bool"`
	BoolOpt      *bool            `json:"boolOpt,omitempty"`
	BoolEmpty    *bool            `json:"boolEmpty,omitempty"`
	Decimal      decimal.Decimal  `json:"decimal"`
	DecimalOpt   *decimal.Decimal `json:"decimalOpt,omitempty"`
	DecimalEmpty *decimal.Decimal `json:"decimalEmpty,omitempty"`
	JSON         json.RawMessage  `json:"json"`
	JSONOpt      *json.RawMessage `json:"jsonOpt,omitempty"`
	JSONEmpty    *json.RawMessage `json:"jsonEmpty,omitempty"`
	Bytes        []byte           `json:"bytes"`
	BytesOpt     *[]byte          `json:"bytesOpt,omitempty"`
	BytesEmpty   *[]byte          `json:"bytesEmpty,omitempty"`
}

func TestRawTypesCustom(t *testing.T) {
	t.Parallel()

	strOpt := "strOpt"
	i := 5
	f := 5.5
	b := false
	d := decimal.NewFromFloat(5.5)
	jsn := json.RawMessage(`{"field":"value"}`)
	jsonOpt := &jsn
	bytes := []byte(`{"field":"value"}`)
	bytesOpt := &bytes

	date, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}

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
			var actual []CustomRawUserModel
			if err := client.Prisma.QueryRaw(`select * from "User"`).Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []CustomRawUserModel{{
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
			var actual []CustomRawUserModel
			if err := client.Prisma.QueryRaw(`select * from "User" where id = $1`, "id2").Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []CustomRawUserModel{{
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
			var actual []CustomRawUserModel
			if err := client.Prisma.QueryRaw(`select * from "User" where id = $1 and email = $2`, "id2", "email2").Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []CustomRawUserModel{{
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
			var actual []CustomRawUserModel
			if err := client.Prisma.QueryRaw(`select * from "User" where "time" = $1`, date).Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []CustomRawUserModel{{
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
			j := json.RawMessage(`{"field":"b"}`)

			var actual []CustomRawUserModel
			if err := client.Prisma.QueryRaw(`select * from "User" where "json"->>'field' = 'b'`, j).Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []CustomRawUserModel{{
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
