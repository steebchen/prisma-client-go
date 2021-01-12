package mock

import (
	"testing"

	"github.com/prisma/prisma-client-go/runtime/builder"
)

type Expectation struct {
	Query   builder.Query
	Want    interface{}
	WantErr error
	Success bool
}

type Query interface {
	extractQuery() builder.Query
}

type Mock struct {
	Expectations *[]Expectation
}

func (m *Mock) Ensure(t *testing.T) {
	if len(*m.Expectations) == 0 {
		t.Fatalf("no expectations defined")
	}
	for _, e := range *m.Expectations {
		if !e.Success {
			t.Fatalf("expectation not met for query `%s` and result `%s`, error `%s`", e.Query.Build(), e.Want, e.WantErr)
		}
	}
}
