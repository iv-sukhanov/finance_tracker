package spendings

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	bot *telegramBot
	log *log.Logger
}

type telegramBot struct {
	bot     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
}
