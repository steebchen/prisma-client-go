package mock

import (
	"testing"

	"github.com/steebchen/prisma-client-go/runtime/builder"
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
			str, err := e.Query.Build()
			if err != nil {
				t.Fatalf("could not build query: %s", err)
			}
			t.Fatalf("expectation not met for query `%s` and result `%s`, error `%s`", str, e.Want, e.WantErr)
		}
	}
}
