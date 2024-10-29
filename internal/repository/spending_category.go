package repository

import (
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/jmoiron/sqlx"
)

type CategoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (s *CategoryRepo) AddCategory(category ftracker.SpendingCategory) (uuid.UUID, error) {
	return uuid.Nil, nil
}
