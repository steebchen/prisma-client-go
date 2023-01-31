package raw

import (
	"encoding/json"
	"fmt"
)

type prismaBytesValue struct {
	Value []uint8 `json:"prisma__value"`
	Type  string  `json:"prisma__type"`
}

type Bytes []byte

func (r *Bytes) UnmarshalJSON(b []byte) error {
	var v prismaBytesValue
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	if v.Type != "bytes" {
		return fmt.Errorf("invalid type %s, expected bytes", v.Type)
	}
	*r = v.Value
	return nil
}
