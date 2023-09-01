package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestCompositeID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "create",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
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
		},
	}, {
		name: "find",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
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

			org, err := client.RepositoryOrganization.FindUnique(
				RepositoryOrganization.RepositoryOrganizationID(
					RepositoryOrganization.PlatformKind.Equals("test"),
					RepositoryOrganization.PlatformID.Equals("test"),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &RepositoryOrganizationModel{
				InnerRepositoryOrganization: InnerRepositoryOrganization{
					PlatformID:   "test",
					PlatformKind: "test",
					Name:         "test",
				},
			}
			assert.Equal(t, expected, org)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			test.RunSerial(t, []test.Database{test.MySQL, test.PostgreSQL, test.SQLite}, func(t *testing.T, db test.Database, ctx context.Context) {
				client := NewClient()
				mockDBName := test.Start(t, db, client.Engine, tt.before)
				defer test.End(t, db, client.Engine, mockDBName)
				tt.run(t, client, context.Background())
			})
		})
	}
}
