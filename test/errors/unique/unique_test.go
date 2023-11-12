package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/steebchen/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestUniqueConstraintViolation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		dbs    []test.Database
		before []string
		run    Func
	}{{
		name: "postgres unique constraint violation",
		dbs:  []test.Database{test.PostgreSQL},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.CreateOne(
				User.Email.Set("john@example.com"),
				User.Username.Set("username"),
			).Exec(ctx)
			assert.Equal(t, nil, err)

			_, err = client.User.CreateOne(
				User.Email.Set("john@example.com"),
				User.Username.Set("username"),
			).Exec(ctx)

			violation, ok := IsErrUniqueConstraint(err)
			//	assert.Equal(t, &ErrUniqueConstraint{
			//		Field: User.Email.Field(),
			//	}, violation)
			assert.Equal(t, User.Email.Field(), violation.Fields[0])

			assert.Equal(t, true, ok)
		},
	}, {
		name: "mysql unique constraint violation",
		dbs:  []test.Database{test.MySQL},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.CreateOne(
				User.Email.Set("john@example.com"),
				User.Username.Set("username"),
			).Exec(ctx)
			assert.Equal(t, nil, err)

			_, err = client.User.CreateOne(
				User.Email.Set("john@example.com"),
				User.Username.Set("username"),
			).Exec(ctx)

			violation, ok := IsErrUniqueConstraint(err)
			//	assert.Equal(t, &ErrUniqueConstraint{
			//		Key: "User_email_key",
			//	}, violation)
			assert.Equal(t, "User_email_key", violation.Key)

			assert.Equal(t, true, ok)
		},
	}, {
		name: "sqlite unique constraint violation",
		dbs:  []test.Database{test.SQLite},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.CreateOne(
				User.Email.Set("john@example.com"),
				User.Username.Set("username"),
			).Exec(ctx)
			assert.Equal(t, nil, err)

			_, err = client.User.CreateOne(
				User.Email.Set("john@example.com"),
				User.Username.Set("username"),
			).Exec(ctx)

			violation, ok := IsErrUniqueConstraint(err)
			//	assert.Equal(t, &ErrUniqueConstraint{
			//		Field: User.Email.Field(),
			//	}, violation)
			assert.Equal(t, User.Email.Field(), violation.Fields[0])

			assert.Equal(t, true, ok)
		},
	}, {
		name: "mongodb unique constraint violation",
		dbs:  []test.Database{test.MongoDB},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.CreateOne(
				User.Email.Set("john@example.com"),
				User.Username.Set("username"),
			).Exec(ctx)
			assert.Equal(t, nil, err)

			_, err = client.User.CreateOne(
				User.Email.Set("john@example.com"),
				User.Username.Set("username"),
			).Exec(ctx)

			violation, ok := IsErrUniqueConstraint(err)
			//	assert.Equal(t, &ErrUniqueConstraint{
			//		Key: "User_email_key",
			//	}, violation)
			assert.Equal(t, "User_email_key", violation.Key)

			assert.Equal(t, true, ok)
		},
	}, {
		name: "nil error should succeed",
		dbs:  []test.Database{test.SQLite},
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.User.CreateOne(
				User.Email.Set("john@example.com"),
				User.Username.Set("username"),
			).Exec(ctx)

			_, ok := IsErrUniqueConstraint(err)

			assert.Equal(t, false, ok)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, tt.dbs, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
