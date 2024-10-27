package ftracker

import "github.com/google/uuid"

type (
	User struct {
		GUID       uuid.UUID `json:"guid" db:"guid"`
		Username   string    `json:"username" db:"username"`
		TelegramID string    `json:"telegram_id" db:"telegram_id"`
	}
)
