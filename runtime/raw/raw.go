package raw

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/prisma/prisma-client-go/engine"
	"github.com/prisma/prisma-client-go/runtime/builder"
	"github.com/prisma/prisma-client-go/runtime/types/raw"
)

type Raw struct {
	Engine engine.Engine
}

func doRaw(engine engine.Engine, action string, query string, params ...interface{}) builder.Query {
	q := builder.NewQuery()
	q.Engine = engine
	q.Operation = "mutation"
	q.Method = action

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
		switch p := param.(type) {
		case time.Time, raw.Time:
			data, err := json.Marshal(p)
			if err != nil {
				panic(err)
			}
			newParams += fmt.Sprintf(`{"prisma__type":"date","prisma__value":%s}`, string(data))
		default:
			newParams += string(builder.Value(p))
		}
	}
	newParams += "]"

	log.Printf("newParams: %s", newParams)

	q.Inputs = append(q.Inputs, builder.Input{
		Name:  "parameters",
		Value: newParams,
	})

	return q
}
