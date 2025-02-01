package ftracker

import (
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		GUID       uuid.UUID `json:"guid" db:"guid"`
		Username   string    `json:"username" db:"username"`
		TelegramID string    `json:"telegram_id" db:"telegram_id"`
		CreatedAt  time.Time `json:"created_at" db:"created_at"`
		UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	}

	SpendingCategory struct {
		GUID        uuid.UUID `json:"guid" db:"guid"`
		UserGUID    uuid.UUID `json:"user_guid" db:"user_guid"`
		Category    string    `json:"category" db:"category"`
		Description string    `json:"description" db:"description"`
		Amount      float64   `json:"amount" db:"amount"`
		CreatedAt   time.Time `json:"created_at" db:"created_at"`
		UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	}

	SpendingRecord struct {
		GUID         uuid.UUID `json:"guid" db:"guid"`
		CategoryGUID uuid.UUID `json:"category_guid" db:"category_guid"`
		Amount       float32   `json:"amount" db:"amount"`
		Description  string    `json:"description" db:"description"`
		CreatedAt    time.Time `json:"created_at" db:"created_at"`
		UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	}
)
