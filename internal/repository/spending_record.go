package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	"github.com/jmoiron/sqlx"
)

type (
	RecordRepo struct {
		db *sqlx.DB
	}

	RecordOptions struct {
		limit          int //-1 for no limit
		timeFrom       time.Time
		timeTo         time.Time
		byTime         bool
		guids          []uuid.UUID
		category_guids []uuid.UUID
	}

	RecordOption func(*RecordOptions)
)

func NewRecordRepository(db *sqlx.DB) *RecordRepo {
	return &RecordRepo{db: db}
}

func (r *RecordRepo) AddRecords(records []ftracker.SpendingRecord) ([]uuid.UUID, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("Repostiory.AddRecords: %w", err)
	}

	stmtIn, err := tx.PrepareNamed(fmt.Sprintf("INSERT INTO %s (category_guid, amount, description) VALUES (:category_guid, :amount, :description) RETURNING guid", spendingRecordsTable))
	if err != nil {
		return nil, fmt.Errorf("Repostiory.AddRecords: %w", err)
	}
	stmtUpd, err := tx.PrepareNamed(fmt.Sprintf("UPDATE %s SET amount = amount + :amount WHERE guid = :category_guid", spendingCategoriesTable))
	if err != nil {
		return nil, fmt.Errorf("Repostiory.AddRecords: %w", err)
	}

	guids := make([]uuid.UUID, len(records))
	for i, record := range records {

		if _, err := stmtUpd.Exec(record); err != nil {
			_err := tx.Rollback()
			if _err != nil {
				panic(_err)
			}
			return nil, fmt.Errorf("Repostiory.AddRecords: %w", err)
		}

		if err := stmtIn.Get(&guids[i], record); err != nil {
			_err := tx.Rollback()
			if _err != nil {
				panic(_err)
			}
			return nil, fmt.Errorf("Repostiory.AddRecords: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return guids, nil
}

func (r *RecordRepo) GetAllRecords() ([]ftracker.SpendingRecord, error) {

	var records []ftracker.SpendingRecord
	err := r.db.Select(&records, fmt.Sprintf("SELECT guid, category_guid, amount, description FROM %s", spendingRecordsTable))
	if err != nil {
		return nil, fmt.Errorf("Repostiory.GetAllRecords: %w", err)
	}

	return records, nil
}

func WithLimit(limit int) RecordOption {
	return func(o *RecordOptions) {
		o.limit = limit
	}
}

func WithTimeFrame(from, to time.Time) RecordOption {
	return func(o *RecordOptions) {
		o.timeFrom = from
		o.timeTo = to
		o.byTime = true
	}
}

func WithGUIDs(guids []uuid.UUID) RecordOption {
	return func(o *RecordOptions) {
		o.guids = guids
	}
}

func WithCategoryGUIDs(guids []uuid.UUID) RecordOption {
	return func(o *RecordOptions) {
		o.category_guids = guids
	}
}

func (r *RecordRepo) GetRecords(optoins ...RecordOption) ([]ftracker.SpendingRecord, error) {
	var opts RecordOptions
	for _, o := range optoins {
		o(&opts)
	}

	var records []ftracker.SpendingRecord
	err := r.db.Select(&records, fmt.Sprintf(
		"SELECT guid, category_guid, amount, description FROM %s %s %s %s %s",
		spendingRecordsTable,
		utils.MakeWhereIn("guid", "AND", utils.UUIDsToStrings(opts.guids)...),
		utils.MakeIn("category_guid", "AND", utils.UUIDsToStrings(opts.category_guids)...),
		utils.MakeTimeFrame("created_at", opts.timeFrom, opts.timeTo, opts.byTime),
		utils.MakeLimit(opts.limit),
	))
	if err != nil {
		return nil, fmt.Errorf("Repostiory.GetRecords: %w", err)
	}

	return records, nil
}

func (r *RecordRepo) GetRecordsByGUIDs(guids []uuid.UUID) ([]ftracker.SpendingRecord, error) {

	var records []ftracker.SpendingRecord
	err := r.db.Select(&records, fmt.Sprintf(
		"SELECT guid, category_guid, amount, description FROM %s %s",
		spendingRecordsTable,
		utils.MakeWhereIn("guid", "", utils.UUIDsToStrings(guids)...)),
	)
	if err != nil {
		return nil, fmt.Errorf("Repostiory.GetRecordsByGUIDs: %w", err)
	}

	return records, nil
}
