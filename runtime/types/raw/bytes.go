package raw

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type prismaBytesValue struct {
	Value string `json:"prisma__value"`
	Type  string `json:"prisma__type"`
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

	dst := make([]byte, base64.StdEncoding.DecodedLen(len(v.Value)))
	n, err := base64.StdEncoding.Decode(dst, []byte(v.Value))
	if err != nil {
		return err
	}
	dst = dst[:n]

	*r = dst
	return nil
}
