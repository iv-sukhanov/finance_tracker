package service

import (
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"

	"github.com/sirupsen/logrus"
)

type User interface {
	AddUsers(users []ftracker.User) ([]uuid.UUID, error)
	GetUsers(opts ...UserOption) ([]ftracker.User, error)
}

type SpendingCategory interface {
	AddCategories(categories []ftracker.SpendingCategory) ([]uuid.UUID, error)
	GetCategories(opts ...CategoryOption) ([]ftracker.SpendingCategory, error)
}

type SpendingRecord interface {
	AddRecords(records []ftracker.SpendingRecord) ([]uuid.UUID, error)
	GetRecords(opts ...RecordOption) ([]ftracker.SpendingRecord, error)
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
