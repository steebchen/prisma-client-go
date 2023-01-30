package raw

import (
	"encoding/json"
	"fmt"
)

type prismaFloatValue struct {
	Value float64 `json:"prisma__value"`
	Type  string  `json:"prisma__type"`
}

type Float float64

func (r *Float) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%f", *r)), nil
}

func (r *Float) UnmarshalJSON(b []byte) error {
	var v prismaFloatValue
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	if v.Type != "double" {
		return fmt.Errorf("invalid type %s, expected int", v.Type)
	}
	*r = Float(v.Value)
	return nil
}
