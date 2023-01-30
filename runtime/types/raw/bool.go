package raw

import (
	"encoding/json"
	"fmt"
)

type prismaBoolValue struct {
	Value bool   `json:"prisma__value"`
	Type  string `json:"prisma__type"`
}

type Bool bool

func (r *Bool) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%t", *r)), nil
}

func (r *Bool) UnmarshalJSON(b []byte) error {
	var v prismaBoolValue
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	if v.Type != "bool" {
		return fmt.Errorf("invalid type %s, expected int", v.Type)
	}
	*r = Bool(v.Value)
	return nil
}
