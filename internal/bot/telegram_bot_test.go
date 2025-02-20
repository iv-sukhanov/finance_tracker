package bot

import (
	"context"
	"os"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/mock/gomock"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"
	"github.com/iv-sukhanov/finance_tracker/internal/service"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
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
	srv := service.New(repo)

	tgbot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	require.NoError(t, err)
	// Create a new bot
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	bot := NewTelegramBot(srv, tgbot, log)

	// Run the bot
	bot.Start(context.Background())
}

func TestTelegramBot_HandleUpdate(t *testing.T) {

	log := logrus.New()
	log.SetLevel(logrus.FatalLevel)

	tt := []struct {
		name             string
		senderBehavior   func(*MockSender)
		sessionsBehavior func(*MockSessions)
		update           tgbotapi.Update
		ctx              context.Context
	}{
		// TODO: Add test cases.
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mockSender := NewMockSender(controller)
			mockSessions := NewMockSessions(controller)

			b := &TelegramBot{
				sender:   mockSender,
				sessions: mockSessions,
				log:      log,
				service:  nil,
				api:      nil,
			}
			b.HandleUpdate(tc.ctx, tc.update)
		})
	}
}
