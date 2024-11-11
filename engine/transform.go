package engine

import (
	"encoding/json"
	"strings"
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

type MongoCursor struct {
	Cursor struct {
		FirstBatch []map[string]interface{} `json:"firstBatch"`
	} `json:"cursor"`
}

func TransformMongoResponse(data []byte) ([]byte, error) {
	var mongoResponse MongoCursor
	err := json.Unmarshal(data, &mongoResponse)
	if err != nil {
		return nil, err
	}

	// Extract `firstBatch` and transform it into a more readable format
	output := make([]map[string]interface{}, 0)
	for _, doc := range mongoResponse.Cursor.FirstBatch {
		// Create a new map to flatten `$oid` and `$date` fields
		flattened := make(map[string]interface{})
		for k, v := range doc {
			switch val := v.(type) {
			case map[string]interface{}:
				// Handle special cases for `$oid` and `$date` fields
				if oid, exists := val["$oid"]; exists {
					flattened[k] = oid
				} else if date, exists := val["$date"]; exists {
					flattened[k] = date
				} else {
					flattened[k] = val
				}
			default:
				flattened[k] = val
			}
		}
		output = append(output, flattened)
	}

	// Marshal the transformed output back to JSON=
	o, err := json.Marshal(output)
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
	case strings.Contains(string(data), `{"cursor":`):
		return TransformMongoResponse(data)
	}

	return data, nil
}
