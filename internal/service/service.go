package service

import (
	"time"

	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"

	"github.com/sirupsen/logrus"
)

type User interface {
	AddUsers(users []ftracker.User) ([]uuid.UUID, error)
	GetUsers(opts ...UserOption) ([]ftracker.User, error)
	WithLimit(limit int) UserOption
	WithGUIDs(guids []uuid.UUID) UserOption
	WithUsernames(usernames []string) UserOption
	WithTelegramIDs(telegramIDs []string) UserOption
}

type SpendingCategory interface {
	AddCategories(categories []ftracker.SpendingCategory) ([]uuid.UUID, error)
	GetCategories(opts ...CategoryOption) ([]ftracker.SpendingCategory, error)
	WithLimit(limit int) CategoryOption
	WithGUIDs(guids []uuid.UUID) CategoryOption
	WithUserGUIDs(guids []uuid.UUID) CategoryOption
	WithCategories(categories []string) CategoryOption
	WithOrder(order CategoryOrder, asc bool) CategoryOption
}

type SpendingRecord interface {
	AddRecords(records []ftracker.SpendingRecord) ([]uuid.UUID, error)
	GetRecords(opts ...RecordOption) ([]ftracker.SpendingRecord, error)
	WithLimit(limit int) RecordOption
	WithGUIDs(guids []uuid.UUID) RecordOption
	WithCategoryGUIDs(guids []uuid.UUID) RecordOption
	WithTimeFrame(from, to time.Time) RecordOption
	WithOrder(order RecordOrder, asc bool) RecordOption
}

type Service struct {
	User
	SpendingCategory
	SpendingRecord

	log *logrus.Entry
}

func New(repo *repository.Repostitory, log *logrus.Entry) *Service {
	return &Service{
		User:             NewUserService(repo),
		SpendingCategory: NewCategoryService(repo),
		SpendingRecord:   NewRecordService(repo),
		log:              log,
	}
}
