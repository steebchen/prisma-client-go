package raw

import (
	"context"
	"fmt"

	"github.com/prisma/prisma-client-go/generator/builder"
)

func (r Actions) ExecuteRaw(query string, params ...interface{}) ExecuteExec {
	return ExecuteExec{
		query: raw(r.Client, "executeRaw", query, params...),
	}
}

type ExecuteExec struct {
	query builder.Query
}

type ExecuteResult struct {
	ExecuteRaw int `json:"executeRaw"`
}

func (r ExecuteExec) Exec(ctx context.Context) (int, error) {
	var result ExecuteResult
	if err := r.query.Exec(ctx, &result); err != nil {
		return 0, fmt.Errorf("could not send raw query: %w", err)
	}

	return result.ExecuteRaw, nil
}
