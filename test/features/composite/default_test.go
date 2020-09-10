package composite

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

func TestComposite(t *testing.T) {
	test.RunParallel(t, []test.Database{test.MySQL, test.PostgreSQL, test.SQLite}, func(t *testing.T, db test.Database, ctx context.Context) {
		client := NewClient()

		mockDB := test.Start(t, db, client.Engine, []string{})
		defer test.End(t, db, client.Engine, mockDB)

		expected := UserModel{
			InternalUser: InternalUser{
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

		assert.Equal(t, expected, user)

		users, err := client.User.FindMany(
			User.FirstName.Equals("a"),
		).Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, []UserModel{expected}, users)

		oneUser, err := client.User.FindOne(
			User.FirstNameMiddleNameLastName(
				User.FirstName.Equals("a"),
				User.MiddleName.Equals("b"),
				User.LastName.Equals("c"),
			),
		).Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expected, oneUser)
	})
}
