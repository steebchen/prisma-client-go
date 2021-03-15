package raw

import (
	"github.com/prisma/prisma-client-go/runtime/builder"
)

type txExec struct {
	query builder.Query
}

func (r txExec) Tx() *txResult {
	v := &txResult{}
	v.query = r.query
	v.query.TxResult = make(chan []byte, 1)
	return v
}

type txResult struct {
	query builder.Query
}

func (r *txResult) ExtractQuery() builder.Query {
	return r.query
}

func (r *txResult) IsTx() {}
