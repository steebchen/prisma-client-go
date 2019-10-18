package jsonrpc

type Request struct {
	ID     int         `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result"`
}

// NewResponse forms a new JSON RPC response to reply to the Prisma CLI commands
func NewResponse(ID int, result interface{}) Response {
	return Response{
		JSONRPC: "2.0",
		ID:      ID,
		Result:  result,
	}
}
