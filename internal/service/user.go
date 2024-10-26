package service

import (
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/repostitory"
)

type UserService struct {
	repo repostitory.User
}

func NewUserService(repo repostitory.User) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) AddUser(user ftracker.User) (uuid.UUID, error) {
	return s.repo.AddUser(user)
}
