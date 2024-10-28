package repository

import (
	"fmt"

	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	sqlh "github.com/iv-sukhanov/finance_tracker/internal/utils/sql"
	typesh "github.com/iv-sukhanov/finance_tracker/internal/utils/types"
	"github.com/jmoiron/sqlx"
)

type UserStorage struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) AddUser(user ftracker.User) (uuid.UUID, error) {

	query := fmt.Sprintf("INSERT INTO %s (username, telegram_id) VALUES ($1, $2) RETURNING guid", usersTable)

	var guid uuid.UUID
	if err := s.db.Get(&guid, query, user.Username, user.TelegramID); err != nil {
		return uuid.Nil, fmt.Errorf("Repository.AddUser: %w", err)
	}
	return guid, nil
}

func (s *UserStorage) GetUsers() ([]ftracker.User, error) {

	query := fmt.Sprintf("SELECT guid, username, telegram_id FROM %s", usersTable)
	users := []ftracker.User{}

	if err := s.db.Select(&users, query); err != nil {
		return nil, fmt.Errorf("Repository.GetUsers: %w", err)
	}
	return users, nil
}

func (s *UserStorage) GetUsersByGUIDs(guids []uuid.UUID) ([]ftracker.User, error) {
	query := fmt.Sprintf("SELECT guid, username, telegram_id FROM %s %s", usersTable,
		sqlh.MakeWhereIn("guid", typesh.UUIDsToStrings(guids)...))
}
