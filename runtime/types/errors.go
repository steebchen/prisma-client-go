package types

import "errors"

// ErrNotFound gets returned when a database record does not exist
var ErrNotFound = errors.New("ErrNotFound")
