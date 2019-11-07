package relations

//go:generate prisma2 generate

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/photongo/generator/test/hooks"
)

type cx = context.Context
type Func func(t *testing.T, client Client, ctx cx)

func str(v string) *string {
	return &v
}

func TestRelations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before string
		run    Func
	}{{
		name: "relations",
		// language=GraphQL
		before: `
			mutation {
				createOneUser(data: {
					id: "relations",
					email: "john@example.com",
					username: "johndoe",
					name: "John",
					stuff: null,
				}) {
					id
				}
			}
		`,
		run: func(t *testing.T, client Client, ctx cx) {
			actual, err := client.Post.FindMany(
				Post.Title.Equals("asdf"),
				Post.User().Email.Equals("john@example.com"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			assert.Equal(t, "John", actual[0].ID)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			hooks.Run(t)

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
