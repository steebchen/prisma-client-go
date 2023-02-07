package raw

import (
	"encoding/json"
	"fmt"
)

type prismaJSONValue struct {
	Value json.RawMessage `json:"prisma__value"`
	Type  string          `json:"prisma__type"`
}

type JSON json.RawMessage

func (r *JSON) UnmarshalJSON(b []byte) error {
	var v prismaJSONValue
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	if v.Type != "json" {
		return fmt.Errorf("invalid type %s, expected json", v.Type)
	}
	*r = JSON(v.Value)
	return nil
}
