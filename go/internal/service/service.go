package service

import (
	"time"

	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"
	"github.com/xuri/excelize/v2"
)

// User defines the interface for user service.
type User interface {
	AddUsers(users []ftracker.User) ([]uuid.UUID, error)
	GetUsers(opts ...UserOption) ([]ftracker.User, error)
	UsersWithLimit(limit int) UserOption
	UsersWithGUIDs(guids []uuid.UUID) UserOption
	UsersWithUsernames(usernames []string) UserOption
	UsersWithTelegramIDs(telegramIDs []string) UserOption
}

// SpendingCategory defines the interface for spending category service.
type SpendingCategory interface {
	AddCategories(categories []ftracker.SpendingCategory) ([]uuid.UUID, error)
	GetCategories(opts ...CategoryOption) ([]ftracker.SpendingCategory, error)
	SpendingCategoriesWithLimit(limit int) CategoryOption
	SpendingCategoriesWithGUIDs(guids []uuid.UUID) CategoryOption
	SpendingCategoriesWithUserGUIDs(guids []uuid.UUID) CategoryOption
	SpendingCategoriesWithCategories(categories []string) CategoryOption
	SpendingCategoriesWithOrder(order CategoryOrder, asc bool) CategoryOption
	CreateExelFromCategories(categories []ftracker.SpendingCategory) (*excelize.File, error)
}

// SpendingRecord defines the interface for spending record service.
type SpendingRecord interface {
	AddRecords(records []ftracker.SpendingRecord) ([]uuid.UUID, error)
	GetRecords(opts ...RecordOption) ([]ftracker.SpendingRecord, error)
	SpendingRecordsWithLimit(limit int) RecordOption
	SpendingRecordsWithGUIDs(guids []uuid.UUID) RecordOption
	SpendingRecordsWithCategoryGUIDs(guids []uuid.UUID) RecordOption
	SpendingRecordsWithTimeFrame(from, to time.Time) RecordOption
	SpendingRecordsWithOrder(order RecordOrder, asc bool) RecordOption
	CreateExelFromRecords(recods []ftracker.SpendingRecord) (*excelize.File, error)
}

// ServiceInterface defines the interface for the service layer.
type ServiceInterface interface {
	User
	SpendingCategory
	SpendingRecord
}

// Service implements the ServiceInterface.
type Service struct {
	User
	SpendingCategory
	SpendingRecord
}

// New creates a new instance of Service with the provided repository.
func New(repo *repository.Repostitory) *Service {
	return &Service{
		User:             NewUserService(repo),
		SpendingCategory: NewCategoryService(repo),
		SpendingRecord:   NewRecordService(repo),
	}
}
