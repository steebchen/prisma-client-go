package raw

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/prisma/prisma-client-go/runtime/builder"
	"github.com/prisma/prisma-client-go/runtime/transaction"
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

func (r QueryExec) Tx() TxQueryResult {
	v := NewTxQueryResult()
	v.query = r.query
	v.query.TxResult = make(chan []byte, 1)
	return v
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

func NewTxQueryResult() TxQueryResult {
	return TxQueryResult{
		result: &transaction.Result{},
	}
}

type TxQueryResult struct {
	query  builder.Query
	result *transaction.Result
}

func (r TxQueryResult) ExtractQuery() builder.Query {
	return r.query
}

func (r TxQueryResult) IsTx() {}

func (r TxQueryResult) Into(v interface{}) error {
	return r.result.Get(r.query.TxResult, &v)
}
