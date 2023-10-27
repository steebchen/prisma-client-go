package types

import (
	"errors"
	"strings"
)

// ErrNotFound gets returned when a database record does not exist
var ErrNotFound = errors.New("ErrNotFound")

type A string
type B string

type F interface {
	~string
}

type ErrUniqueConstraint[T F] struct {
	Field T
}

const prismaUniqueConstraint = "Unique constraint failed on the fields: (`%s`)"

// CheckUniqueConstraint returns on a unique constraint error or violation with error info
// Use as follows:
//
//	user, err := db.User.CreateOne(...).Exec(cxt)
//	if err != nil {
//		if info, err := db.UniqueConstraintError(); err != nil {
//			log.Printf("unique constraint on the field: %s", info.Field)
//		}
//	}
//
// Ideally this will be replaced with Prisma-generated errors in the future
func CheckUniqueConstraint[T F](err error) (*ErrUniqueConstraint[T], bool) {
	// TODO use regex
	if !strings.Contains(err.Error(), prismaUniqueConstraint) {
		return nil, false
	}
	return &ErrUniqueConstraint[T]{
		Field: "asdf",
	}, true
}

// ----------
// THIS IS GENERATED CODE
// ----------

type Fields string

// TODO check what JS client uses for fields exports

const UserModelNameField Fields = "user.name"

type RealErrUniqueConstraint = ErrUniqueConstraint[Fields]

func CheckUniqueConstraintError(err error) (*ErrUniqueConstraint[Fields], bool) {
	return CheckUniqueConstraint[Fields](err)
}
