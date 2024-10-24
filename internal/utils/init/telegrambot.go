package inith

import (
	"fmt"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// InitBot initializes telegram bot from arguments
func InitBot(token string, debugMode bool) (bot *tgbotapi.BotAPI, err error) {

	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("InitBotFromArgs: %w", err)
	}
	bot.Debug = debugMode

	return bot, nil
}

// InitBotFromEnv initializes telegram bot from environment variables
func InitBotFromEnv() (bot *tgbotapi.BotAPI, err error) {

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	mode := os.Getenv("TELEGRAM_DEBUG_MODE")

	if token == "" || mode == "" {
		var builder strings.Builder
		builder.Grow(10) // "token " -> 6, "mode" -> 4
		if token == "" {
			builder.WriteString("token ")
		}
		if mode == "" {
			builder.WriteString("mode")
		}
		return nil, fmt.Errorf("InitBot: args are empty: %s", builder.String())
	}

	return InitBot(token, mode == "true")
}
