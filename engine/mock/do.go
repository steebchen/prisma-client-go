package mock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/prisma/prisma-client-go/engine"
)

func (e *Engine) Do(ctx context.Context, payload interface{}, v interface{}) error {
	e.expMu.Lock()
	defer e.expMu.Unlock()

	expectations := *e.expectations

	var n = -1
	for i, e := range expectations {
		req := payload.(engine.GQLRequest)
		if e.Query.Build() == req.Query {
			n = i
			break
		}
	}
	if n == -1 {
		panic("could not find query")
	}
	var retErr error = nil
	switch {
	case expectations[n].Want != nil:
		r, err := json.Marshal(expectations[n].Want)
		if err != nil {
			return fmt.Errorf("error happened at unmarshaling expectation want: %w", err)
		}
		if err := json.Unmarshal(r, &v); err != nil {
			return fmt.Errorf("error happened at marshaling expectation want: %w", err)
		}
	case expectations[n].WantErr != nil:
		retErr = expectations[n].WantErr
	default:
		panic("need to define either Want or WantErr")
	}
	expectations[n].Success = true
	*e.expectations = expectations
	return retErr
}

func (e *Engine) Batch(ctx context.Context, payload interface{}, v interface{}) error {
	// TODO
	panic("TODO")
}
