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
	AddUsers(user ftracker.User) (uuid.UUID, error)
	//GetUsers(opts ...UserOption) ([]ftracker.User, error)
}

type SpendingCategory interface {
	AddCategories(category []ftracker.SpendingCategory) ([]uuid.UUID, error)
	//GetCategories(opts ...CategoryOption) ([]ftracker.SpendingCategory, error)
}

type SpendingRecord interface {
	AddRecords(records []ftracker.SpendingRecord) ([]uuid.UUID, error)
	GetRecords(opts ...RecordOption) ([]ftracker.SpendingRecord, error)
}

type Repostitory struct {
	User
	SpendingCategory
	SpendingRecord
}

func NewRepostitory(db *sqlx.DB) *Repostitory {
	return &Repostitory{
		//User:             NewUserRepository(db),
		SpendingCategory: NewCategoryRepository(db),
		SpendingRecord:   NewRecordRepository(db),
	}
}
