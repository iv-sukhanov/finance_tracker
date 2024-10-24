package spendings

import (
	"github.com/jackc/pgx/v5"
)

func newRepository(db *pgx.Conn) *repository {
	return &repository{db: db}
}
