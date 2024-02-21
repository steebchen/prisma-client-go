package transaction

import (
	"context"
	"fmt"

	"github.com/steebchen/prisma-client-go/engine"
	"github.com/steebchen/prisma-client-go/engine/protocol"
	"github.com/steebchen/prisma-client-go/runtime/builder"
)

type TX struct {
	Engine engine.Engine
}

// Deprecated: use Transaction instead
type Param = Transaction

type Transaction interface {
	IsTx()
	ExtractQuery() builder.Query
}

func (r TX) Transaction(queries ...Transaction) Exec {
	return Exec{
		engine:  r.Engine,
		queries: queries,
	}
}

type Exec struct {
	queries  []Transaction
	engine   engine.Engine
	requests []protocol.GQLRequest
}

func (r Exec) Exec(ctx context.Context) error {
	r.requests = make([]protocol.GQLRequest, len(r.queries))
	for i, query := range r.queries {
		str, err := query.ExtractQuery().Build()
		if err != nil {
			return err
		}
		r.requests[i] = protocol.GQLRequest{
			Query:     str,
			Variables: map[string]interface{}{},
		}
	}

	for _, q := range r.queries {
		//goland:noinspection GoDeferInLoop
		defer close(q.ExtractQuery().TxResult)
	}

	var result protocol.GQLBatchResponse
	payload := protocol.GQLBatchRequest{
		Batch:       r.requests,
		Transaction: true,
	}
	if err := r.engine.Batch(ctx, payload, &result); err != nil {
		return fmt.Errorf("could not send raw query: %w", err)
	}
	if len(result.Errors) > 0 {
		first := result.Errors[0]
		return fmt.Errorf("pql error: %s", first.RawMessage())
	}
	for i, inner := range result.Result {
		if len(inner.Errors) > 0 {
			first := result.Errors[0]
			return fmt.Errorf("pql error: %s", first.RawMessage())
		}

		r.queries[i].ExtractQuery().TxResult <- inner.Data.Result
	}
	return nil
}
