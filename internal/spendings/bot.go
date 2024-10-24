package spendings

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	offset  = 0
	timeout = 60
)

func newTelegramBot(t *tgbotapi.BotAPI) *telegramBot {
	updates := t.GetUpdatesChan(tgbotapi.UpdateConfig{Offset: offset, Timeout: timeout})
	return &telegramBot{bot: t, updates: updates}
}
