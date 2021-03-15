package raw

import (
	"context"
	"fmt"

	"github.com/prisma/prisma-client-go/runtime/builder"
	"github.com/prisma/prisma-client-go/runtime/types"
)

func (r Raw) ExecuteRaw(query string, params ...interface{}) ExecuteExec {
	q := raw(r.Engine, "executeRaw", query, params...)
	return ExecuteExec{
		query: q,
		txExec: txExec{
			query: q,
		},
	}
}

type ExecuteExec struct {
	query builder.Query
	txExec
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
