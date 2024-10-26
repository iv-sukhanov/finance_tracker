package repostitory

import (
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
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
	return uuid.New(), nil
}
