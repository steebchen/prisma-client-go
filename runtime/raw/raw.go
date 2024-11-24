package raw

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"

	"github.com/steebchen/prisma-client-go/engine"
	"github.com/steebchen/prisma-client-go/logger"
	"github.com/steebchen/prisma-client-go/runtime/builder"
	"github.com/steebchen/prisma-client-go/runtime/types/raw"
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
		newParams += convertType(param)
	}
	newParams += "]"

	logger.Debug.Printf("raw params: %s", newParams)

	q.Inputs = append(q.Inputs, builder.Input{
		Name:  "parameters",
		Value: newParams,
	})

	return q
}

func doCommandRaw(engine engine.Engine, action string, cmd string) builder.Query {
	q := builder.NewQuery()
	q.Engine = engine
	q.Operation = "mutation"
	q.Method = action

	q.Inputs = append(q.Inputs, builder.Input{
		Name:  "command",
		Value: cmd,
	})

	return q
}

func convertType(input interface{}) string {
	data, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	switch p := input.(type) {
	case time.Time, *time.Time, raw.DateTime, *raw.DateTime:
		return fmt.Sprintf(`{"prisma__type":"date","prisma__value":%s}`, string(data))
	case decimal.Decimal, *decimal.Decimal, raw.Decimal, *raw.Decimal:
		return fmt.Sprintf(`{"prisma__type":"decimal","prisma__value":%q}`, string(data))
	case json.RawMessage, *json.RawMessage, raw.JSON, *raw.JSON:
		encoded := base64.URLEncoding.EncodeToString(data)
		return fmt.Sprintf(`{"prisma__type":"json","prisma__value":%q}`, encoded)
	case []byte, *[]byte, raw.Bytes, *raw.Bytes:
		encoded := base64.URLEncoding.EncodeToString(data)
		return fmt.Sprintf(`{"prisma__type":"bytes","prisma__value":%q}`, encoded)
	default:
		return string(builder.Value(p))
	}
}
