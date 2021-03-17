package transaction

import (
	"context"
	"fmt"

	"github.com/prisma/prisma-client-go/engine"
	"github.com/prisma/prisma-client-go/runtime/builder"
)

type TX struct {
	Engine engine.Engine
}

type Param interface {
	IsTx()
	ExtractQuery() builder.Query
}

func (r TX) Transaction(queries ...Param) Exec {
	requests := make([]engine.GQLRequest, len(queries))
	for i, query := range queries {
		requests[i] = engine.GQLRequest{
			Query:     query.ExtractQuery().Build(),
			Variables: map[string]interface{}{},
		}
	}
	return Exec{
		engine:   r.Engine,
		requests: requests,
		queries:  queries,
	}
}

type Exec struct {
	queries  []Param
	engine   engine.Engine
	requests []engine.GQLRequest
}

func (r Exec) Exec(ctx context.Context) error {
	for _, q := range r.queries {
		//goland:noinspection GoDeferInLoop
		defer close(q.ExtractQuery().TxResult)
	}

	var result engine.GQLBatchResponse
	payload := engine.GQLBatchRequest{
		Batch:       r.requests,
		Transaction: true,
	}
	if err := r.engine.Batch(ctx, payload, &result); err != nil {
		return fmt.Errorf("could not send raw query: %w", err)
	}
	if len(result.Errors) > 0 {
		first := result.Errors[0]
		return fmt.Errorf("pql error: %s", first.Message)
	}
	for i, inner := range result.Result {
		if len(inner.Errors) > 0 {
			first := result.Errors[0]
			return fmt.Errorf("pql error: %s", first.Message)
		}

		r.queries[i].ExtractQuery().TxResult <- inner.Data.Result
	}
	return nil
}
