package raw

import (
	"context"
	"encoding/json"
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

func (r ExecuteExec) Tx() TxExecuteResult {
	v := TxExecuteResult{}
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

type TxExecuteResult struct {
	query builder.Query
}

func (r TxExecuteResult) ExtractQuery() builder.Query {
	return r.query
}

func (r TxExecuteResult) IsTx() {}

func (r TxExecuteResult) Result() *types.BatchResult {
	var v int
	data, ok := <-r.query.TxResult
	if !ok {
		return nil
	}
	if err := json.Unmarshal(data, &v); err != nil {
		panic(err)
	}
	return &types.BatchResult{
		Count: v,
	}
}
