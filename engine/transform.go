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

// TransformResponse for raw queries
func TransformResponse(data []byte) ([]byte, error) {
	// TODO properly detect a json response
	if !strings.HasPrefix(string(data), `{"columns":[`) {
		return data, nil
	}

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
