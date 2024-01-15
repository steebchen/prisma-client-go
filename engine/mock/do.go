package mock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/steebchen/prisma-client-go/engine/protocol"
)

func (e *Engine) Do(_ context.Context, payload interface{}, v interface{}) error {
	e.expMu.Lock()
	defer e.expMu.Unlock()

	expectations := *e.expectations

	n := -1
	for i, e := range expectations {
		req := payload.(protocol.GQLRequest)
		str, err := e.Query.Build()
		if err != nil {
			return err
		}
		if str == req.Query {
			n = i
			break
		}
	}
	if n == -1 {
		panic("could not find query")
	}
	var retErr error
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

func (e *Engine) Batch(context.Context, interface{}, interface{}) error {
	// TODO
	panic("TODO")
}
