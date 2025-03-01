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
	AddUsers(users []ftracker.User) ([]uuid.UUID, error)
	GetUsers(opts UserOptions) ([]ftracker.User, error)
}

type SpendingCategory interface {
	AddCategories(category []ftracker.SpendingCategory) ([]uuid.UUID, error)
	GetCategories(opts CategoryOptions) ([]ftracker.SpendingCategory, error)
}

type SpendingRecord interface {
	AddRecords(records []ftracker.SpendingRecord) ([]uuid.UUID, error)
	GetRecords(opts RecordOptions) ([]ftracker.SpendingRecord, error)
}

type Repostitory struct {
	User
	SpendingCategory
	SpendingRecord
}

func New(db *sqlx.DB) *Repostitory {
	return &Repostitory{
		User:             NewUserRepository(db),
		SpendingCategory: NewCategoryRepository(db),
		SpendingRecord:   NewRecordRepository(db),
	}
}
