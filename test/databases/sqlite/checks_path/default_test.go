package raw

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/steebchen/prisma-client-go/test"
)

type cx = context.Context
type Func func(t *testing.T, client *PrismaClient, ctx cx)

func TestSqliteChecks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		before []string
		run    Func
	}{{
		name: "check for connection URL",
		run: func(t *testing.T, client *PrismaClient, ctx cx) {
			assert.Equal(t, "sqlite:dev.db", schemaConnectionURL)
		},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()

			mockDB := test.Start(t, test.SQLite, client.Engine, tt.before)
			defer test.End(t, test.SQLite, client.Engine, mockDB)

			tt.run(t, client, context.Background())
		})
	}
}
