package repository

import (
	"fmt"

	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (s *UserRepo) AddUser(user ftracker.User) (uuid.UUID, error) {

	stmt, err := s.db.PrepareNamed(fmt.Sprintf("INSERT INTO %s (username, telegram_id) VALUES (:username, :telegram_id) RETURNING guid", usersTable))
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Repository.AddUser: %w", err)
	}

	var guid uuid.UUID
	if err := stmt.Get(&guid, user); err != nil {
		return uuid.Nil, fmt.Errorf("Repository.AddUser: %w", err)
	}
	return guid, nil
}

func (s *UserRepo) GetUsers() ([]ftracker.User, error) {

	query := fmt.Sprintf("SELECT guid, username, telegram_id FROM %s", usersTable)
	users := []ftracker.User{}

	if err := s.db.Select(&users, query); err != nil {
		return nil, fmt.Errorf("Repository.GetUsers: %w", err)
	}
	return users, nil
}

func (s *UserRepo) GetUsersByGUIDs(guids []uuid.UUID) ([]ftracker.User, error) {
	query := fmt.Sprintf("SELECT guid, username, telegram_id FROM %s %s", usersTable,
		utils.MakeWhereIn("guid", utils.UUIDsToStrings(guids)...))
	users := []ftracker.User{}

	if err := s.db.Select(&users, query); err != nil {
		return nil, fmt.Errorf("Repository.GetUsersByGUIDs: %w", err)
	}
	return users, nil
}
