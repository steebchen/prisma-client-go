package raw

import (
	"context"
	"fmt"

	"github.com/prisma/prisma-client-go/runtime/builder"
	"github.com/prisma/prisma-client-go/runtime/transaction"
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

func (r ExecuteExec) Tx() TxExecuteResult {
	v := NewTxExecuteResult()
	v.query = r.query
	v.query.TxResult = make(chan []byte, 1)
	return v
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

func NewTxExecuteResult() TxExecuteResult {
	return TxExecuteResult{
		result: &transaction.Result{},
	}
}

type TxExecuteResult struct {
	query  builder.Query
	result *transaction.Result
}

func (r TxExecuteResult) ExtractQuery() builder.Query {
	return r.query
}

func (r TxExecuteResult) IsTx() {}

func (r TxExecuteResult) Result() *types.BatchResult {
	var v int
	if err := r.result.Get(r.query.TxResult, &v); err != nil {
		panic(err)
	}
	return &types.BatchResult{
		Count: v,
	}
}
