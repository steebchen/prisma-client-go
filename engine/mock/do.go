package mock

import (
	"context"
	"encoding/json"
	"fmt"
)

func (e *Engine) Do(ctx context.Context, query string, v interface{}) error {
	expectations := *e.expectations

	var n = -1
	for i, e := range expectations {
		if e.Query.Build() == query {
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
