package raw

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type prismaBigIntValue struct {
	Value string `json:"prisma__value"`
	Type  string `json:"prisma__type"`
}

type BigInt int64

func (r *BigInt) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", *r)), nil
}

func (r *BigInt) UnmarshalJSON(b []byte) error {
	var v prismaBigIntValue
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	if v.Type != "bigint" {
		return fmt.Errorf("invalid type %s, expected bigint", v.Type)
	}
	i, err := strconv.Atoi(v.Value)
	if err != nil {
		return err
	}
	*r = BigInt(i)
	return nil
}
