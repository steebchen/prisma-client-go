package errors

type UniqueConstraintViolation struct {
	Field string
}

func IsUniqueConstraintViolation() (*UniqueConstraintViolation, bool) {
	return nil, false
}
