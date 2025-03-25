package repository

import (
	"fmt"

	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	"github.com/jmoiron/sqlx"
)

type (
	// CategoryRepo implements the SpendingCategory interface.
	CategoryRepo struct {
		db *sqlx.DB
	}

	// CategoryOptions defines the options for retrieving spending categories.
	CategoryOptions struct {
		Limit      int
		GUIDs      []uuid.UUID
		UserGUIDs  []uuid.UUID
		Categories []string
		Order      CategoryOrder
	}

	// CategoryOrder defines the order in which categories can be sorted.
	CategoryOrder struct {
		Column string
		Asc    bool
	}
)

// NewCategoryRepository creates a new instance of CategoryRepo with the provided database connection.
func NewCategoryRepository(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

// GetCategories retrieves a list of spending categories from the database based on the provided options.
//
// Parameters:
//   - opts: A struct containing filtering, ordering, and limiting options for the query.
//
// Returns:
//   - A slice of SpendingCategory objects that match the query criteria.
//   - An error if the query fails, or nil if successful.
func (c *CategoryRepo) GetCategories(opts CategoryOptions) ([]ftracker.SpendingCategory, error) {

	whereClause := utils.BindWithOp("AND", true,
		utils.MakeIn("guid", utils.UUIDsToStrings(opts.GUIDs)...),
		utils.MakeIn("user_guid", utils.UUIDsToStrings(opts.UserGUIDs)...),
		utils.MakeIn("category", opts.Categories...),
	)

	query := fmt.Sprintf("SELECT guid, user_guid, category, description, amount, created_at, updated_at FROM %s %s %s %s",
		spendingCategoriesTable,
		whereClause,
		utils.MakeOrderBy(opts.Order.Column, opts.Order.Asc),
		utils.MakeLimit(opts.Limit),
	)

	var categories []ftracker.SpendingCategory
	err := c.db.Select(&categories, query)
	if err != nil {
		return nil, fmt.Errorf("Repostiory.GetCategories: %w", err)
	}

	return categories, nil
}

// AddCategories inserts multiple spending categories into the database and returns their generated UUIDs.
//
// Parameters:
//   - categories: A slice of SpendingCategory objects to be added to the database.
//
// Returns:
//   - A slice of UUIDs corresponding to the inserted categories.
//   - An error if the operation fails at any point.
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
