package composite

import (
	"context"
	"github.com/prisma/prisma-client-go/test"
	"testing"
)

func TestFindManyRelationUpdate(t *testing.T) {
	t.Logf("success (compiling means succeeded)")

	test.RunParallel(t, []test.Database{test.MySQL, test.PostgreSQL, test.SQLite}, func(t *testing.T, db test.Database, ctx context.Context) {
		client := NewClient()
		mockDBName := test.Start(t, test.SQLite, client.Engine, []string{})
		defer test.End(t, test.SQLite, client.Engine, mockDBName)

		_, err := client.Post.FindMany(
			Post.ID.Equals("ckfxy9dh00003c5rp97l1cica"),
			//Post.UpdatedAt.Before(timestamp),
		).Update(
			Post.User.Link(User.ID.Equals("123")),
		).Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})
}
