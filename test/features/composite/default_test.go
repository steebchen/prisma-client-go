package composite

import (
	"context"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
	"github.com/steebchen/prisma-client-go/test/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestComposite(t *testing.T) {
	// no-op compile time test
	User.SomethingIDAnotherIDStuff(
		User.SomethingID.Equals(""),
		User.AnotherIDStuff.Equals(""),
	)

	// custom name test
	User.AnotherIDStuffSomethingID(
		User.AnotherIDStuff.Equals(""),
		User.SomethingID.Equals(""),
	)

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name:   "composite FirstNameMiddleNameLastName",
		before: nil,
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			expected := &UserModel{
				InnerUser: InnerUser{
					FirstName:  "a",
					MiddleName: "b",
					LastName:   "c",
				},
			}

			user, err := client.User.CreateOne(
				User.FirstName.Set("a"),
				User.MiddleName.Set("b"),
				User.LastName.Set("c"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, expected, user)

			users, err := client.User.FindMany(
				User.FirstName.Equals("a"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, []UserModel{*expected}, users)

			oneUser, err := client.User.FindUnique(
				User.FirstNameMiddleNameLastName(
					User.FirstName.Equals("a"),
					User.MiddleName.Equals("b"),
					User.LastName.Equals("c"),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, expected, oneUser)
		},
	}, {
		name:   "composite FirstNameLastName",
		before: nil,
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			expected := &UserModel{
				InnerUser: InnerUser{
					FirstName:  "a",
					MiddleName: "b",
					LastName:   "c",
				},
			}

			user, err := client.User.CreateOne(
				User.FirstName.Set("a"),
				User.MiddleName.Set("b"),
				User.LastName.Set("c"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, expected, user)

			users, err := client.User.FindMany(
				User.FirstName.Equals("a"),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, []UserModel{*expected}, users)

			oneUser, err := client.User.FindUnique(
				User.FirstNameLastName(
					User.FirstName.Equals("a"),
					User.LastName.Equals("c"),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatal(err)
			}

			massert.Equal(t, expected, oneUser)
		},
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, test.Databases, func(t *testing.T, db test.Database, ctx context.Context) {
				if db == test.MongoDB {
					// TODO
					t.Skip()
				}
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
