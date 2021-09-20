package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

const RFC3339Milli = "2006-01-02T15:04:05.999Z07:00"

type BatchResult struct {
	Count int `json:"count"`
}

// DateTime is a type alias for time.Time
type DateTime = time.Time

// Decimal points to github.com/shopspring/decimal.Decimal, as Go does not have a native decimal type
type Decimal = decimal.Decimal

// Bytes is a type alias for []byte
type Bytes = []byte

// BigInt is a type alias for int64
type BigInt int64

// UnmarshalJSON converts the Prisma QE value of string to int64
func (m *BigInt) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("BigInt: UnmarshalJSON on nil pointer")
	}
	str, err := strconv.Unquote(string(data))
	if err != nil {
		return fmt.Errorf("BigInt: unquote: %w", err)
	}
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fmt.Errorf("BigInt: UnmarshalJSON error: %w", err)
	}
	*m = BigInt(i)
	return nil
}

// JSON is a new type which implements the correct internal prisma (un)marshaller
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
		return fmt.Errorf("JSON: UnmarshalJSON error: %w", err)
	}
	*m = append((*m)[0:0], str...)
	return nil
}
