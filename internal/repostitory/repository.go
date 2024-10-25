package repostitory

import "github.com/jmoiron/sqlx"

type SpendingType interface {
}

type SpendingRecord interface {
}

type Repostitory struct {
	SpendingType
	SpendingRecord

	db *sqlx.DB
}

func NewRepostitory(db *sqlx.DB) *Repostitory {
	return &Repostitory{db: db}
}
