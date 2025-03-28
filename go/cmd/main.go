package main

import (
	"context"
	"os"

	tbot "github.com/iv-sukhanov/finance_tracker/internal/bot"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"
	"github.com/iv-sukhanov/finance_tracker/internal/service"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	log := utils.NewLogger(argLoggerLevel).WithField("app", argAppName)

	db, closeDB, err := utils.NewPostgresDB(utils.ParamsPostgresDB{
		User:     argPostgresUser,
		Password: argPostgresPassword,
		Host:     argPostgresHost,
		Port:     argPostgresPort,
		DBName:   argPostgresNameDB,
		AppName:  argAppName,
	})
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to DB")
	}
	defer closeDB()

	log.Info("Connected to DB", db.Stats())

	bot, err := utils.NewBot(argTelegramBotToken, argTelegramBotMode == "true")
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize telegram bot")
	}

	repo := repository.New(db)
	src := service.New(repo)
	telegramBot := tbot.New(src, bot, log.Logger)

	telegramBot.Start(context.Background())
}
