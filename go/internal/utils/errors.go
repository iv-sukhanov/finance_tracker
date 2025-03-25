package utils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ErrSQLUniqueViolation = "23505"
)

// GetItitialError traverses the chain of wrapped errors and returns the
// initial (root) error in the chain. If the provided error is nil, it
// returns nil.
//
// Parameters:
//   - err: The error to traverse.
//
// Returns:
//   - initialError: The root error in the chain, or nil if the input error is nil.
func GetItitialError(err error) (initialError error) {
	for err != nil {
		initialError = err
		err = errors.Unwrap(err)
	}
	return nil
}

// GetSQLErrorCode extracts the SQL state code from a PostgreSQL error.
//
// Parameters:
//   - err: The error to extract the SQL state code from.
//
// Returns:
//   - A string representing the SQL state code if the error is a PostgreSQL
//     error, or an empty string otherwise.
func GetSQLErrorCode(err error) string {
	if res, ok := err.(*pgconn.PgError); ok {
		return res.SQLState()
	}
	return ""
}
