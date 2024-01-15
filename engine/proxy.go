package engine

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/steebchen/prisma-client-go/binaries"
	"github.com/steebchen/prisma-client-go/engine/protocol"
	"github.com/steebchen/prisma-client-go/logger"
	"github.com/steebchen/prisma-client-go/runtime/types"
)

func NewDataProxyEngine(schema, connectionURL string) *DataProxyEngine {
	return &DataProxyEngine{
		Schema:        schema,
		connectionURL: connectionURL,
		http:          &http.Client{},
	}
}

type DataProxyEngine struct {
	// http is the internal http client
	http *http.Client

	// url holds the query-engine url
	url string

	// connectionURL is the env var for the datasource url
	// this is needed internally to extract the api key
	connectionURL string

	// Schema contains the prisma Schema
	Schema string

	// apiKey contains the parsed prisma data proxy api key from the connection string
	apiKey string
}

func (e *DataProxyEngine) Connect() error {
	// Example uri: https://aws-eu-west-1.prisma-data.com/2.26.0/412bf0a1742a576d699fbd5102a4f725557eff3992995f2e18febce128794961/
	hash := hashSchema(e.Schema)
	logger.Debug.Printf("local schema hash %s", hash)

	logger.Debug.Printf("parsing connection string from database url %s", e.connectionURL)

	u, err := url.Parse(e.connectionURL)
	if err != nil {
		return fmt.Errorf("parse prisma string: %w", err)
	}
	e.apiKey = u.Query().Get("api_key")
	if e.apiKey == "" {
		return fmt.Errorf("could not parse api key from data proxy prisma connection string")
	}

	e.url = getCloudURI(u.Host, hash)
	logger.Debug.Printf("using %s as remote URI", e.url)
	if err := e.uploadSchema(context.Background()); err != nil {
		return fmt.Errorf("upload schema: %w", err)
	}

	return nil
}

func (e *DataProxyEngine) uploadSchema(ctx context.Context) error {
	logger.Debug.Printf("uploading schema...")
	b64Schema := encodeSchema(e.Schema)
	res, err := e.request(ctx, "PUT", "/schema", []byte(b64Schema))
	if err != nil {
		return fmt.Errorf("put schema: %w", err)
	}
	logger.Debug.Printf("schema upload response: %s", res)
	type SchemaResponse struct {
		SchemaHash string `json:"schemaHash"`
	}
	var response SchemaResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return fmt.Errorf("schema response err: %w", err)
	}
	logger.Debug.Printf("remote schema hash %s", response.SchemaHash)
	logger.Debug.Printf("schema upload done.")
	return nil
}

func (e *DataProxyEngine) Disconnect() error {
	return nil
}

func (e *DataProxyEngine) Do(ctx context.Context, payload interface{}, into interface{}) error {
	startReq := time.Now()
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("payload marshal: %w", err)
	}

	body, err := e.retryableRequest(ctx, "POST", "/graphql", data)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	logger.Debug.Printf("[timing] query engine request took %s", time.Since(startReq))

	startParse := time.Now()

	var response protocol.GQLResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("json gql resopnse unmarshal: %w", err)
	}

	if len(response.Errors) > 0 {
		first := response.Errors[0]
		if first.RawMessage() == internalUpdateNotFoundMessage ||
			first.RawMessage() == internalDeleteNotFoundMessage {
			return types.ErrNotFound
		}
		return fmt.Errorf("pql error: %s", first.RawMessage())
	}

	if err := json.Unmarshal(response.Data.Result, into); err != nil {
		return fmt.Errorf("json data result unmarshal: %w", err)
	}

	logger.Debug.Printf("[timing] request unmarshal took %s", time.Since(startParse))

	return nil
}

func (e *DataProxyEngine) Batch(ctx context.Context, payload interface{}, into interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("payload marshal: %w", err)
	}

	body, err := e.retryableRequest(ctx, "POST", "/graphql", data)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if err := json.Unmarshal(body, &into); err != nil {
		return fmt.Errorf("json body unmarshal: %w", err)
	}

	return nil
}

func (e *DataProxyEngine) Name() string {
	return "data-proxy"
}

func (e *DataProxyEngine) request(ctx context.Context, method string, path string, payload []byte) ([]byte, error) {
	logger.Debug.Printf("requesting %s", e.url+path)
	auth := func(req *http.Request) {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.apiKey))
	}
	return request(ctx, e.http, method, e.url+path, payload, auth)
}

func (e *DataProxyEngine) retryableRequest(ctx context.Context, method string, path string, payload []byte) ([]byte, error) {
	res, err := e.request(ctx, method, path, payload)
	if err != nil {
		if !errors.Is(err, errNotFound) {
			return nil, err
		}
		logger.Debug.Printf("got status not found in data proxy request; re-uploading schema")
		if err := e.uploadSchema(ctx); err != nil {
			return nil, fmt.Errorf("upload schema after 400 request: %w", err)
		}
		logger.Debug.Printf("schema re-upload succeeded")
		return e.request(ctx, method, path, payload)
	}
	return res, nil
}

func hashSchema(schema string) string {
	b64Schema := encodeSchema(schema)
	sum := sha256.Sum256([]byte(b64Schema))
	return fmt.Sprintf("%x", sum)
}

func encodeSchema(schema string) string {
	return fmt.Sprint(base64.StdEncoding.EncodeToString([]byte(schema + "\n")))
}

func getCloudURI(host, schemaHash string) string {
	return "https://" + path.Join(host, binaries.PrismaVersion, schemaHash)
}
