package db

import (
	"context"
	"testing"

	"github.com/prisma/prisma-client-go/test"
)

func TestConflict(t *testing.T) {
	test.RunParallel(t, []test.Database{test.SQLite}, func(t *testing.T, db test.Database, ctx context.Context) {
		client := NewClient()
		mockDBName := test.Start(t, db, client.Engine, []string{})
		defer test.End(t, db, client.Engine, mockDBName)

		// noop, just test for conflicts in code generation
	})
}
