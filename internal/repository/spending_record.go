package repository

import (
	"fmt"

	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	"github.com/jmoiron/sqlx"
)

type RecordRepo struct {
	db *sqlx.DB
}

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

func (r *RecordRepo) GetRecordsByGUIDs(guids []uuid.UUID) ([]ftracker.SpendingRecord, error) {

	var records []ftracker.SpendingRecord
	err := r.db.Select(&records, fmt.Sprintf(
		"SELECT guid, category_guid, amount, description FROM %s %s",
		spendingRecordsTable,
		utils.MakeWhereIn("guid", utils.UUIDsToStrings(guids)...)),
	)
	if err != nil {
		return nil, fmt.Errorf("Repostiory.GetRecordsByGUIDs: %w", err)
	}

	return records, nil
}
