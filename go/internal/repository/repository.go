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

// User defines the interface for user repository.
type User interface {
	AddUsers(users []ftracker.User) ([]uuid.UUID, error)
	GetUsers(opts UserOptions) ([]ftracker.User, error)
}

// SpendingCategory defines the interface for spending category repository.
type SpendingCategory interface {
	AddCategories(category []ftracker.SpendingCategory) ([]uuid.UUID, error)
	GetCategories(opts CategoryOptions) ([]ftracker.SpendingCategory, error)
}

// SpendingRecord defines the interface for spending record repository.
type SpendingRecord interface {
	AddRecords(records []ftracker.SpendingRecord) ([]uuid.UUID, error)
	GetRecords(opts RecordOptions) ([]ftracker.SpendingRecord, error)
}

// Repository implements the interfaces for user, spending category, and spending record repositories.
type Repostitory struct {
	User
	SpendingCategory
	SpendingRecord
}

// NewUserRepository creates a new instance of User repository.
func New(db *sqlx.DB) *Repostitory {
	return &Repostitory{
		User:             NewUserRepository(db),
		SpendingCategory: NewCategoryRepository(db),
		SpendingRecord:   NewRecordRepository(db),
	}
}
