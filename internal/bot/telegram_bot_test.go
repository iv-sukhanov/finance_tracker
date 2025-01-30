package bot

import (
	"context"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"
	"github.com/iv-sukhanov/finance_tracker/internal/service"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_Run(t *testing.T) {

	basePath, err := os.Getwd()
	require.NoError(t, err)
	basePath += "/../../migrations/"
	testContainerDB, stop, err := utils.NewPGContainer(
		basePath + "000001_init.up.sql",
	)
	require.NoError(t, err)
	defer stop()

	repo := repository.NewRepostitory(testContainerDB)
	srv := service.New(repo, logrus.New().WithField("test", "finance_tracker"))

	tgbot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	require.NoError(t, err)
	// Create a new bot
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	bot := NewTelegramBot(srv, tgbot, log)

	// Run the bot
	bot.Start(context.Background())
}
