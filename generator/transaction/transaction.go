package transaction

import (
	"context"
	"fmt"

	"github.com/prisma/prisma-client-go/engine"
	"github.com/prisma/prisma-client-go/generator/builder"
)

type TX struct {
	Engine engine.Engine
}

type Params interface {
	ExtractQuery() builder.Query
}

func (r TX) Transaction(queries ...Params) Exec {
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
	}
}

type Exec struct {
	engine   engine.Engine
	requests []engine.GQLRequest
}

func (r Exec) Exec(ctx context.Context) error {
	var result engine.GQLBatchResponse
	payload := engine.GQLBatchRequest{
		Batch:       r.requests,
		Transaction: true,
	}
	if err := r.engine.Batch(ctx, payload, &result); err != nil {
		return fmt.Errorf("could not send raw query: %w", err)
	}
	return nil
}
