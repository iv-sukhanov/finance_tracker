package utils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ErrSQLUniqueViolation = "23505"
)

func GetItitialError(err error) (initialError error) {
	for err != nil {
		initialError = err
		err = errors.Unwrap(err)
	}
	return initialError
}

func GetSQLErrorCode(err error) string {
	if res, ok := err.(*pgconn.PgError); ok {
		return res.SQLState()
	}
	return ""
}
