package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Do sends the http Request to the query engine and unmarshals the response
func (e *Engine) Do(ctx context.Context, query string, response interface{}) error {
	payload := GQLRequest{
		Query:     query,
		Variables: map[string]interface{}{},
	}

	body, err := e.Request(ctx, "POST", "/", &payload)
	if err != nil {
		return fmt.Errorf("Request failed: %w", err)
	}

	// TODO temporary hack, actually parse the response
	if str := string(body); strings.Contains(str, "errors: \"[{\"error") {
		return fmt.Errorf("pql error: %s", str)
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	return nil
}

func (e *Engine) Request(ctx context.Context, method string, path string, payload interface{}) ([]byte, error) {
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

	rawResponse, err := e.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("raw post: %w", err)
	}
	defer func() {
		err := rawResponse.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	responseBody, err := ioutil.ReadAll(rawResponse.Body)

	if err != nil {
		return nil, fmt.Errorf("raw read: %w", err)
	}

	if rawResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code %d with response %s", rawResponse.StatusCode, responseBody)
	}

	return responseBody, nil
}
