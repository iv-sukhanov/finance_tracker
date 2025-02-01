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
		Limit      int
		GUIDs      []uuid.UUID
		UserGUIDs  []uuid.UUID
		Categories []string
		Order      CategoryOrder
	}

	CategoryOption func(*CategoryOptions)
	CategoryOrder  int
)

const (
	DefaultOrder CategoryOrder = iota
	LastModifiedOrder
	AmountOrder
	AlphabeticalOrder
)

func NewCategoryRepository(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (c *CategoryRepo) GetCategories(opts CategoryOptions) ([]ftracker.SpendingCategory, error) {

	whereClause := utils.BindWithOp("AND", true,
		utils.MakeIn("guid", utils.UUIDsToStrings(opts.GUIDs)...),
		utils.MakeIn("user_guid", utils.UUIDsToStrings(opts.UserGUIDs)...),
		utils.MakeIn("category", opts.Categories...),
	)

	query := fmt.Sprintf("SELECT guid, user_guid, category, description, amount, created_at, updated_at FROM %s %s %s %s",
		spendingCategoriesTable,
		whereClause,
		makeOrderBy(opts.Order),
		utils.MakeLimit(opts.Limit),
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

func makeOrderBy(order CategoryOrder) string {
	switch order {
	case DefaultOrder:
		return ""
	case LastModifiedOrder:
		return "ORDER BY updated_at DESC"
	case AmountOrder:
		return "ORDER BY amount DESC"
	case AlphabeticalOrder:
		return "ORDER BY category ASC"
	default:
		return ""
	}
}
