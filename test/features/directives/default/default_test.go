package enums

import (
	"context"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
	"github.com/steebchen/prisma-client-go/test/helpers/massert"
)

func TestEnums(t *testing.T) {
	test.RunParallel(t, []test.Database{test.MySQL, test.PostgreSQL, test.MongoDB}, func(t *testing.T, db test.Database, ctx context.Context) {
		client := NewClient()

		mockDB := test.Start(t, db, client.Engine, []string{})
		defer test.End(t, db, client.Engine, mockDB)

		stringOpt := "stringOpt"
		boolOpt := false
		testB := TestB
		expected := &UserModel{
			InnerUser: InnerUser{
				ID:        "123",
				String:    "string",
				StringOpt: &stringOpt,
				Bool:      true,
				BoolOpt:   &boolOpt,
				Test:      TestA,
				TestOpt:   &testB,
			},
		}

		created, err := client.User.CreateOne(
			User.ID.Set("123"),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("fail %s", err)
		}

		massert.Equal(t, expected, created)

		actual, err := client.User.FindMany().Exec(ctx)
		if err != nil {
			t.Fatalf("fail %s", err)
		}

		massert.Equal(t, []UserModel{*expected}, actual)
	})
}
