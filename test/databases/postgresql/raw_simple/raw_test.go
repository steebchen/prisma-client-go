package raw

import (
	"context"
	"testing"

	"github.com/prisma/prisma-client-go/test"
	"github.com/prisma/prisma-client-go/test/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

type CustomRawUserModel struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func TestRawSimple(t *testing.T) {
	t.Parallel()

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
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []UserModel
			if err := client.Prisma.QueryRaw(`select * from "User"`).Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []CustomRawUserModel{{
				ID:       "id1",
				Email:    "email1",
				Username: "a",
			}, {
				ID:       "id2",
				Email:    "email2",
				Username: "b",
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
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []UserModel
			if err := client.Prisma.QueryRaw(`select * from "User" where id = $1`, "id2").Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []CustomRawUserModel{{
				ID:       "id2",
				Email:    "email2",
				Username: "b",
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
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			var actual []UserModel
			if err := client.Prisma.QueryRaw(`select * from "User" where id = $1 and username = $2`, "id2", "b").Exec(ctx, &actual); err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []CustomRawUserModel{{
				ID:       "id2",
				Email:    "email2",
				Username: "b",
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
			result, err := client.Prisma.ExecuteRaw(`insert into "User" ("id", "email", "username") values ($1,$2,$3)`, "a", "a", "a").Exec(ctx)
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
