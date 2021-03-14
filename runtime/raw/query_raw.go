package raw

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/prisma/prisma-client-go/runtime/builder"
)

func (r Raw) QueryRaw(query string, params ...interface{}) QueryExec {
	return QueryExec{
		query: raw(r.Engine, "queryRaw", query, params...),
	}
}

type QueryExec struct {
	query builder.Query
}

func (r QueryExec) ExtractQuery() builder.Query {
	return r.query
}

type QueryResult struct {
	QueryRaw json.RawMessage `json:"queryRaw"`
}

func (r QueryExec) Exec(ctx context.Context, into interface{}) error {
	if err := r.query.Exec(ctx, &into); err != nil {
		return fmt.Errorf("could not send raw query: %w", err)
	}

	return nil
}
