package engine

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/prisma/prisma-client-go/binaries"
	"github.com/prisma/prisma-client-go/logger"
	"github.com/prisma/prisma-client-go/runtime/types"
)

func NewDataProxyEngine(schema, schemaEnvVarName string) *DataProxyEngine {
	return &DataProxyEngine{
		Schema:           schema,
		schemaEnvVarName: schemaEnvVarName,
		http:             &http.Client{},
	}
}

type DataProxyEngine struct {
	// http is the internal http client
	http *http.Client

	// url holds the query-engine url
	url string

	// schemaEnvVarName is the env var for the datasource url
	// this is needed internally to extract the api key
	schemaEnvVarName string

	// Schema contains the prisma Schema
	Schema string

	// apiKey contains the parsed prisma data proxy api key from the connection string
	apiKey string
}

func (e *DataProxyEngine) Connect() error {
	// Example uri: https://aws-eu-west-1.prisma-data.com/2.26.0/412bf0a1742a576d699fbd5102a4f725557eff3992995f2e18febce128794961/
	hash := hashSchema(e.Schema)
	logger.Debug.Printf("local schema hash %s", hash)

	logger.Debug.Printf("parsing connection string from database url %s", e.schemaEnvVarName)

	connectionString := os.Getenv(e.schemaEnvVarName)
	if connectionString == "" {
		return fmt.Errorf("no connection string found")
	}
	u, err := url.Parse(connectionString)
	if err != nil {
		return fmt.Errorf("parse prisma string: %w", err)
	}
	e.apiKey = u.Query().Get("api_key")
	if e.apiKey == "" {
		return fmt.Errorf("could not parse api key from data proxy prisma connection string")
	}

	e.url = getCloudURI(u.Host, hash)
	// TODO temp: remove this once prisma data proxy works with recent stable versions – right now only 3.0.1 works
	e.url = "https://" + path.Join(u.Host, "3.0.1", hash)
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

	body, err := e.request(ctx, "POST", "/graphql", data)
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
			return types.ErrNotFound
		}
		return fmt.Errorf("pql error: %s", first.Message)
	}

	if err := json.Unmarshal(response.Data.Result, into); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	logger.Debug.Printf("[timing] request unmarshal took %s", time.Since(startParse))

	return nil
}

func (e *DataProxyEngine) Batch(ctx context.Context, payload interface{}, into interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("payload marshal: %w", err)
	}

	body, err := e.request(ctx, "POST", "/graphql", data)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if err := json.Unmarshal(body, &into); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	return nil
}

func (e *DataProxyEngine) Name() string {
	return "data-proxy"
}

func (e *DataProxyEngine) request(ctx context.Context, method string, path string, payload []byte) ([]byte, error) {
	// TODO use specific log level
	if logger.Enabled {
		logger.Debug.Printf("prisma engine payload: `%s`", payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, e.url+path, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("raw post: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.apiKey))
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

	if rawResponse.StatusCode != http.StatusOK && rawResponse.StatusCode != http.StatusCreated {
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
