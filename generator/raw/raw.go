package raw

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/prisma/prisma-client-go/generator/builder"
)

type Actions struct {
	Client builder.Client
}

func (r Actions) Raw(query string, params ...interface{}) Exec {
	q := builder.NewQuery()
	q.Client = r.Client
	q.Operation = "mutation"
	q.Method = "executeRaw"

	q.Inputs = append(q.Inputs, builder.Input{
		Name:  "query",
		Value: query,
	})

	// convert params to a string with an array
	var newParams = "["
	for i, param := range params {
		if i > 0 {
			newParams += ","
		}
		newParams += string(builder.Value(param))
	}
	newParams += "]"

	if len(newParams) > 0 {
		q.Inputs = append(q.Inputs, builder.Input{
			Name:  "parameters",
			Value: newParams,
		})
	}

	return Exec{
		query: q,
	}
}

type Exec struct {
	query builder.Query
}

type Result struct {
	Data struct {
		ExecuteRaw json.RawMessage `json:"executeRaw"`
	} `json:"data"`
}

func (r Exec) Exec(ctx context.Context, into interface{}) error {
	var result Result
	err := r.query.Exec(ctx, &result)
	if err != nil {
		return fmt.Errorf("could not send raw query: %w", err)
	}

	if err := json.Unmarshal(result.Data.ExecuteRaw, into); err != nil {
		return fmt.Errorf("could not decode result.ExecuteRaw: %w", err)
	}

	return nil
}
