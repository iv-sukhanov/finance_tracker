package service

import (
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"
)

type (
	// CategoryService implements the SpendingCategory interface.
	CategoryService struct {
		repo repository.SpendingCategory
	}

	// CategoryOption is a function to modify the CategoryOptions.
	CategoryOption func(*repository.CategoryOptions)

	// CategoryOrder defines the order in which categories can be sorted
	// It is some sort of enum for the order of categories.
	CategoryOrder int
)

const (
	OrderCategoriesDefault     CategoryOrder = iota // default order
	OrderCategoriesByCategory                       // order alphabetically by category name
	OrderCategoriesByAmount                         // order by amount
	OrderCategoriesByCreatedAt                      // order by created_at
	OrderCategoriesByUpdatedAt                      //order by updated_at
)

// NewCategoryService creates a new instance of CategoryService with the provided repository.
func NewCategoryService(repo repository.SpendingCategory) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

// SpendingCategoriesWithLimit is a function that sets the limit for the number of categories to be returned.
func (CategoryService) SpendingCategoriesWithLimit(limit int) CategoryOption {
	return func(o *repository.CategoryOptions) {
		o.Limit = limit
	}
}

// SpendingCategoriesWithOrder is a function that sets the order of the categories to be returned.
func (CategoryService) SpendingCategoriesWithOrder(order CategoryOrder, asc bool) CategoryOption {

	repOrder := repository.CategoryOrder{Asc: asc}
	switch order {
	case OrderCategoriesByCategory:
		repOrder.Column = "category"
	case OrderCategoriesByAmount:
		repOrder.Column = "amount"
	case OrderCategoriesByCreatedAt:
		repOrder.Column = "created_at"
	case OrderCategoriesByUpdatedAt:
		repOrder.Column = "updated_at"
	}

	return func(o *repository.CategoryOptions) {
		o.Order = repOrder
	}
}

// SpendingCategoriesWithGUIDs is a function that sets the GUIDs of the categories to be returned.
func (CategoryService) SpendingCategoriesWithGUIDs(guids []uuid.UUID) CategoryOption {
	return func(o *repository.CategoryOptions) {
		o.GUIDs = guids
	}
}

// SpendingCategoriesWithUserGUIDs is a function that sets the user GUIDs of the categories to be returned.
func (CategoryService) SpendingCategoriesWithUserGUIDs(guids []uuid.UUID) CategoryOption {
	return func(o *repository.CategoryOptions) {
		o.UserGUIDs = guids
	}
}

// SpendingCategoriesWithCategories is a function that sets the category names of the categories to be returned.
func (CategoryService) SpendingCategoriesWithCategories(categories []string) CategoryOption {
	return func(o *repository.CategoryOptions) {
		o.Categories = categories
	}
}

// GetCategories retrieves a list of spending categories based on the provided options.
//
// Parameters:
//   - options: A variadic list of CategoryOption functions used to configure the query options.
//
// Returns:
//   - []ftracker.SpendingCategory: A slice of spending categories that match the query options.
//   - error: An error if the operation fails, or nil if successful.
func (s *CategoryService) GetCategories(options ...CategoryOption) ([]ftracker.SpendingCategory, error) {
	var opts repository.CategoryOptions
	for _, option := range options {
		option(&opts)
	}

	return s.repo.GetCategories(opts)
}

// AddCategories adds a list of spending categories to the repository.
//
// Parameters:
//   - categories: A slice of SpendingCategory objects to be added.
//
// Returns:
//   - A slice of UUIDs representing the IDs of the newly added categories.
//   - An error if the operation fails, or nil if it succeeds.
func (s *CategoryService) AddCategories(categories []ftracker.SpendingCategory) ([]uuid.UUID, error) {
	return s.repo.AddCategories(categories)
}
