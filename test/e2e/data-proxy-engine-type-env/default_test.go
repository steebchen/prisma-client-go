//go:build e2e
// +build e2e

// no-op copy the schema here, but use a template prisma schema file so it won't be generated from the common test
// schema generation, but rather generate it here manually in order to
//go:generate cp schema.template.prisma schema.out.prisma
//go:generate sh generate.sh

// This test checks whether the data proxy works with the PRISMA_CLIENT_ENGINE_TYPE=dataproxy being set
package db

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/steebchen/prisma-client-go/test"
	"github.com/steebchen/prisma-client-go/test/helpers/massert"
)

const id = "23230653-a467-47b7-aaf9-98d422da3d9e"

func str(v string) *string {
	return &v
}

func TestE2ERemoteDataProxyEngineTypeEnv(t *testing.T) {
	t.Skip("temporarily paused")
	test.RunSerial(t, []test.Database{test.PostgreSQL}, func(t *testing.T, db test.Database, ctx context.Context) {
		t.Skip("data proxy is unmaintained")

		client := NewClient()
		if err := client.Connect(); err != nil {
			t.Fatalf("fail %s", err)
		}

		createdAt, _ := time.Parse(time.RFC3339, "2021-09-22T09:32:31.706Z")
		updatedAt, _ := time.Parse(time.RFC3339, "2021-09-22T09:32:31.707Z")

		expected := &UserModel{
			InnerUser: InnerUser{
				ID:        id,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				Name:      str("Bertrand"),
				Email:     "Dakota.Waelchi@gmail.com",
			},
		}

		actual, err := client.User.FindUnique(
			User.ID.Equals(id),
		).Exec(ctx)
		if err != nil {
			t.Fatalf("fail %s", err)
		}

		v, _ := json.MarshalIndent(actual, "", "  ")
		log.Printf("data proxy response: %s", v)

		massert.Equal(t, expected, actual)
	})
}
