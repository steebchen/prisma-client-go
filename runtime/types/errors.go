package types

import (
	"errors"

	"github.com/steebchen/prisma-client-go/engine/protocol"
)

// ErrNotFound gets returned when a database record does not exist
var ErrNotFound = errors.New("ErrNotFound")

// IsErrNotFound is true if the error is a ErrNotFound, which gets returned when a database record does not exist
// This can happen when you call `FindUnique` on a record, or update or delete a single record which doesn't exist.
func IsErrNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

type F interface {
	~string
}

type ErrUniqueConstraint[T F] struct {
	// Message is the error message
	Message string
	// Fields only shows on Postgres
	Fields []T
	// Key only shows on MySQL
	Key string
}

// CheckUniqueConstraint returns on a unique constraint error or violation with error info
// Ideally this will be replaced with Prisma-generated errors in the future
func CheckUniqueConstraint[T F](err error) (*ErrUniqueConstraint[T], bool) {
	if err == nil {
		return nil, false
	}

	var ufr *protocol.UserFacingError
	if ok := errors.As(err, &ufr); !ok {
		return nil, false
	}

	if ufr.ErrorCode != "P2002" {
		return nil, false
	}

	// postgres
	if items, ok := ufr.Meta.Target.([]interface{}); ok {
		var fields []T
		for _, f := range items {
			field, ok := f.(string)
			if ok {
				fields = append(fields, T(field))
			}
		}
		return &ErrUniqueConstraint[T]{
			Fields: fields,
		}, true
	}

	// mysql
	if item, ok := ufr.Meta.Target.(string); ok {
		return &ErrUniqueConstraint[T]{
			Key: item,
		}, true
	}

	return nil, false
}
