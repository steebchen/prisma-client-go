package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/steebchen/prisma-client-go/engine"
	"github.com/steebchen/prisma-client-go/engine/protocol"
	"github.com/steebchen/prisma-client-go/logger"
)

type MethodFormat string

const (
	FindRaw      MethodFormat = "findRaw"
	AggregateRaw MethodFormat = "aggregateRaw"
)

var (
	MethodFormatMaping = map[MethodFormat]string{
		FindRaw:      "find%sRaw",      // find{Model}Raw
		AggregateRaw: "aggregate%sRaw", // aggregate{Model}Raw
	}
)

type Input struct {
	Name     string
	Fields   []Field
	Value    interface{}
	WrapList bool
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

	// List saves whether the fields are a list of items
	List bool

	// WrapList saves whether the field should be wrapped in an individual object
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

func (q Query) Build() (string, error) {
	var builder strings.Builder

	builder.WriteString(q.Operation + " " + q.Name)
	builder.WriteString("{")
	builder.WriteString("result: ")

	str, err := q.BuildInner()
	if err != nil {
		return "", err
	}
	builder.WriteString(str)

	builder.WriteString("}")

	return builder.String(), nil
}

func (q Query) BuildInner() (string, error) {
	var builder strings.Builder
	switch MethodFormat(q.Method) {
	case FindRaw:
		builder.WriteString(fmt.Sprintf(MethodFormatMaping[FindRaw], q.Model))
	case AggregateRaw:
		builder.WriteString(fmt.Sprintf(MethodFormatMaping[AggregateRaw], q.Model))
	default:
		builder.WriteString(q.Method + q.Model)
	}

	if len(q.Inputs) > 0 {
		str, err := q.buildInputs(q.Inputs)
		if err != nil {
			return "", err
		}
		builder.WriteString(str)
	}

	builder.WriteString(" ")

	if len(q.Outputs) > 0 {
		str, err := q.buildOutputs(q.Outputs)
		if err != nil {
			return "", err
		}
		builder.WriteString(str)
	}

	return builder.String(), nil
}

func (q Query) buildInputs(inputs []Input) (string, error) {
	var builder strings.Builder

	builder.WriteString("(")

	for index, i := range inputs {
		builder.WriteString(i.Name)

		builder.WriteString(":")

		if i.Value != nil {
			builder.Write(Value(i.Value))
		} else {
			if i.WrapList {
				builder.WriteString("[")
			}
			str, err := q.buildFields(i.WrapList, i.WrapList, i.Fields)
			if err != nil {
				return "", err
			}
			builder.WriteString(str)
			if i.WrapList {
				builder.WriteString("]")
			}
		}

		if index < len(inputs)-1 {
			builder.WriteString(",")
		}
	}

	builder.WriteString(")")

	return builder.String(), nil
}

func (q Query) buildOutputs(outputs []Output) (string, error) {
	var builder strings.Builder

	builder.WriteString("{")

	for _, o := range outputs {
		builder.WriteString(o.Name + " ")

		if len(o.Inputs) > 0 {
			str, err := q.buildInputs(o.Inputs)
			if err != nil {
				return "", err
			}
			builder.WriteString(str)
		}

		if len(o.Outputs) > 0 {
			str, err := q.buildOutputs(o.Outputs)
			if err != nil {
				return "", err
			}
			builder.WriteString(str)
		}
	}

	builder.WriteString("}")

	return builder.String(), nil
}

var ErrDuplicateField = fmt.Errorf("duplicate field (https://github.com/steebchen/prisma-client-go/issues/1095)")

func (q Query) buildFields(list bool, wrapList bool, fields []Field) (string, error) {
	var builder strings.Builder

	if !list {
		builder.WriteString("{")
	}

	var final []Field
	// remember the order in which the unique fields where added to the map
	var uniqueNames []string

	// check for duplicate fields so that multiple queries on the same field will be shared
	// this is necessary for json filters and more
	uniques := make(map[string]*Field)
	for i, f := range fields {
		if _, ok := uniques[f.Name]; ok {
			// check if field is a model operation
			if f.Fields != nil && f.Name != "AND" && f.Name != "OR" && f.Name != "NOT" {
				// field already exists, join sub-fields
				uniques[f.Name].Fields = append(uniques[f.Name].Fields, f.Fields...)
			} else {
				// if it's a list or just contains a value, just add it, which may result in a duplicate
				// this is necessary for some operations, e.g. linking multiple records
				final = append(final, f)
			}
		} else {
			uniques[f.Name] = &fields[i]
			uniqueNames = append(uniqueNames, f.Name)
		}
	}

	// use the list of unique names to add the unique fields in a deterministic order
	for _, name := range uniqueNames {
		final = append(final, *uniques[name])
	}

	for _, f := range final {
		if err := checkFields(f, f.Fields); err != nil {
			return "", err
		}

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
			str, err := q.buildFields(f.List, f.WrapList, f.Fields)
			if err != nil {
				return "", err
			}
			builder.WriteString(str)
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

	return builder.String(), nil
}

func checkFields(parent Field, fields []Field) error {
	uniqueObjectFields := make(map[string]Field)
	for _, f := range fields {
		if f.Value != nil && !f.List && !parent.List {
			if _, ok := uniqueObjectFields[f.Name]; ok {
				return fmt.Errorf("%w: %q", ErrDuplicateField, f.Name)
			}
			uniqueObjectFields[f.Name] = f
		}
	}
	return nil
}

func (q Query) Exec(ctx context.Context, into interface{}) error {
	str, err := q.Build()
	if err != nil {
		return err
	}
	payload := protocol.GQLRequest{
		Query:     str,
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
