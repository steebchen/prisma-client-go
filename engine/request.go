package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/prisma/prisma-client-go/generator/runtime"
	"github.com/prisma/prisma-client-go/logger"
)

var internalUpdateNotFoundMessage = "Error occurred during query execution:\nInterpretationError(\"Error for binding" +
	" \\'0\\'\", Some(QueryGraphBuilderError(RecordNotFound(\"Record to update not found.\"))))"
var internalDeleteNotFoundMessage = "Error occurred during query execution:\nInterpretationError(\"Error for binding" +
	" \\'0\\'\", Some(QueryGraphBuilderError(RecordNotFound(\"Record to delete does not exist.\"))))"

// Do sends the http Request to the query engine and unmarshals the response
func (e *QueryEngine) Do(ctx context.Context, payload interface{}, v interface{}) error {
	startReq := time.Now()

	body, err := e.Request(ctx, "POST", "/", payload)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	logger.Debug.Printf("[timing] query engine request took %s", time.Since(startReq))

	startParse := time.Now()

	var response GQLResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	if len(response.Errors) > 0 {
		first := response.Errors[0]
		if first.Message == internalUpdateNotFoundMessage ||
			first.Message == internalDeleteNotFoundMessage {
			return runtime.ErrNotFound
		}
		return fmt.Errorf("pql error: %s", first.Message)
	}

	if err := json.Unmarshal(response.Data.Result, v); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	logger.Debug.Printf("[timing] request unmarshaling took %s", time.Since(startParse))

	return nil
}

func (e *QueryEngine) Request(ctx context.Context, method string, path string, payload interface{}) ([]byte, error) {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("payload marshal: %w", err)
	}

	req, err := http.NewRequest(method, e.url+path, bytes.NewBuffer(requestBody))

	if err != nil {
		return nil, fmt.Errorf("raw post: %w", err)
	}

	req.Header.Set("content-type", "application/json")
	req = req.WithContext(ctx)

	startReq := time.Now()
	rawResponse, err := e.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("raw post: %w", err)
	}
	defer func() {
		if err := rawResponse.Body.Close(); err != nil {
			panic(err)
		}
	}()
	reqDuration := time.Since(startReq)
	logger.Debug.Printf("[timing] query engine raw request took %s", reqDuration)

	responseBody, err := ioutil.ReadAll(rawResponse.Body)

	if err != nil {
		return nil, fmt.Errorf("raw read: %w", err)
	}

	if rawResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code %d with response %s", rawResponse.StatusCode, responseBody)
	}

	if logger.Enabled {
		if elapsedRaw := rawResponse.Header["X-Elapsed"]; len(elapsedRaw) > 0 {
			elapsed, _ := strconv.Atoi(elapsedRaw[0])
			duration := time.Duration(elapsed) * time.Microsecond
			logger.Debug.Printf("[timing] elapsed: %s", duration)

			diff := reqDuration - duration
			logger.Debug.Printf("[timing] just http: %s", diff)
			logger.Debug.Printf("[timing] http percentage: %.2f%%", float64(diff)/float64(reqDuration)*100)
		}
	}

	return responseBody, nil
}
