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
		Limit         int
		TimeFrom      time.Time
		TimeTo        time.Time
		ByTime        bool
		GUIDs         []uuid.UUID
		CategoryGUIDs []uuid.UUID
		Order         RecordOrder
	}

	RecordOption func(*RecordOptions)
	RecordOrder  struct {
		Column string
		Asc    bool
	}
)

func NewRecordRepository(db *sqlx.DB) *RecordRepo {
	return &RecordRepo{db: db}
}

func (r *RecordRepo) GetRecords(opts RecordOptions) ([]ftracker.SpendingRecord, error) {

	whereClause := utils.BindWithOp("AND", true,
		utils.MakeIn("guid", utils.UUIDsToStrings(opts.GUIDs)...),
		utils.MakeIn("category_guid", utils.UUIDsToStrings(opts.CategoryGUIDs)...),
		utils.MakeTimeFrame("updated_at", opts.TimeFrom, opts.TimeTo, opts.ByTime),
	)

	query := fmt.Sprintf(
		"SELECT guid, category_guid, amount, description, created_at, updated_at FROM %s %s %s %s",
		spendingRecordsTable,
		whereClause,
		utils.MakeOrderBy(opts.Order.Column, opts.Order.Asc),
		utils.MakeLimit(opts.Limit),
	)

	var records []ftracker.SpendingRecord
	err := r.db.Select(&records, query)
	if err != nil {
		return nil, fmt.Errorf("Repostiory.GetRecords: %w", err)
	}

	return records, nil
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
