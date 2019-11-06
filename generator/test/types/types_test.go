package types

//go:generate prisma2 generate

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type cx = context.Context
type Func func(t *testing.T, client Client, ctx cx)

func cmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		exit, ok := err.(*exec.ExitError)
		if !ok {
			return fmt.Errorf("command %s %s failed: %w", name, args, err)
		}

		if !exit.Success() {
			return fmt.Errorf("%s %s exited with status code %d and output %s: %w", name, args, exit.ExitCode(), string(out), err)
		}
	}

	return nil
}

func TestTypes(t *testing.T) {
	t.Parallel()

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
				User.ID.Set(id),
				User.Str.Set(str),
				User.Int.Set(5),
				User.Float.Set(5.5),
				User.Bool.Set(true),
				User.Date.Set(date),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := UserModel{
				user{
					ID:    id,
					Str:   str,
					Int:   5,
					Float: 5.5,
					Bool:  true,
					Date:  date,
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
		name: "basic equals",
		// language=GraphQL
		before: `
			mutation {
				a: createOneUser(data: {
					id: "id",
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
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				user{
					ID:    "id",
					Str:   "str",
					Int:   5,
					Float: 5.5,
					Bool:  true,
					Date:  date,
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
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				user{
					ID:    "id",
					Str:   "alongstring",
					Int:   5,
					Float: 5.5,
					Bool:  true,
					Date:  date,
				},
			}}

			assert.Equal(t, expected, users)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := cmd("rm", "-rf", "dev.sqlite"); err != nil {
				log.Fatal(err)
			}
			if err := cmd("rm", "-rf", "migrations"); err != nil {
				log.Fatal(err)
			}

			if err := cmd("prisma2", "lift", "save", "--create-db", "--name", "init"); err != nil {
				t.Fatalf("could not run lift save %s", err)
			}
			if err := cmd("prisma2", "lift", "up"); err != nil {
				t.Fatalf("could not run lift up %s", err)
			}

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
				response, err := client.gql.Raw(ctx, tt.before, map[string]interface{}{})
				log.Printf("mock response query %+v", response)
				if err != nil {
					t.Fatalf("could not send mock query %s %+v", err, response)
				}
			}

			tt.run(t, client, ctx)
		})
	}
}
