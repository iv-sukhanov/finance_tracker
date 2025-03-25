package repository

import (
	"fmt"

	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	"github.com/jmoiron/sqlx"
)

type (
	// UserRepo implements the User interface.
	UserRepo struct {
		db *sqlx.DB
	}

	// UserOptions defines the options for retrieving users.
	UserOptions struct {
		Limit       int
		GUIDs       []uuid.UUID
		Usernames   []string
		TelegramIDs []string
	}
)

// NewUserRepository creates a new instance of UserRepo with the provided database connection.
func NewUserRepository(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

// GetUsers retrieves a list of users from the database based on the provided options.
//
// Parameters:
//   - opts: A struct containing filtering options
//
// Returns:
//   - A slice of User objects matching the specified criteria.
//   - An error if the query fails or any other issue occurs.
func (s *UserRepo) GetUsers(opts UserOptions) ([]ftracker.User, error) {

	whereClause := utils.BindWithOp("AND", true,
		utils.MakeIn("guid", utils.UUIDsToStrings(opts.GUIDs)...),
		utils.MakeIn("username", opts.Usernames...),
		utils.MakeIn("telegram_id", opts.TelegramIDs...),
	)

	query := fmt.Sprintf(
		"SELECT guid, username, telegram_id, created_at, updated_at FROM %s %s %s",
		usersTable,
		whereClause,
		utils.MakeLimit(opts.Limit),
	)

	var users []ftracker.User
	err := s.db.Select(&users, query)
	if err != nil {
		return nil, fmt.Errorf("Repository.GetUsers: %w", err)
	}

	return users, nil
}

// AddUsers inserts multiple users into the database and returns their generated UUIDs.
//
// Parameters:
//   - users: A slice of ftracker.User objects to be inserted into the database.
//
// Returns:
//   - A slice of uuid.UUID representing the generated UUIDs for the inserted users.
//   - An error if any issue occurs during the operation.
func (s *UserRepo) AddUsers(users []ftracker.User) ([]uuid.UUID, error) {

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("Repository.AddUsers: %w", err)
	}

	stmt, err := s.db.PrepareNamed(fmt.Sprintf("INSERT INTO %s (username, telegram_id) VALUES (:username, :telegram_id) RETURNING guid", usersTable))
	if err != nil {
		return nil, fmt.Errorf("Repository.AddUsers: %w", err)
	}

	guids := make([]uuid.UUID, len(users))
	for i, u := range users {
		if err := stmt.Get(&guids[i], u); err != nil {
			_err := tx.Rollback()
			if _err != nil {
				panic(_err)
			}
			return nil, fmt.Errorf("Repository.AddUsers: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	return guids, nil
}
