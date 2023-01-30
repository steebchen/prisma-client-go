package raw

import (
	"encoding/json"
	"fmt"
)

type prismaStringValue struct {
	Value string `json:"prisma__value"`
	Type  string `json:"prisma__type"`
}

type String string

func (r *String) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", *r)), nil
}

func (r *String) UnmarshalJSON(data []byte) error {
	var v prismaStringValue
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v.Type != "string" {
		return fmt.Errorf("invalid type %s, expected string", v.Type)
	}
	*r = String(v.Value)
	return nil
}
