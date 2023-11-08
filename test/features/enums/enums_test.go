package enums

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/steebchen/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestEnums(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "create",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			admin := RoleAdmin
			mod := RoleModerator

			stuffCASING := StuffCASING
			stuffDifferent := StuffDifferent
			stuffHaHa := StuffHaHa
			expected := &UserModel{
				InnerUser: InnerUser{
					ID:      "123",
					Role:    admin,
					RoleOpt: &mod,
					Stuff1:  &stuffCASING,
					Stuff2:  &stuffDifferent,
					Stuff3:  &stuffHaHa,
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
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, []test.Database{test.MySQL, test.PostgreSQL, test.MongoDB}, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
