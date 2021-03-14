package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

const RFC3339Milli = "2006-01-02T15:04:05.999Z07:00"

type BatchResult struct {
	Count int `json:"count"`
}

// DateTime is a type alias for time.Time
type DateTime = time.Time

type JSON json.RawMessage

// MarshalJSON returns m as the JSON encoding of m.
func (m JSON) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%q", m)), nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *JSON) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("JSON: UnmarshalJSON on nil pointer")
	}
	str, err := strconv.Unquote(string(data))
	if err != nil {
		return errors.New("JSON: UnmarshalJSON unquote error")
	}
	*m = append((*m)[0:0], str...)
	return nil
}

// Direction describes
type Direction string

const (
	ASC  Direction = "asc"
	DESC Direction = "desc"
)
