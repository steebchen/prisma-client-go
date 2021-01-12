package raw

import (
	"context"
	"fmt"

	"github.com/prisma/prisma-client-go/generator/builder"
)

func (r Actions) ExecuteRaw(query string, params ...interface{}) ExecuteExec {
	return ExecuteExec{
		query: raw(r.Engine, "executeRaw", query, params...),
	}
}

type ExecuteExec struct {
	query builder.Query
}

func (r ExecuteExec) Exec(ctx context.Context) (int, error) {
	var result int
	if err := r.query.Exec(ctx, &result); err != nil {
		return -1, fmt.Errorf("could not send raw query: %w", err)
	}
	return result, nil
}
