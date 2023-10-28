package types

import (
	"errors"
	"regexp"
)

// ErrNotFound gets returned when a database record does not exist
var ErrNotFound = errors.New("ErrNotFound")

type F interface {
	~string
}

type ErrUniqueConstraint[T F] struct {
	// Field only shows on Postgres
	Field T
	// Key only shows on MySQL
	Key string
}

const fieldKey = "field"

var prismaMySQLUniqueConstraint = regexp.MustCompile("Unique constraint failed on the constraint: `(?P<" + fieldKey + ">.+)`")
var prismaPostgresUniqueConstraint = regexp.MustCompile("Unique constraint failed on the fields: \\(`(?P<" + fieldKey + ">.+)`\\)")

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
	if match, ok := findMatch(err, prismaMySQLUniqueConstraint); ok {
		return &ErrUniqueConstraint[T]{
			Key: match,
		}, true
	}
	if match, ok := findMatch(err, prismaPostgresUniqueConstraint); ok {
		return &ErrUniqueConstraint[T]{
			Field: T(match),
		}, true
	}
	return nil, false
}

func findMatch(err error, regex *regexp.Regexp) (string, bool) {
	result := regex.FindStringSubmatch(err.Error())
	if result == nil {
		return "", false
	}

	index := regex.SubexpIndex(fieldKey)
	field := result[index]
	return field, true
}
