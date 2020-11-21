package mock

import (
	"testing"

	"github.com/prisma/prisma-client-go/generator/builder"
)

type Expectation struct {
	Query   builder.Query
	Want    interface{}
	WantErr error
	Success bool
}

type Query interface {
	ExtractQuery() builder.Query
}

type Mock struct {
	Expectations *[]Expectation
}

func (m *Mock) Expect(query Query) *Exec {
	return &Exec{
		mock:  m,
		query: query.ExtractQuery(),
	}
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

type Exec struct {
	mock  *Mock
	query builder.Query
}

func (m *Exec) Returns(v interface{}) {
	*m.mock.Expectations = append(*m.mock.Expectations, Expectation{
		Query: m.query,
		Want:  &v,
	})
}

func (m *Exec) Errors(err error) {
	*m.mock.Expectations = append(*m.mock.Expectations, Expectation{
		Query:   m.query,
		WantErr: err,
	})
}
