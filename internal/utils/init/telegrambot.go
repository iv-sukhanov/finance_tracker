package inith

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// InitBot initializes telegram bot from arguments
func NewBot(token string, debugMode bool) (bot *tgbotapi.BotAPI, err error) {

	if token == "" {
		return nil, fmt.Errorf("InitBot: token is empty")
	}

	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("InitBotFromArgs: %w", err)
	}
	bot.Debug = debugMode

	return bot, nil
}
