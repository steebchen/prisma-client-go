package builder

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/prisma/photongo/generator/runtime"

	"github.com/prisma/photongo/logger"
)

type Input struct {
	Name   string
	Fields []Field
	Value  interface{}
}

// output can be a single Name or can have nested fields
type Output struct {
	Name string

	// inputs (optional) to provide arguments to a field
	Inputs []Input

	Outputs []Output
}

type Field struct {
	// The Name of the field.
	Name string

	// an Action for input fields, e.g. `contains`
	Action string

	// whether the fields is a list of items
	List bool

	// whether the a list field should be wrapped in an object
	WrapList bool

	// Value contains the field value. if nil, fields will contain a subselection.
	Value interface{}

	// Fields contains a subselection of fields. If not nil, value will be undefined.
	Fields []Field
}

type Client interface {
	Do(ctx context.Context, query string, into interface{}) error
}

type Query struct {
	// The generic Photon Client
	Client Client

	// operation describes the PQL operation: query, mutation or subscription
	Operation string

	// Name describes the operation; useful for tracing
	Name string

	// method describes a crud operation
	Method string

	// model contains the Prisma model Name
	Model string

	// inputs contains function arguments
	Inputs []Input

	// outputs contains the return fields
	Outputs []Output
}

func (q Query) buildQuery() string {
	var builder strings.Builder

	builder.WriteString(q.Operation + " " + q.Name)
	builder.WriteString("{")

	builder.WriteString(q.Build())

	builder.WriteString("}")

	return builder.String()
}

func (q Query) Build() string {
	var builder strings.Builder

	builder.WriteString(q.Method + q.Model)

	if len(q.Inputs) > 0 {
		builder.WriteString(buildInputs(q.Inputs))
	}

	builder.WriteString(" ")

	builder.WriteString(buildOutputs(q.Outputs))

	return builder.String()
}

func buildInputs(inputs []Input) string {
	var builder strings.Builder

	builder.WriteString("(")

	for _, i := range inputs {
		builder.WriteString(i.Name)

		builder.WriteString(":")

		if i.Value != nil {
			builder.WriteString(Value(i.Name, i.Value))
		} else {
			builder.WriteString(buildFields(false, false, i.Fields))
		}

		builder.WriteString(",")
	}

	builder.WriteString(")")

	return builder.String()
}

func buildOutputs(outputs []Output) string {
	var builder strings.Builder

	builder.WriteString("{")

	for _, o := range outputs {
		builder.WriteString(o.Name + " ")

		if len(o.Inputs) > 0 {
			log.Printf("building inputs: %d %+v", len(o.Inputs), o.Inputs)
			builder.WriteString(buildInputs(o.Inputs))
		}

		if len(o.Outputs) > 0 {
			builder.WriteString(buildOutputs(o.Outputs))
		}
	}

	builder.WriteString("}")

	return builder.String()
}

func buildFields(list bool, wrapList bool, fields []Field) string {
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

		if f.Name != "" && f.Action != "" {
			builder.WriteString("_" + f.Action)
		}

		if f.Name != "" {
			builder.WriteString(":")
		}

		if f.List {
			builder.WriteString("[")
		}

		if f.Fields != nil {
			builder.WriteString(buildFields(f.List, f.WrapList, f.Fields))
		}

		builder.WriteString(Value(f.Name, f.Value))

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

func (q Query) Exec(ctx context.Context, v interface{}) error {
	if q.Client == nil {
		panic("client.Connect() needs to be called before sending queries")
	}

	s := q.buildQuery()

	// TODO use specific log level
	if logger.Enabled {
		logger.Debug.Printf("prisma query: `%s`", s)
	}

	return q.Client.Do(ctx, s, &v)
}

func Value(name string, value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf(`%q`, v)
	case *string:
		if v == nil {
			return "null"
		}
		return fmt.Sprintf(`%q`, *v)
	case bool:
		return fmt.Sprintf(`%t`, v)
	case *bool:
		if v == nil {
			return "null"
		}
		return fmt.Sprintf(`%t`, *v)
	case int:
		return fmt.Sprintf(`%d`, v)
	case *int:
		if v == nil {
			return "null"
		}
		return fmt.Sprintf(`%d`, *v)
	case float64:
		return fmt.Sprintf(`%f`, v)
	case *float64:
		if v == nil {
			return "null"
		}
		return fmt.Sprintf(`%f`, *v)
	case runtime.DateTime:
		return fmt.Sprintf(`"%s"`, v.UTC().Format(runtime.RFC3339Milli))
	case *runtime.DateTime:
		if v == nil {
			return "null"
		}
		return fmt.Sprintf(`"%s"`, v.UTC().Format(runtime.RFC3339Milli))
	case runtime.Direction:
		return string(v)
	case nil:
		return ""

	// TODO handle enums
	// {{ range $t := $.DMMF.Datamodel.Enums }}
	// case {{ $t.Name}}:
	// 	return fmt.Sprintf(`%s`, v)
	// case *{{ $t.Name}}:
	// 	if v == nil {
	// 		return "null"
	// 	}
	// 	return fmt.Sprintf(`%s`, *v)
	// {{ end }}

	default:
		panic(fmt.Errorf("no branch for field %s of type %T", name, v))
	}
}
