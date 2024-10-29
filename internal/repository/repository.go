package repository

import (
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable              = "users"
	spendingCategoriesTable = "spending_categories"
	spendingRecordsTable    = "spending_records"
)

type User interface {
	AddUser(user ftracker.User) (uuid.UUID, error)
	GetUsers() ([]ftracker.User, error)
	GetUsersByGUIDs(guids []uuid.UUID) ([]ftracker.User, error)
}

type SpendingCategory interface {
	AddCategories(category []ftracker.SpendingCategory) ([]uuid.UUID, error)
}

type SpendingRecord interface {
	AddRecords(records []ftracker.SpendingRecord) ([]uuid.UUID, error)
}

type Repostitory struct {
	User
	SpendingCategory
	SpendingRecord
}

func NewRepostitory(db *sqlx.DB) *Repostitory {
	return &Repostitory{
		User:             NewUserRepository(db),
		SpendingCategory: NewCategoryRepository(db),
		SpendingRecord:   NewRecordRepository(db),
	}
}
