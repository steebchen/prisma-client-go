package raw

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"

	"github.com/prisma/prisma-client-go/engine"
	"github.com/prisma/prisma-client-go/logger"
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
		data, err := json.Marshal(param)
		if err != nil {
			panic(err)
		}
		switch p := param.(type) {
		case time.Time, *time.Time, raw.DateTime, *raw.DateTime:
			newParams += fmt.Sprintf(`{"prisma__type":"date","prisma__value":%s}`, string(data))
		case decimal.Decimal, *decimal.Decimal, raw.Decimal, *raw.Decimal:
			newParams += fmt.Sprintf(`{"prisma__type":"decimal","prisma__value":%q}`, string(data))
		case json.RawMessage, *json.RawMessage, raw.JSON, *raw.JSON:
			encoded := base64.URLEncoding.EncodeToString(data)
			newParams += fmt.Sprintf(`{"prisma__type":"json","prisma__value":%q}`, encoded)
		case []byte, *[]byte, raw.Bytes, *raw.Bytes:
			encoded := base64.URLEncoding.EncodeToString(data)
			newParams += fmt.Sprintf(`{"prisma__type":"bytes","prisma__value":%q}`, encoded)
		default:
			newParams += string(builder.Value(p))
		}
	}
	newParams += "]"

	logger.Debug.Printf("raw params: %s", newParams)

	q.Inputs = append(q.Inputs, builder.Input{
		Name:  "parameters",
		Value: newParams,
	})

	return q
}
