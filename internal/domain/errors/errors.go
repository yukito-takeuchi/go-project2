package errors

import "errors"

var (
	ErrItemNotFound   = errors.New("item not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrDatabaseError  = errors.New("database error")
	ErrDuplicateEntry = errors.New("duplicate entry")
)

func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrItemNotFound)
}

func IsDatabaseError(err error) bool {
	return errors.Is(err, ErrDatabaseError)
}

func IsValidationError(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}
