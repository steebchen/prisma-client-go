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

type Time struct {
	time.Time
}

func (r *Time) String() string {
	return r.Time.String()
}

func (r *Time) MarshalJSON() ([]byte, error) {
	return r.Time.MarshalJSON()
}

func (r *Time) UnmarshalJSON(data []byte) error {
	var v prismaTimeValue
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v.Type != "date" && v.Type != "datetime" {
		return fmt.Errorf("invalid type %s, expected date", v.Type)
	}
	*r = Time{Time: v.Value}
	return nil
}
