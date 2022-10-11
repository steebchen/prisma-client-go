package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/prisma/prisma-client-go/logger"
	"github.com/prisma/prisma-client-go/runtime/types"
	"strconv"
	"time"
)

var prismaQueryEngineMap = map[int64]*QueryEngine{}

type EngineFactory interface {
	GetPrismaQueryEngine() (*QueryEngine, error)
	ReloadPrismaQueryEngine() error
	QuerySchema(param GQLRequest, result interface{}) error
}

// QueryEngineFactory 创建工厂结构体并实现工厂接口
type QueryEngineFactory struct {
	Key      int64  `json:"key"`
	DBSchema string `json:"DBSchema"`
}

func NewQueryEngineFactory(key int64, dbSchema string) QueryEngineFactory {
	return QueryEngineFactory{
		Key:      key,
		DBSchema: dbSchema,
	}
}

func (q *QueryEngineFactory) GetPrismaQueryEngine() (*QueryEngine, error) {
	// 如果不存在
	if _, ok := prismaQueryEngineMap[q.Key]; !ok {
		// 创建
		content, err := Pull(q.DBSchema)
		if err != nil {
			logger.Debug.Printf("connect fail err : ", err)
			return nil, err
		}
		queryEngine := NewQueryEngine(content, false)
		if err := queryEngine.ConnectSDK(); err != nil {
			logger.Debug.Printf("connect fail err : ", err)
			return nil, err
		}
		prismaQueryEngineMap[q.Key] = queryEngine
	}
	return prismaQueryEngineMap[q.Key], nil
}

func (q *QueryEngineFactory) ReloadPrismaQueryEngine() error {
	// 先销毁旧的引擎
	prismaQueryEngineMap[q.Key].Disconnect()
	// 创建新引擎
	content, err := Pull(q.DBSchema)
	if err != nil {
		logger.Debug.Printf("connect fail err : ", err)
		return err
	}
	queryEngine := NewQueryEngine(content, false)
	if err := queryEngine.ConnectSDK(); err != nil {
		logger.Debug.Printf("connect fail err : ", err)
		return err
	}
	// 存入
	prismaQueryEngineMap[q.Key] = queryEngine
	return nil
}

func DisConnectEngine() {
	for key, engine := range prismaQueryEngineMap {
		engine.Disconnect()
		delete(prismaQueryEngineMap, key)
	}
}

func (q *QueryEngineFactory) QuerySchema(param GQLRequest, result interface{}) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Debug.Println(err)
		}
	}()
	ctx := context.TODO()
	queryEngine, err := q.GetPrismaQueryEngine()
	if err != nil {
		return err
	}
	err = queryEngine.DoQuery(ctx, param, result)
	if err != nil {
		return err
	}
	return nil
}

type GQLResult struct {
	Data       interface{}            `json:"data"`
	Errors     []GQLError             `json:"errors"`
	Extensions map[string]interface{} `json:"extensions"`
}

// DoQuery sends the http Request to the query engine and unmarshals the response
func (e *QueryEngine) DoQuery(ctx context.Context, payload interface{}, v interface{}) error {
	startReq := time.Now()

	body, err := e.Request(ctx, "POST", "/", payload)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	logger.Debug.Printf("[timing] query engine request took %s", time.Since(startReq))

	startParse := time.Now()

	var response GQLResult
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	if len(response.Errors) > 0 {
		first := response.Errors[0]
		if first.RawMessage() == internalUpdateNotFoundMessage ||
			first.RawMessage() == internalDeleteNotFoundMessage {
			return types.ErrNotFound
		}
		return fmt.Errorf("pql error: %s", first.RawMessage())
	}

	if err := json.Unmarshal([]byte(InterfaceToString(response.Data)), v); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	logger.Debug.Printf("[timing] request unmarshaling took %s", time.Since(startParse))

	return nil
}

func InterfaceToString(i interface{}) string {
	switch v := i.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case bool:
		return strconv.FormatBool(v)
	default:
		bytes, _ := json.Marshal(v)
		return string(bytes)
	}
}
