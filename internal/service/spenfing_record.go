package service

import (
	"time"

	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"
)

type (
	RecordService struct {
		repo repository.SpendingRecord
	}

	RecordOption func(*repository.RecordOptions)
	RecordOrder  int
)

const (
	OrderRecordsDefault RecordOrder = iota
	OrderRecordsByAmount
	OrderRecordsByCreatedAt
	OrderRecordsByUpdatedAt
)

func NewRecordService(repo repository.SpendingRecord) *RecordService {
	return &RecordService{
		repo: repo,
	}
}

func (RecordService) WithLimit(limit int) RecordOption {
	return func(o *repository.RecordOptions) {
		o.Limit = limit
	}
}

func (RecordService) WithTimeFrame(from, to time.Time) RecordOption {
	return func(o *repository.RecordOptions) {
		o.TimeFrom = from
		o.TimeTo = to
		o.ByTime = true
	}
}

func (RecordService) WithGUIDs(guids []uuid.UUID) RecordOption {
	return func(o *repository.RecordOptions) {
		o.GUIDs = guids
	}
}

func (RecordService) WithCategoryGUIDs(guids []uuid.UUID) RecordOption {
	return func(o *repository.RecordOptions) {
		o.CategoryGUIDs = guids
	}
}

func (RecordService) WithOrder(order RecordOrder, asc bool) RecordOption {
	repOrder := repository.RecordOrder{Asc: asc}
	switch order {
	case OrderRecordsByAmount:
		repOrder.Column = "amount"
	case OrderRecordsByCreatedAt:
		repOrder.Column = "created_at"
	case OrderRecordsByUpdatedAt:
		repOrder.Column = "updated_at"
	}

	return func(o *repository.RecordOptions) {
		o.Order = repOrder
	}
}

func (s *RecordService) GetRecords(options ...RecordOption) ([]ftracker.SpendingRecord, error) {
	var opts repository.RecordOptions
	for _, option := range options {
		option(&opts)
	}

	return s.repo.GetRecords(opts)
}

func (s *RecordService) AddRecords(records []ftracker.SpendingRecord) ([]uuid.UUID, error) {
	return s.repo.AddRecords(records)
}
