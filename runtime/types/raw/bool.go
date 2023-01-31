package raw

import (
	"encoding/json"
	"fmt"
)

type prismaBoolValue struct {
	// value is raw message as mysql represents bool using ints
	Value interface{} `json:"prisma__value"`
	Type  string      `json:"prisma__type"`
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
	if v.Type != "bool" && v.Type != "int" {
		return fmt.Errorf("invalid type %s, expected bool", v.Type)
	}
	var n bool
	switch d := v.Value.(type) {
	case float64:
		if d == 1 {
			n = true
		} else if d == 0 {
			n = false
		} else {
			return fmt.Errorf("invalid value: %f", d)
		}
	case bool:
		n = d
	default:
		return fmt.Errorf("invalid type: %T", d)
	}
	*r = Bool(n)
	return nil
}
