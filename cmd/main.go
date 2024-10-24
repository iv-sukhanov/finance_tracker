package main

import (
	"context"
	"iv-sukhanov/finance_tracker/internal/spendings"
	inith "iv-sukhanov/finance_tracker/internal/utils/init"
)

func main() {

	log := inith.NewLogger()

	db, closeDB, err := inith.ConnectDBFromEnv(context.Background())
	if err != nil {
		log.WithError(err).Fatal("ConnectDBFromEnv: failed to connect to db")
	}
	defer closeDB()

	bot, err := inith.InitBotFromEnv()
	if err != nil {
		log.WithError(err).Fatal("InitBotFromEnv: failed to init bot")
	}

	service, err := spendings.New(log, db, bot)
	if err != nil {
		log.WithError(err).Fatal("spendings.New: failed to create service")
	}
	service.RunBot()
}
