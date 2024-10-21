package spendings

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// InitBot initializes telegram bot from arguments
func (s *Service) InitBotFromArgs(token string, debugMode bool) {

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		s.log.WithError(err).Fatal("Failed to initialize bot")
	}

	bot.Debug = debugMode
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	s.bot = &telegramBot{bot: bot, updates: updates}
}

// InitBot initializes telegram bot from environment variables
func (s *Service) InitBot() {

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	mode := os.Getenv("TELEGRAM_DEBUG_MODE")

	if token == "" || mode == "" {
		s.log.WithField("TELEGRAM_BOT_TOKEN", token).WithField("TELEGRAM_DEBUG_MODE", mode).Fatal("some of telegram env variables are not set")
	}

	s.InitBotFromArgs(token, mode == "true")
}

func (s *Service) RunBot() {
	for update := range s.bot.updates {
		if update.Message == nil {
			continue
		}

		s.log.WithField("Message", update.Message.Text).Debug("Received update")
	}
}
