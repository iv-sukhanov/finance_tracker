package repository

import (
	"fmt"

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

func (s *CategoryRepo) AddCategories(categories []ftracker.SpendingCategory) ([]uuid.UUID, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("Repostiory.AddCategory: %w", err)
	}

	stmt, err := tx.PrepareNamed(fmt.Sprintf("INSERT INTO %s (user_guid, category, description, amount) VALUES (:user_guid, :category, :description, :amount) RETURNING guid", spendingCategoriesTable))
	if err != nil {
		return nil, fmt.Errorf("Repostiory.AddCategory: %w", err)
	}

	guids := make([]uuid.UUID, len(categories))
	for i, category := range categories {
		if err := stmt.Get(&guids[i], category); err != nil {
			_err := tx.Rollback()
			if _err != nil {
				panic(_err)
			}
			return nil, fmt.Errorf("Repostiory.AddCategory: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return guids, nil
}
