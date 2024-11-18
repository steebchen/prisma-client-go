package raw

import (
	"context"
	"fmt"

	"github.com/steebchen/prisma-client-go/runtime/builder"
)

func (r Raw) RunCommandRaw(cmd interface{}) RunCommandExec {
	return RunCommandExec{
		query: doCommandRaw(r.Engine, "runCommandRaw", fmt.Sprintf("%v", cmd)),
	}
}

type RunCommandExec struct {
	query builder.Query
}

func (r RunCommandExec) ExtractQuery() builder.Query {
	return r.query
}

func (r RunCommandExec) Tx() TxQueryResult {
	v := NewTxQueryResult()
	v.query = r.query
	v.query.TxResult = make(chan []byte, 1)
	return v
}

func (r RunCommandExec) Exec(ctx context.Context, into interface{}) error {
	if err := r.query.Exec(ctx, &into); err != nil {
		return err
	}

	return nil
}
