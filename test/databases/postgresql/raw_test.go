package postgresql

//go:generate go run github.com/prisma/prisma-client-go generate

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test/hooks"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

const containerName = "go-client-postgres"

func setup(t *testing.T) {
	teardown(t)

	if err := hooks.Cmd("docker", "stop", "postgres"); err != nil {
		log.Println(err)
	}

	if err := hooks.Cmd("docker", "rm", "postgres", "--force"); err != nil {
		log.Println(err)
	}

	if err := hooks.Cmd("docker", "run", "--name", "postgres", "-p", "5432:5432", "-e", "POSTGRES_PASSWORD=pw", "-d", "postgres"); err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)
}

func teardown(t *testing.T) {
	if err := hooks.Cmd("docker", "stop", containerName); err != nil {
		log.Println(err)
	}

	if err := hooks.Cmd("docker", "rm", containerName, "--force"); err != nil {
		log.Println(err)
	}
}

func TestRaw(t *testing.T) {
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
				a: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "a",
				}) {
					id
				}
			}
		`, `
			mutation {
				b: createOneUser(data: {
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
			err := client.Raw(`SELECT * FROM "User"`).Exec(ctx, &actual)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				RawUser: RawUser{
					ID:       "id1",
					Email:    "email1",
					Username: "a",
				},
			}, {
				RawUser: RawUser{
					ID:       "id2",
					Email:    "email2",
					Username: "b",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "raw query with parameter",
		// language=GraphQL
		before: []string{`
			mutation {
				a: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "a",
				}) {
					id
				}
			}
		`, `
			mutation {
				b: createOneUser(data: {
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
			err := client.Raw(`SELECT * FROM "User" WHERE id = $1`, "id2").Exec(ctx, &actual)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				RawUser: RawUser{
					ID:       "id2",
					Email:    "email2",
					Username: "b",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}, {
		name: "raw query with multiple parameters",
		// language=GraphQL
		before: []string{`
			mutation {
				a: createOneUser(data: {
					id: "id1",
					email: "email1",
					username: "a",
				}) {
					id
				}
			}
		`, `
			mutation {
				b: createOneUser(data: {
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
			err := client.Raw(`SELECT * FROM "User" WHERE id = $1 AND email = $2`, "id2", "email2").Exec(ctx, &actual)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []UserModel{{
				RawUser: RawUser{
					ID:       "id2",
					Email:    "email2",
					Username: "b",
				},
			}}

			assert.Equal(t, expected, actual)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setup(t)

			client := NewClient()

			hooks.Start(t, client.Engine, tt.before)
			defer hooks.End(t, client.Engine)

			tt.run(t, client, context.Background())

			// teardown(t)
		})
	}
}
