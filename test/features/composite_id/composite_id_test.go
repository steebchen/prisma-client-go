package db

import (
	"context"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
)

// TODO add test for find unique by composite id

func TestCompositeID(t *testing.T) {
	test.RunParallel(t, []test.Database{test.MySQL, test.PostgreSQL}, func(t *testing.T, db test.Database, ctx context.Context) {
		client := NewClient()

		mockDB := test.Start(t, db, client.Engine, []string{})
		defer test.End(t, db, client.Engine, mockDB)

		_, err := client.Repository.CreateOne(
			Repository.PlatformID.Set("test"),
			Repository.PlatformKind.Set("test"),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("fail %s", err)
		}

		_, err = client.RepositoryOrganization.CreateOne(
			RepositoryOrganization.PlatformID.Set("test"),
			RepositoryOrganization.PlatformKind.Set("test"),
			RepositoryOrganization.Name.Set("test"),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("fail %s", err)
		}
	})
}
