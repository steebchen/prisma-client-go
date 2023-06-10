package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/prisma/prisma-client-go/logger"
	"github.com/prisma/prisma-client-go/runtime/types"
)

var internalUpdateNotFoundMessage = "Error occurred during query execution: InterpretationError(\"Error for binding '0'\", Some(QueryGraphBuilderError(RecordNotFound(\"Record to update not found.\"))))"
var internalDeleteNotFoundMessage = "Error occurred during query execution: InterpretationError(\"Error for binding '0'\", Some(QueryGraphBuilderError(RecordNotFound(\"Record to delete does not exist.\"))))"

// Do sends the http Request to the query engine and unmarshals the response
func (e *QueryEngine) Do(ctx context.Context, payload interface{}, v interface{}) error {
	startReq := time.Now()

	body, err := e.Request(ctx, "POST", "/", payload)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	logger.Debug.Printf("[timing] query engine request took %s", time.Since(startReq))
	logger.Debug.Printf("[timing] query engine response %s", body)

	startParse := time.Now()

	var response GQLResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("json gql response unmarshal: %w", err)
	}

	if len(response.Errors) > 0 {
		first := response.Errors[0]
		if first.RawMessage() == internalUpdateNotFoundMessage ||
			first.RawMessage() == internalDeleteNotFoundMessage {
			return types.ErrNotFound
		}
		return fmt.Errorf("pql error: %s", first.RawMessage())
	}

	response.Data.Result, err = transformResponse(response.Data.Result)
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
	body, err := e.Request(ctx, "POST", "/", payload)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	body, err = transformResponse(body)
	if err != nil {
		return fmt.Errorf("transform response: %w", err)
	}

	if err := json.Unmarshal(body, &v); err != nil {
		return fmt.Errorf("json body unmarshal: %w", err)
	}

	return nil
}

func (e *QueryEngine) Request(ctx context.Context, method string, path string, payload interface{}) ([]byte, error) {
	if e.disconnected {
		logger.Info.Printf("A query was executed after Disconnect() was called. Make sure to not send any queries after disconnecting the client.")
		return nil, fmt.Errorf("client is disconnected")
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("payload marshal: %w", err)
	}

	// TODO use specific log level
	if logger.Enabled {
		logger.Debug.Printf("prisma engine payload: `%s`", requestBody)
	}

	return request(ctx, e.http, method, e.url+path, requestBody, func(req *http.Request) {
		req.Header.Set("content-type", "application/json")
	})
}
