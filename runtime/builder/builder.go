package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/prisma/prisma-client-go/engine"
	"github.com/prisma/prisma-client-go/logger"
)

type Input struct {
	Name   string
	Fields []Field
	Value  interface{}
}

// Output can be a single Name or can have nested fields
type Output struct {
	Name string

	// Inputs (optional) to provide arguments to a field
	Inputs []Input

	Outputs []Output
}

type Field struct {
	// The Name of the field.
	Name string

	// List saves whether the fields is a list of items
	List bool

	// WrapList saves whether the a list field should be wrapped in an object
	WrapList bool

	// Value contains the field value. if nil, fields will contain a subselection.
	Value interface{}

	// Fields contains a subselection of fields. If not nil, value will be undefined.
	Fields []Field
}

func NewQuery() Query {
	return Query{
		Start: time.Now(),
	}
}

type Query struct {
	// Engine holds the implementation of how queries are processed
	Engine engine.Engine

	// Operation describes the PQL operation: query, mutation or subscription
	Operation string

	// Name describes the operation; useful for tracing
	Name string

	// Method describes a crud operation
	Method string

	// Model contains the Prisma model Name
	Model string

	// Inputs contains function arguments
	Inputs []Input

	// Outputs contains the return fields
	Outputs []Output

	// Start time of the request for tracing
	Start time.Time

	TxResult chan []byte
}

func (q Query) Build() string {
	var builder strings.Builder

	builder.WriteString(q.Operation + " " + q.Name)
	builder.WriteString("{")
	builder.WriteString("result: ")

	builder.WriteString(q.BuildInner())

	builder.WriteString("}")

	return builder.String()
}

func (q Query) BuildInner() string {
	var builder strings.Builder

	builder.WriteString(q.Method + q.Model)

	if len(q.Inputs) > 0 {
		builder.WriteString(q.buildInputs(q.Inputs))
	}

	builder.WriteString(" ")

	if len(q.Outputs) > 0 {
		builder.WriteString(q.buildOutputs(q.Outputs))
	}

	return builder.String()
}

func (q Query) buildInputs(inputs []Input) string {
	var builder strings.Builder

	builder.WriteString("(")

	for _, i := range inputs {
		builder.WriteString(i.Name)

		builder.WriteString(":")

		if i.Value != nil {
			builder.Write(Value(i.Value))
		} else {
			builder.WriteString(q.buildFields(false, false, i.Fields))
		}

		builder.WriteString(",")
	}

	builder.WriteString(")")

	return builder.String()
}

func (q Query) buildOutputs(outputs []Output) string {
	var builder strings.Builder

	builder.WriteString("{")

	for _, o := range outputs {
		builder.WriteString(o.Name + " ")

		if len(o.Inputs) > 0 {
			builder.WriteString(q.buildInputs(o.Inputs))
		}

		if len(o.Outputs) > 0 {
			builder.WriteString(q.buildOutputs(o.Outputs))
		}
	}

	builder.WriteString("}")

	return builder.String()
}

func (q Query) buildFields(list bool, wrapList bool, fields []Field) string {
	var builder strings.Builder

	if !list {
		builder.WriteString("{")
	}

	for _, f := range fields {
		if wrapList {
			builder.WriteString("{")
		}

		if f.Name != "" {
			builder.WriteString(f.Name)
		}

		if f.Name != "" {
			builder.WriteString(":")
		}

		if f.List {
			builder.WriteString("[")
		}

		if f.Fields != nil {
			builder.WriteString(q.buildFields(f.List, f.WrapList, f.Fields))
		}

		if f.Value != nil {
			builder.Write(Value(f.Value))
		}

		if f.List {
			builder.WriteString("]")
		}

		if wrapList {
			builder.WriteString("}")
		}

		builder.WriteString(",")
	}

	if !list {
		builder.WriteString("}")
	}

	return builder.String()
}

func (q Query) Exec(ctx context.Context, into interface{}) error {
	payload := engine.GQLRequest{
		Query:     q.Build(),
		Variables: map[string]interface{}{},
	}
	return q.Do(ctx, payload, into)
}

func (q Query) Do(ctx context.Context, payload interface{}, into interface{}) error {
	if q.Engine == nil {
		return fmt.Errorf("client.Prisma.Connect() needs to be called before sending queries")
	}

	logger.Debug.Printf("[timing] building %q", time.Since(q.Start))

	err := q.Engine.Do(ctx, payload, into)
	now := time.Now()
	totalDuration := now.Sub(q.Start)
	logger.Debug.Printf("[timing] TOTAL %q", totalDuration)
	return err
}

func Value(value interface{}) []byte {
	v, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	return v
}
