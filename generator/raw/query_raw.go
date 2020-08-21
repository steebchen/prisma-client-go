package raw

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/prisma/prisma-client-go/generator/builder"
)

func (r Actions) QueryRaw(query string, params ...interface{}) QueryExec {
	return QueryExec{
		query: raw(r.Client, "queryRaw", query, params...),
	}
}

type QueryExec struct {
	query builder.Query
}

type QueryResult struct {
	Data struct {
		QueryRaw json.RawMessage `json:"queryRaw"`
	} `json:"data"`
}

func (r QueryExec) Exec(ctx context.Context, into interface{}) error {
	var result QueryResult
	if err := r.query.Exec(ctx, &result); err != nil {
		return fmt.Errorf("could not send raw query: %w", err)
	}

	if err := json.Unmarshal(result.Data.QueryRaw, into); err != nil {
		return fmt.Errorf("could not decode result.QueryRaw: %w", err)
	}

	return nil
}
