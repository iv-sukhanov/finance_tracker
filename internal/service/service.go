package service

import (
	"github.com/iv-sukhanov/finance_tracker/internal/repostitory"

	"github.com/sirupsen/logrus"
)

type SpendingType interface {
}

type SpendingRecord interface {
}

type Service struct {
	SpendingType
	SpendingRecord

	repo *repostitory.Repostitory
	log  *logrus.Entry
}

func New(repo *repostitory.Repostitory, log *logrus.Entry) *Service {
	return &Service{repo: repo, log: log}
}
