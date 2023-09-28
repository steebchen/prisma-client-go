package engine

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/steebchen/prisma-client-go/logger"
)

var errNotFound = fmt.Errorf("not found; re-upload schema")

func request(ctx context.Context, client *http.Client, method string, url string, payload []byte, apply func(*http.Request)) ([]byte, error) {
	if logger.Enabled {
		logger.Debug.Printf("prisma engine payload: `%s`", payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("raw post: %w", err)
	}

	apply(req)

	req = req.WithContext(ctx)

	startReq := time.Now()
	rawResponse, err := client.Do(req)
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

	responseBody, err := io.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("raw read: %w", err)
	}

	if rawResponse.StatusCode == http.StatusNotFound {
		logger.Debug.Printf("status not found with response body %s", responseBody)
		return nil, errNotFound
	}

	if rawResponse.StatusCode != http.StatusOK && rawResponse.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("http status code %d with response %s", rawResponse.StatusCode, responseBody)
	}

	if logger.Enabled {
		logger.Debug.Printf("prisma engine response: `%s`", responseBody)

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
