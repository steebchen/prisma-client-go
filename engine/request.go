package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/steebchen/prisma-client-go/engine/protocol"
	"github.com/steebchen/prisma-client-go/logger"
	"github.com/steebchen/prisma-client-go/runtime/types"
)

var internalUpdateNotFoundMessage = "Error occurred during query execution: InterpretationError(\"Error for binding '0'\", Some(QueryGraphBuilderError(RecordNotFound(\"Record to update not found.\"))))"
var internalDeleteNotFoundMessage = "Error occurred during query execution: InterpretationError(\"Error for binding '0'\", Some(QueryGraphBuilderError(RecordNotFound(\"Record to delete does not exist.\"))))"

// Do sends the http Request to the query engine and unmarshals the response
func (e *QueryEngine) Do(ctx context.Context, payload interface{}, v interface{}) error {
	startReq := time.Now()

	body, err := e.Request(ctx, "POST", "/", payload, true)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	logger.Debug.Printf("[timing] query engine request took %s", time.Since(startReq))
	logger.Debug.Printf("[timing] query engine response %s", body)

	startParse := time.Now()

	var response protocol.GQLResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("json gql response unmarshal: %w", err)
	}

	if len(response.Errors) > 0 {
		e := response.Errors[0]
		if e.RawMessage() == internalUpdateNotFoundMessage ||
			e.RawMessage() == internalDeleteNotFoundMessage {
			return types.ErrNotFound
		}

		if e.UserFacingError != nil {
			return fmt.Errorf("user facing error: %w", e.UserFacingError)
		}

		return fmt.Errorf("internal error: %s", e.RawMessage())
	}

	response.Data.Result, err = TransformResponse(response.Data.Result)
	if err != nil {
		return fmt.Errorf("transform response: %w", err)
	}

	if err := json.Unmarshal(response.Data.Result, v); err != nil {
		return fmt.Errorf("json data result unmarshal: %w", err)
	}

	logger.Debug.Printf("[timing] request unmarshaling took %s", time.Since(startParse))

	return nil
}

// Batch sends a batch request to the query engine; used for transactions
func (e *QueryEngine) Batch(ctx context.Context, payload interface{}, v interface{}) error {
	body, err := e.Request(ctx, "POST", "/", payload, true)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	body, err = TransformResponse(body)
	if err != nil {
		return fmt.Errorf("transform response: %w", err)
	}

	if err := json.Unmarshal(body, &v); err != nil {
		return fmt.Errorf("json body unmarshal: %w", err)
	}

	return nil
}

func (e *QueryEngine) Request(ctx context.Context, method string, path string, payload interface{}, requiresConnection bool) ([]byte, error) {
	if !e.connected && requiresConnection {
		logger.Info.Printf("A query was executed before Connect() was called. Make sure to call .Prisma.Connect() before sending any queries.")
		return nil, fmt.Errorf("client is not connected yet")
	}

	e.mu.RLock()
	if e.disconnected {
		e.mu.RUnlock()
		logger.Info.Printf("A query was executed after Disconnect() was called. Make sure to not send any queries after calling .Prisma.Disconnect() the client.")
		return nil, fmt.Errorf("client is already disconnected")
	}
	e.mu.RUnlock()

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("payload marshal: %w", err)
	}

	return request(ctx, e.http, method, e.httpURL+path, requestBody, func(req *http.Request) {
		req.Header.Set("content-type", "application/json")
	})
}
