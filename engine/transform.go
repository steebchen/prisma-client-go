package engine

import (
	"encoding/json"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Input struct {
	Columns []string        `json:"columns"`
	Types   []string        `json:"types"`
	Rows    [][]interface{} `json:"rows"`
}

func TransformSQLResponse(data []byte) ([]byte, error) {
	var input Input
	err := json.Unmarshal(data, &input)
	if err != nil {
		return nil, err
	}

	output := make([]map[string]interface{}, 0)

	for _, row := range input.Rows {
		m := make(map[string]interface{})
		for i, column := range input.Columns {
			m[column] = row[i]
		}
		output = append(output, m)
	}

	o, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func TransformMongoResponse(data []byte) ([]byte, error) {
	var result []map[string]interface{}

	if err := bson.UnmarshalExtJSON(data, false, &result); err != nil {
		return nil, err
	}

	for _, doc := range result {
		if doc["id"] == nil {
			doc["id"] = doc["_id"]
		}
	}

	o, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return o, nil
}

// TransformResponse for raw queries
func TransformResponse(data []byte) ([]byte, error) {
	// TODO properly detect a json response
	switch {
	case strings.HasPrefix(string(data), `{"columns":[`):
		return TransformSQLResponse(data)

	// https://github.com/mongodb/mongo-go-driver/blob/91abd887f6b44ab56f47e58430f57b1be1996ceb/bson/extjson_wrappers.go#L18
	case strings.Contains(string(data), `{"$oid":`),
		strings.Contains(string(data), `{"$date":`),
		strings.Contains(string(data), `{"$numberInt":`),
		strings.Contains(string(data), `{"$numberLong":`),
		strings.Contains(string(data), `{"$symbol":`),
		strings.Contains(string(data), `{"$numberDouble":`),
		strings.Contains(string(data), `{"$numberDecimal":`),
		strings.Contains(string(data), `{"$binary":`),
		strings.Contains(string(data), `{"$code":`),
		strings.Contains(string(data), `{"$scope":`),
		strings.Contains(string(data), `{"$timestamp":`),
		strings.Contains(string(data), `{"$regularExpression":`),
		strings.Contains(string(data), `{"$dbPointer":`),
		strings.Contains(string(data), `{"$minKey":`),
		strings.Contains(string(data), `{"$maxKey":`),
		strings.Contains(string(data), `{"$undefined":`):
		return TransformMongoResponse(data)
	}

	return data, nil
}
