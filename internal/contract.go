package ftracker

import "github.com/google/uuid"

type (
	User struct {
		GUID       uuid.UUID `json:"guid" db:"guid"`
		Name       string    `json:"name" db:"name"`
		TelegramID string    `json:"telegram_id" db:"telegram_id"`
	}
)
