package repository

import (
	"fmt"

	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	"github.com/jmoiron/sqlx"
)

type (
	CategoryRepo struct {
		db *sqlx.DB
	}

	CategoryOptions struct {
		limit      int
		guids      []uuid.UUID
		userGUIDs  []uuid.UUID
		categories []string
	}

	CategoryOption func(*CategoryOptions)
)

func NewCategoryRepository(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (CategoryRepo) WithLimit(limit int) CategoryOption {
	return func(o *CategoryOptions) {
		o.limit = limit
	}
}

func (CategoryRepo) WithGUIDs(guids []uuid.UUID) CategoryOption {
	return func(o *CategoryOptions) {
		o.guids = guids
	}
}

func (CategoryRepo) WithUserGUIDs(guids []uuid.UUID) CategoryOption {
	return func(o *CategoryOptions) {
		o.userGUIDs = guids
	}
}

func (CategoryRepo) WithCategories(categories []string) CategoryOption {
	return func(o *CategoryOptions) {
		o.categories = categories
	}
}

func (c *CategoryRepo) GetCategories(opts CategoryOptions) ([]ftracker.SpendingCategory, error) {

	whereClause := utils.BindWithOp("AND", true,
		utils.MakeIn("guid", utils.UUIDsToStrings(opts.guids)...),
		utils.MakeIn("user_guid", utils.UUIDsToStrings(opts.userGUIDs)...),
		utils.MakeIn("category", opts.categories...),
	)

	query := fmt.Sprintf("SELECT guid, user_guid, category, description, amount FROM %s %s %s",
		spendingCategoriesTable,
		whereClause,
		utils.MakeLimit(opts.limit),
	)

	var categories []ftracker.SpendingCategory
	err := c.db.Select(&categories, query)
	if err != nil {
		return nil, fmt.Errorf("Repostiory.GetCategories: %w", err)
	}

	return categories, nil
}

func (c *CategoryRepo) AddCategories(categories []ftracker.SpendingCategory) ([]uuid.UUID, error) {
	tx, err := c.db.Beginx()
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
