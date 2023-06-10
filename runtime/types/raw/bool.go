package raw

import (
	"encoding/json"
	"fmt"
)

type Boolean bool

func (r *Boolean) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var n bool
	switch d := v.(type) {
	case float64:
		// MySQL uses tinyint for booleans
		switch d {
		case 1:
			n = true
		case 0:
			n = false
		default:
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
