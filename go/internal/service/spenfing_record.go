package service

import (
	"time"

	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"
)

type (
	// RecordService implements the SpendingRecord interface.
	RecordService struct {
		repo repository.SpendingRecord
	}

	// RecordOption is a function to modify the RecordOptions.
	RecordOption func(*repository.RecordOptions)

	// RecordOrder defines the order in which records can be sorted
	// It is some sort of enum for the order of records.
	RecordOrder int
)

const (
	OrderRecordsDefault     RecordOrder = iota // default order
	OrderRecordsByAmount                       // order by amount
	OrderRecordsByCreatedAt                    // order by created_at
	OrderRecordsByUpdatedAt                    // order by updated_at
)

// NewRecordService creates a new instance of RecordService with the provided repository.
func NewRecordService(repo repository.SpendingRecord) *RecordService {
	return &RecordService{
		repo: repo,
	}
}

// SpendingRecordsWithLimit is a function that sets the limit for the number of records to be returned.
func (RecordService) SpendingRecordsWithLimit(limit int) RecordOption {
	return func(o *repository.RecordOptions) {
		o.Limit = limit
	}
}

// SpendingRecordsWithTimeFrame is a function that sets the time frame for the records to be returned.
func (RecordService) SpendingRecordsWithTimeFrame(from, to time.Time) RecordOption {
	return func(o *repository.RecordOptions) {
		o.TimeFrom = from
		o.TimeTo = to
		o.ByTime = true
	}
}

// SpendingRecordsWithGUIDs is a function that sets the GUIDs for the records to be returned.
func (RecordService) SpendingRecordsWithGUIDs(guids []uuid.UUID) RecordOption {
	return func(o *repository.RecordOptions) {
		o.GUIDs = guids
	}
}

// SpendingRecordsWithCategoryGUIDs is a function that sets the category GUIDs for the records to be returned.
func (RecordService) SpendingRecordsWithCategoryGUIDs(guids []uuid.UUID) RecordOption {
	return func(o *repository.RecordOptions) {
		o.CategoryGUIDs = guids
	}
}

// SpendingRecordsWithOrder is a function that sets the order of the records to be returned.
func (RecordService) SpendingRecordsWithOrder(order RecordOrder, asc bool) RecordOption {
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

// GetRecords retrieves a list of spending records based on the provided options.
//
// Parameters:
//   - options: A variadic list of RecordOption functions used to configure the query options.
//
// Returns:
//   - []ftracker.SpendingRecord: A slice of spending records that match the query options.
//   - error: An error if the operation fails, otherwise nil.
func (s *RecordService) GetRecords(options ...RecordOption) ([]ftracker.SpendingRecord, error) {
	var opts repository.RecordOptions
	for _, option := range options {
		option(&opts)
	}

	return s.repo.GetRecords(opts)
}

// AddRecords adds multiple spending records to the repository.
//
// Parameters:
//   - records: A slice of SpendingRecord objects to be added.
//
// Returns:
//   - A slice of UUIDs representing the IDs of the newly added records.
//   - An error if the operation fails, or nil if it succeeds.
func (s *RecordService) AddRecords(records []ftracker.SpendingRecord) ([]uuid.UUID, error) {
	return s.repo.AddRecords(records)
}
