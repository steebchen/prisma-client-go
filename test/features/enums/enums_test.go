package enums

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/steebchen/prisma-client-go/runtime/builder"
	"github.com/steebchen/prisma-client-go/test"
	"github.com/steebchen/prisma-client-go/test/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestEnums(t *testing.T) {
	t.Parallel()

	// casing test no-ops
	_ = StuffLast7D
	_ = StuffLast30D
	_ = StuffSlack
	_ = StuffLast7DAnd
	_ = StuffLast30DAnd
	_ = StuffID

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "create",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			admin := RoleAdmin
			mod := RoleModerator

			stuffCasing := StuffCasing
			stuffDifferent := StuffDifferent
			stuffHaHa := StuffHaHa
			expected := &UserModel{
				InnerUser: InnerUser{
					ID:      "123",
					Role:    admin,
					RoleOpt: &mod,
					Stuff1:  &stuffCasing,
					Stuff2:  &stuffDifferent,
					Stuff3:  &stuffHaHa,
				},
			}

			created, err := client.User.CreateOne(
				User.Role.Set(RoleAdmin),
				User.Stuff1.Set(StuffCasing),
				User.Stuff2.Set(StuffDifferent),
				User.Stuff3.Set(StuffHaHa),
				User.ID.Set("123"),
				User.RoleOpt.Set(RoleModerator),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, expected, created)

			actual, err := client.User.FindFirst(
				User.Role.Equals(RoleAdmin),
				User.Role.In([]Role{RoleAdmin}),
				User.RoleOpt.Equals(RoleModerator),
				User.RoleOpt.In([]Role{RoleModerator}),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, expected, actual)
		},
	}, {
		name: "many or with and wrapper",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.CreateOne(
				User.Role.Set(RoleAdmin),
				User.ID.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			_, err = client.User.CreateOne(
				User.Role.Set(RoleModerator),
				User.ID.Set("456"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			_, err = client.User.CreateOne(
				User.Role.Set(RoleUser),
				User.ID.Set("789"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			actual, err := client.User.FindMany(
				User.Or(
					User.And(
						User.Role.Equals(RoleUser),
					),
					User.And(
						User.Role.Equals(RoleAdmin),
					),
				),
			).OrderBy(
				User.ID.Order(SortOrderAsc),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, []UserModel{
				{
					InnerUser: InnerUser{
						ID:   "123",
						Role: RoleAdmin,
					},
				},
				{
					InnerUser: InnerUser{
						ID:   "789",
						Role: RoleUser,
					},
				},
			}, actual)
		},
	}, {
		name: "many or direct",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.CreateOne(
				User.Role.Set(RoleAdmin),
				User.ID.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			_, err = client.User.CreateOne(
				User.Role.Set(RoleModerator),
				User.ID.Set("456"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			_, err = client.User.CreateOne(
				User.Role.Set(RoleUser),
				User.ID.Set("789"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			_, err = client.User.FindMany(
				User.Or(
					User.Role.Equals(RoleUser),
					User.Role.Equals(RoleAdmin),
				),
			).OrderBy(
				User.ID.Order(SortOrderAsc),
			).Exec(ctx)

			assert.Equal(t, builder.ErrDuplicateField, errors.Unwrap(err))
		},
	}, {
		name: "in",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.CreateOne(
				User.Role.Set(RoleAdmin),
				User.ID.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			_, err = client.User.CreateOne(
				User.Role.Set(RoleModerator),
				User.ID.Set("456"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			_, err = client.User.CreateOne(
				User.Role.Set(RoleUser),
				User.ID.Set("789"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			actual, err := client.User.FindMany(
				User.Role.In([]Role{RoleUser, RoleAdmin}),
			).OrderBy(
				User.ID.Order(SortOrderAsc),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			massert.Equal(t, []UserModel{
				{
					InnerUser: InnerUser{
						ID:   "123",
						Role: RoleAdmin,
					},
				},
				{
					InnerUser: InnerUser{
						ID:   "789",
						Role: RoleUser,
					},
				},
			}, actual)
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
