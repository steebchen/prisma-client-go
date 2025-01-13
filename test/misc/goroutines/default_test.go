package db

import (
	"context"
	"testing"
	"time"

	"github.com/steebchen/prisma-client-go/test"
	"github.com/steebchen/prisma-client-go/test/helpers/massert"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

// The purpose of this test is to have a goroutine continuously making queries (whether or not they succeed),
// and detect race conditions, after disconnecting.
func TestDisconnectConcurrent(t *testing.T) {
	test.RunSerial(t, []test.Database{test.MySQL, test.PostgreSQL}, func(t *testing.T, db test.Database, ctx context.Context) {
		client := NewClient()
		mockDBName := test.Start(t, db, client.Engine, []string{})

		a := "a"
		b := "b"
		created, err := client.User.CreateOne(
			User.A.Set(a),
			User.B.Set(b),
			User.ID.Set("123"),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("fail %s", err)
		}

		expected := &UserModel{
			InnerUser: InnerUser{
				ID: "123",
				A:  a,
				B:  &b,
			},
		}

		massert.Equal(t, expected, created)

		// Query database concurrently
		closeCh := make(chan struct{})
		go func() {
			loop := true
			for loop {
				time.Sleep(time.Millisecond * 100)
				select {
				case <-closeCh:
					loop = false
				default:
					_, _ = client.User.FindUnique(User.ID.Equals(created.ID)).Exec(ctx)
				}
			}
		}()

		// Wait, and close the connection
		time.Sleep(time.Millisecond * 300)
		test.End(t, db, client.Engine, mockDBName)
		closeCh <- struct{}{}
	})
}
