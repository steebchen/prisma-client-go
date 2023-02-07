package raw

import (
	"encoding/json"
	"fmt"
	"time"
)

type prismaTimeValue struct {
	Value time.Time `json:"prisma__value"`
	Type  string    `json:"prisma__type"`
}

type DateTime struct {
	time.Time
}

func (r *DateTime) String() string {
	return r.Time.String()
}

func (r *DateTime) MarshalJSON() ([]byte, error) {
	return r.Time.MarshalJSON()
}

func (r *DateTime) UnmarshalJSON(data []byte) error {
	var v prismaTimeValue
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v.Type != "date" && v.Type != "datetime" {
		return fmt.Errorf("invalid type %s, expected date or datetime", v.Type)
	}
	*r = DateTime{Time: v.Value}
	return nil
}
