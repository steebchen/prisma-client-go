package enums

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

func TestEnums(t *testing.T) {
	test.RunParallel(t, []test.Database{test.MySQL, test.PostgreSQL}, func(t *testing.T, db test.Database, ctx context.Context) {
		client := NewClient()

		mockDB := test.Start(t, db, client.Engine, []string{})
		defer test.End(t, db, client.Engine, mockDB)

		admin := RoleAdmin
		mod := RoleModerator
		expected := &UserModel{
			InnerUser: InnerUser{
				ID:      "123",
				Role:    admin,
				RoleOpt: &mod,
				Stuff1:  StuffCASING,
				Stuff2:  StuffDifferent,
				Stuff3:  StuffHaHa,
			},
		}

		created, err := client.User.CreateOne(
			User.Role.Set(RoleAdmin),
			User.Stuff1.Set(StuffCASING),
			User.Stuff2.Set(StuffDifferent),
			User.Stuff3.Set(StuffHaHa),
			User.ID.Set("123"),
			User.RoleOpt.Set(RoleModerator),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("fail %s", err)
		}

		assert.Equal(t, expected, created)

		actual, err := client.User.FindMany(
			User.Role.Equals(RoleAdmin),
			User.Role.In([]Role{RoleAdmin}),
			User.RoleOpt.Equals(RoleModerator),
			User.RoleOpt.In([]Role{RoleModerator}),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("fail %s", err)
		}

		assert.Equal(t, []UserModel{*expected}, actual)
	})
}
