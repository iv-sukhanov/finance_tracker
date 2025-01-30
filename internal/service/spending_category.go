package service

import (
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"
)

type (
	CategoryService struct {
		repo repository.SpendingCategory
	}

	CategoryOption func(*repository.CategoryOptions)
)

func NewCategoryService(repo repository.SpendingCategory) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (CategoryService) WithLimit(limit int) CategoryOption {
	return func(o *repository.CategoryOptions) {
		o.Limit = limit
	}
}

func (CategoryService) WithOrder(order repository.CategoryOrder) CategoryOption {
	return func(o *repository.CategoryOptions) {
		o.Order = order
	}
}

func (CategoryService) WithGUIDs(guids []uuid.UUID) CategoryOption {
	return func(o *repository.CategoryOptions) {
		o.GUIDs = guids
	}
}

func (CategoryService) WithUserGUIDs(guids []uuid.UUID) CategoryOption {
	return func(o *repository.CategoryOptions) {
		o.UserGUIDs = guids
	}
}

func (CategoryService) WithCategories(categories []string) CategoryOption {
	return func(o *repository.CategoryOptions) {
		o.Categories = categories
	}
}
func (s *CategoryService) GetCategories(options ...CategoryOption) ([]ftracker.SpendingCategory, error) {
	var opts repository.CategoryOptions
	for _, option := range options {
		option(&opts)
	}

	return s.repo.GetCategories(opts)
}

func (s *CategoryService) AddCategories(categories []ftracker.SpendingCategory) ([]uuid.UUID, error) {
	return s.repo.AddCategories(categories)
}
