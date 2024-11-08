package bot

import (
	"os"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"
)

func Test_Run(t *testing.T) {

	tgbot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	require.NoError(t, err)
	// Create a new bot
	bot := NewTelegramBot(nil, tgbot)

	// Run the bot
	bot.Start()
}
