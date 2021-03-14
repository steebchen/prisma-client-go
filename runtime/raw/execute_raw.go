package raw

import (
	"context"
	"fmt"

	"github.com/prisma/prisma-client-go/runtime/builder"
	"github.com/prisma/prisma-client-go/runtime/types"
)

func (r Raw) ExecuteRaw(query string, params ...interface{}) ExecuteExec {
	return ExecuteExec{
		query: raw(r.Engine, "executeRaw", query, params...),
	}
}

type ExecuteExec struct {
	query builder.Query
}

func (r ExecuteExec) ExtractQuery() builder.Query {
	return r.query
}

func (r ExecuteExec) Exec(ctx context.Context) (*types.BatchResult, error) {
	var count int
	if err := r.query.Exec(ctx, &count); err != nil {
		return nil, fmt.Errorf("could not send raw query: %w", err)
	}
	return &types.BatchResult{
		Count: count,
	}, nil
}
