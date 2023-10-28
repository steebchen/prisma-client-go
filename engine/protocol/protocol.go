package protocol

import (
	"encoding/json"
	"strings"
)

// GQLResponse is the default GraphQL response
type GQLResponse struct {
	Data       Data                   `json:"data"`
	Errors     []GQLError             `json:"errors"`
	Extensions map[string]interface{} `json:"extensions"`
}

type Data struct {
	Result json.RawMessage `json:"result"`
}

type GQLBatchResponse struct {
	Errors []GQLError    `json:"errors"`
	Result []GQLResponse `json:"batchResult"`
}

// GQLRequest is the payload for GraphQL queries
type GQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// GQLBatchRequest is the payload for GraphQL queries
type GQLBatchRequest struct {
	Batch       []GQLRequest `json:"batch"`
	Transaction bool         `json:"transaction"`
}

type UserFacingError struct {
	IsPanic   bool   `json:"is_panic"`
	Message   string `json:"message"`
	Meta      Meta   `json:"meta"`
	ErrorCode string `json:"error_code"`
}

func (e *UserFacingError) Error() string {
	return e.Message
}

type Meta struct {
	Target interface{} `json:"target"` // can be of type []string or string
}

// GQLError is a GraphQL Message
type GQLError struct {
	Message         string           `json:"error"`
	UserFacingError *UserFacingError `json:"user_facing_error"`
	Path            []string         `json:"path"`
}

func (e *GQLError) Error() string {
	return e.Message
}

func (e *GQLError) RawMessage() string {
	return strings.ReplaceAll(e.Message, "\n", " ")
}
