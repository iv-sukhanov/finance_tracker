package service

import (
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"
)

type (
	UserService struct {
		repo repository.User
	}

	UserOption func(*repository.UserOptions)
)

func NewUserService(repo repository.User) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (UserService) UsersWithLimit(limit int) UserOption {
	return func(o *repository.UserOptions) {
		o.Limit = limit
	}
}

func (UserService) UsersWithGUIDs(guids []uuid.UUID) UserOption {
	return func(o *repository.UserOptions) {
		o.GUIDs = guids
	}
}

func (UserService) UsersWithUsernames(usernames []string) UserOption {
	return func(o *repository.UserOptions) {
		o.Usernames = usernames
	}
}

func (UserService) UsersWithTelegramIDs(telegramIDs []string) UserOption {
	return func(o *repository.UserOptions) {
		o.TelegramIDs = telegramIDs
	}
}

func (s *UserService) GetUsers(options ...UserOption) ([]ftracker.User, error) {
	var opts repository.UserOptions
	for _, option := range options {
		option(&opts)
	}

	return s.repo.GetUsers(opts)
}

func (s *UserService) AddUsers(users []ftracker.User) ([]uuid.UUID, error) {
	return s.repo.AddUsers(users)
}
