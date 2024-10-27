package service

import (
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"

	"github.com/sirupsen/logrus"
)

type User interface {
	AddUser(user ftracker.User) (uuid.UUID, error)
}

type SpendingType interface {
}

type SpendingRecord interface {
}

type Service struct {
	User
	SpendingType
	SpendingRecord

	log *logrus.Entry
}

func New(repo *repository.Repostitory, log *logrus.Entry) *Service {
	return &Service{
		User: NewUserService(repo),
		log:  log,
	}
}
