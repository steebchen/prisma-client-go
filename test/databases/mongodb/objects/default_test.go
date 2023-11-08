package raw

import (
	"context"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
	"github.com/steebchen/prisma-client-go/test/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestObjects(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "types",
		// language=GraphQL
		before: []string{`
			mutation {
				result: createOneUser(data: {
					id: "id1",
					email: "email1",
					info: {
						age: 0,
					},
					infoOpt: {
						age: 0,
					},
					list: {
						create: [
							{
								age: 0,
							},
						]
					}
				}) {
					id
				}
			}
		`},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			user, err := client.User.FindUnique(
				User.ID.Equals("id1"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, &UserModel{
				InnerUser: InnerUser{
					ID:    "id1",
					Email: "email1",
					Info: InfoType{
						Age: 0,
					},
					InfoOpt: &InfoType{
						Age: 0,
					},
					List: []InfoType{
						{
							Age: 0,
						},
					},
				},
			}, user)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()

			mockDB := test.Start(t, test.MongoDB, client.Engine, tt.before)
			defer test.End(t, test.MongoDB, client.Engine, mockDB)

			tt.run(t, client, context.Background())
		})
	}
}
