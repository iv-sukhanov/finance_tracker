package ftracker

import "github.com/google/uuid"

type (
	User struct {
		GUID       uuid.UUID `json:"guid" db:"guid"`
		Username   string    `json:"username" db:"username"`
		TelegramID string    `json:"telegram_id" db:"telegram_id"`
	}

	SpendingCategory struct {
		GUID        uuid.UUID `json:"guid" db:"guid"`
		UserGUID    uuid.UUID `json:"user_guid" db:"user_guid"`
		Category    string    `json:"category" db:"category"`
		Description string    `json:"description" db:"description"`
		Amount      float64   `json:"amount" db:"amount"`
	}

	SpendingRecord struct {
		GUID         uuid.UUID `json:"guid" db:"guid"`
		CategoryGUID uuid.UUID `json:"category_guid" db:"category_guid"`
		Amount       float32   `json:"amount" db:"amount"`
		Description  string    `json:"description" db:"description"`
	}
)
