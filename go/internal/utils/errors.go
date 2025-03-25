package utils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ErrSQLUniqueViolation = "23505"
)

// GetItitialError returns the initial error from a wrapped error chain.
func GetItitialError(err error) (initialError error) {
	for err != nil {
		initialError = err
		err = errors.Unwrap(err)
	}
	return initialError
}

// GetSQLErrorCode returns the SQL error code from a pgconn.PgError.
// If the error is not a pgconn.PgError, it returns an empty string.
func GetSQLErrorCode(err error) string {
	if res, ok := err.(*pgconn.PgError); ok {
		return res.SQLState()
	}
	return ""
}
