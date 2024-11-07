package repository

import (
	"fmt"

	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	"github.com/jmoiron/sqlx"
)

type (
	UserRepo struct {
		db *sqlx.DB
	}

	UserOptions struct {
		Limit       int
		GUIDs       []uuid.UUID
		Usernames   []string
		TetegramIDs []string
	}
)

func NewUserRepository(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (s *UserRepo) GetUsers(opts UserOptions) ([]ftracker.User, error) {

	whereClause := utils.BindWithOp("AND", true,
		utils.MakeIn("guid", utils.UUIDsToStrings(opts.GUIDs)...),
		utils.MakeIn("username", opts.Usernames...),
		utils.MakeIn("telegram_id", opts.TetegramIDs...),
	)

	query := fmt.Sprintf(
		"SELECT guid, username, telegram_id FROM %s %s %s",
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
