package ftracker

import (
	"time"

	"github.com/google/uuid"
)

type (
	//Represents a user
	//GUID - unique identifier of the user
	//Username - telegram username of the user
	//TelegramID - telegram id of the user
	//CreatedAt - time when the user was created
	//UpdatedAt - time when the user was updated last time
	User struct {
		GUID       uuid.UUID `json:"guid" db:"guid"`
		Username   string    `json:"username" db:"username"`
		TelegramID string    `json:"telegram_id" db:"telegram_id"`
		CreatedAt  time.Time `json:"created_at" db:"created_at"`
		UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	}

	//SpendingCategory represents a spending category
	//GUID - unique identifier of the category
	//UserGUID - unique identifier of the user to whom the category belongs
	//Category - name of the category
	//Description - description of the category
	//Amount - amount of money spent in the category
	//CreatedAt - time when the category was created
	//UpdatedAt - time when the category was updated last time
	SpendingCategory struct {
		GUID        uuid.UUID `json:"guid" db:"guid"`
		UserGUID    uuid.UUID `json:"user_guid" db:"user_guid"`
		Category    string    `json:"category" db:"category"`
		Description string    `json:"description" db:"description"`
		Amount      uint64    `json:"amount" db:"amount"`
		CreatedAt   time.Time `json:"created_at" db:"created_at"`
		UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	}

	//SpendingRecord represents a spending record
	//GUID - unique identifier of the record
	//CategoryGUID - unique identifier of the category to which the record belongs
	//Amount - amount of money spent in the record
	//Description - description of the record
	//CreatedAt - time when the record was created
	//UpdatedAt - time when the record was updated last time
	SpendingRecord struct {
		GUID         uuid.UUID `json:"guid" db:"guid"`
		CategoryGUID uuid.UUID `json:"category_guid" db:"category_guid"`
		Amount       uint32    `json:"amount" db:"amount"`
		Description  string    `json:"description" db:"description"`
		CreatedAt    time.Time `json:"created_at" db:"created_at"`
		UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	}
)
