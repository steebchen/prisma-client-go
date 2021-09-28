//go:build e2e
// +build e2e

package db

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/prisma/prisma-client-go/test"
)

const id = "23230653-a467-47b7-aaf9-98d422da3d9e"

func str(v string) *string {
	return &v
}

func TestE2ERemoteDataProxy(t *testing.T) {
	test.RunSerial(t, []test.Database{test.PostgreSQL}, func(t *testing.T, db test.Database, ctx context.Context) {
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

		assert.Equal(t, expected, actual)
	})
}
