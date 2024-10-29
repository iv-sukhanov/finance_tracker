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
		Amount      uint32    `json:"amount" db:"amount"`
	}
)
