package runtime

import (
	"time"
)

const RFC3339Milli = "2006-01-02T15:04:05.999Z07:00"

// DateTime is a type alias for time.Time
type DateTime = time.Time

// Direction describes
type Direction string

const (
	ASC  Direction = "ASC"
	DESC Direction = "DESC"
)
