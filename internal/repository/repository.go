package repository

import (
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable           = "users"
	spendingTypesTable   = "spending_types"
	spendingRecordsTable = "spending_records"
)

type User interface {
	AddUser(user ftracker.User) (uuid.UUID, error)
}

type SpendingType interface {
}

type SpendingRecord interface {
}

type Repostitory struct {
	User
	SpendingType
	SpendingRecord
}

func NewRepostitory(db *sqlx.DB) *Repostitory {
	return &Repostitory{
		User: NewUserRepository(db),
	}
}
