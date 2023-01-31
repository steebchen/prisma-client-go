package raw

import (
	"encoding/json"
	"fmt"
)

type prismaIntValue struct {
	Value int    `json:"prisma__value"`
	Type  string `json:"prisma__type"`
}

type Int int

func (r *Int) UnmarshalJSON(b []byte) error {
	var v prismaIntValue
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	if v.Type != "int" {
		return fmt.Errorf("invalid type %s, expected int", v.Type)
	}
	*r = Int(v.Value)
	return nil
}
