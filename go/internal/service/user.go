package service

import (
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"
)

type (
	// UserService implements the User interface.
	UserService struct {
		repo repository.User
	}

	// UserOption is a function to modify the UserOptions.
	UserOption func(*repository.UserOptions)
)

// NewUserService creates a new instance of UserService with the provided repository.
func NewUserService(repo repository.User) *UserService {
	return &UserService{
		repo: repo,
	}
}

// UsersWithLimit is a function that sets the limit for the number of users to be returned.
func (UserService) UsersWithLimit(limit int) UserOption {
	return func(o *repository.UserOptions) {
		o.Limit = limit
	}
}

// UsersWithGUIDs is a function that sets the GUIDs for the users to be returned.
func (UserService) UsersWithGUIDs(guids []uuid.UUID) UserOption {
	return func(o *repository.UserOptions) {
		o.GUIDs = guids
	}
}

// UsersWithUsernames is a function that sets the usernames for the users to be returned.
func (UserService) UsersWithUsernames(usernames []string) UserOption {
	return func(o *repository.UserOptions) {
		o.Usernames = usernames
	}
}

// UsersWithTelegramIDs is a function that sets the telegram IDs for the users to be returned.
func (UserService) UsersWithTelegramIDs(telegramIDs []string) UserOption {
	return func(o *repository.UserOptions) {
		o.TelegramIDs = telegramIDs
	}
}

// GetUsers retrieves a list of users based on the provided options.
//
// Parameters:
//   - options: A variadic list of UserOption functions that modify the
//     UserOptions struct to filter or customize the user retrieval.
//
// Returns:
//   - []ftracker.User: A slice of User objects that match the specified options.
//   - error: An error if the operation fails, or nil if successful.
func (s *UserService) GetUsers(options ...UserOption) ([]ftracker.User, error) {
	var opts repository.UserOptions
	for _, option := range options {
		option(&opts)
	}

	return s.repo.GetUsers(opts)
}

// AddUsers adds a list of users to the repository.
//
// Parameters:
//   - users: A slice of User objects to be added.
//
// Returns:
//   - []uuid.UUID: A slice of UUIDs representing the GUIDs of the added users.
//   - error: An error if the operation fails, or nil if successful.
func (s *UserService) AddUsers(users []ftracker.User) ([]uuid.UUID, error) {
	return s.repo.AddUsers(users)
}
