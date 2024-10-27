package main

import (
	"os"

	tbot "github.com/iv-sukhanov/finance_tracker/internal/bot"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"
	"github.com/iv-sukhanov/finance_tracker/internal/service"
	inith "github.com/iv-sukhanov/finance_tracker/internal/utils/init"
)

var (
	argPostgresUser     = os.Getenv("POSTGRES_USER")
	argPostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	argPostgresHost     = os.Getenv("POSTGRES_HOST")
	argPostgresPort     = os.Getenv("POSTGRES_PORT")
	argPostgresNameDB   = os.Getenv("POSTGRES_DB")
	argLoggerLevel      = os.Getenv("LOG_LEVEL")
	argAppName          = os.Getenv("APP_NAME")
	argTelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	argTelegramBotMode  = os.Getenv("TELEGRAM_DEBUG_MODE")
)

func main() {

	log := inith.NewLogger(argLoggerLevel).WithField("app", argAppName)

	db, err := inith.NewPostgresDB(inith.ParamsPostgresDB{
		User:     argPostgresUser,
		Password: argPostgresPassword,
		Host:     argPostgresHost,
		Port:     argPostgresPort,
		DBName:   argPostgresNameDB,
	})
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to DB")
	}

	bot, err := inith.NewBot(argTelegramBotToken, argTelegramBotMode == "true")
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize telegram bot")
	}

	repo := repository.NewRepostitory(db)
	src := service.New(repo, log)
	telegramBot := tbot.NewTelegramBot(src, bot)

	telegramBot.Start()
}
