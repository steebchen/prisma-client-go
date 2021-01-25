package engine

import (
	"encoding/json"
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

type Meta struct {
	Target []string `json:"target"`
}

// GQLError is a GraphQL Error
type GQLError struct {
	Error           string           `json:"error"`
	UserFacingError *UserFacingError `json:"user_facing_error"`
	Path            []string         `json:"path"`
}
