package engine

// GQLResponse is the payload for a GraphQL response
type GQLResponse struct {
	Data       interface{}            `json:"data"`
	Errors     []GQLError             `json:"errors"`
	Extensions map[string]interface{} `json:"extensions"`
}

// GQLRequest is the payload for GraphQL queries
type GQLRequest struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	OperationName *string                `json:"operationName"`
}

// GQLError is a GraphQL Error
type GQLError struct {
	Message    string                 `json:"error"` // note: the query-engine uses 'error' instead of 'message'
	Path       []string               `json:"path"`
	Extensions map[string]interface{} `json:"query"`
}
