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

type Boolean bool

func (r *Boolean) UnmarshalJSON(b []byte) error {
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
	*r = Boolean(n)
	return nil
}
