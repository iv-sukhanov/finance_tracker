package utils

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// NewBot initializes a new Telegram Bot API client using the provided token.
// It also allows enabling or disabling debug mode for the bot.
//
// Parameters:
//   - token: The authentication token for the Telegram Bot API. Must not be empty.
//   - debugMode: A boolean flag to enable or disable debug mode.
//
// Returns:
//   - bot: A pointer to the initialized tgbotapi.BotAPI instance.
//   - err: An error if the initialization fails, or nil if successful.
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
