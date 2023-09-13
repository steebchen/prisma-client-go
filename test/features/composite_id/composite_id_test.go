package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func str(v string) *string {
	return &v
}

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

			_, err = client.Organization.CreateOne(
				Organization.PlatformID.Set("test"),
				Organization.PlatformKind.Set("test"),
				Organization.Name.Set("test"),
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

			_, err = client.Organization.CreateOne(
				Organization.PlatformID.Set("test"),
				Organization.PlatformKind.Set("test"),
				Organization.Name.Set("test"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			org, err := client.Organization.FindUnique(
				Organization.OrganizationID(
					Organization.PlatformKind.Equals("test"),
					Organization.PlatformID.Equals("test"),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := &OrganizationModel{
				InnerOrganization: InnerOrganization{
					PlatformID:   "test",
					PlatformKind: "test",
					Name:         "test",
				},
			}
			assert.Equal(t, expected, org)
		},
	}, {
		name: "link",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.Repository.CreateOne(
				Repository.PlatformID.Set("test"),
				Repository.PlatformKind.Set("test"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			_, err = client.Organization.CreateOne(
				Organization.PlatformID.Set("a"),
				Organization.PlatformKind.Set("kind"),
				Organization.Name.Set("a"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			_, err = client.Organization.CreateOne(
				Organization.PlatformID.Set("b"),
				Organization.PlatformKind.Set("kind"),
				Organization.Name.Set("b"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			_, err = client.Repository.FindUnique(
				Repository.RepositoryID(
					Repository.PlatformKind.Equals("test"),
					Repository.PlatformID.Equals("test"),
				),
			).Update(
				Repository.Org.Link(
					Organization.OrganizationID(
						Organization.PlatformKind.Equals("kind"),
						Organization.PlatformID.Equals("b"),
					),
				),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			org, err := client.Organization.FindMany().With(
				Organization.Repositories.Fetch(),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			expected := []OrganizationModel{{
				InnerOrganization: InnerOrganization{
					PlatformID:   "a",
					PlatformKind: "kind",
					Name:         "a",
				},
				RelationsOrganization: RelationsOrganization{
					Repositories: []RepositoryModel{},
				},
			}, {
				InnerOrganization: InnerOrganization{
					PlatformID:   "b",
					PlatformKind: "kind",
					Name:         "b",
				},
				RelationsOrganization: RelationsOrganization{
					Repositories: []RepositoryModel{
						{
							InnerRepository: InnerRepository{
								PlatformID:   "test",
								PlatformKind: "kind",
								OrgID:        str("b"),
							},
						},
					},
				},
			}}
			assert.Equal(t, expected, org)
		},
	}, {
		name: "create with specific model layout",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.Company.CreateOne(
				Company.ID.Set("123"),
				Company.Name.Set("name"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			_, err = client.Access.CreateOne(
				Access.CompanyRelation.Link(
					Company.ID.Equals("123"),
				),
				Access.Email.Set("email"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}
		},
	}, {
		name: "create with another specific model layout",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			_, err := client.Team.CreateOne(
				Team.Path.Set("123"),
				Team.Name.Set("name"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}

			_, err = client.Document.CreateOne(
				Document.TeamRelation.Link(
					Team.Path.Equals("123"),
				),
				Document.Name.Set("email"),
				Document.Type.Set("123"),
			).Exec(ctx)
			if err != nil {
				t.Fatalf("fail %s", err)
			}
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
